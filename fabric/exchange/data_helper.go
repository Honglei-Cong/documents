/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"unicode/utf8"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type typeAttributeInfo struct {
	Key bool `json:"key"`
}

type typeInfoStruct struct {
	Keys    map[string]typeAttributeInfo `json:"keys"`
	Name    string                       `json:"name"`
	HasSeq  bool                         `json:"has_seq"`
	NextSeq uint64                       `json:"next_seq"`
}

type KVObject interface {
	GetID() string
	SetID(id string)
	IsSeqID() bool
}

type QueryKV struct {
	Key       string `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
	Value     []byte `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

type QueryIterator interface {
	Close() error
	HasNext() bool
	Next() (*QueryKV, error)
}

var errNotFound = errors.New("结果未找到.")

const (
	typeKeyPrefix         = "\x10"
	compositeKeyNamespace = "\x11"
	minUnicodeRuneValue   = 0            //U+0000
	maxUnicodeRuneValue   = utf8.MaxRune //U+10FFFF - maximum (and unallocated) code point
)

func createCompositeKey(objectType string, attributes []string) (string, error) {
	ck := compositeKeyNamespace + objectType + string(minUnicodeRuneValue)
	for _, att := range attributes {
		ck += att + string(minUnicodeRuneValue)
	}
	return ck, nil
}

func getTypeName(obj interface{}) string {
	typ := reflect.TypeOf(obj)
	if typ.Kind() == reflect.Ptr {
		return typ.Elem().Name()
	} else if typ.Kind() == reflect.Struct {
		return typ.Name()
	}
	return ""
}

func getTypeKey(typename string) string {
	return typeKeyPrefix + typename
}
func getObjectSubkey(stub shim.ChaincodeStubInterface, typename string, subkeys []string) (string, error) {
	return createCompositeKey(getTypeKey(typename), subkeys)
}

func getTypeInfo(stub shim.ChaincodeStubInterface, typename string) (*typeInfoStruct, error) {
	bytes, err := stub.GetState(getTypeKey(typename))
	if err != nil {
		return nil, err
	}
	typeInfo := &typeInfoStruct{
		Keys: make(map[string]typeAttributeInfo),
	}
	err = json.Unmarshal(bytes, &typeInfo)
	if err != nil {
		return nil, fmt.Errorf("解析成json失败 %s", string(bytes))
	}
	if typeInfo.Name != typename {
		return nil, errors.New("类型信息不正确")
	}

	return typeInfo, nil
}

func putTypeInfo(stub shim.ChaincodeStubInterface, typename string, typeInfo *typeInfoStruct) error {
	bytes, err := json.Marshal(typeInfo)
	if err != nil {
		return err
	}
	return stub.PutState(getTypeKey(typename), bytes)
}

func getObjectKVs(stub shim.ChaincodeStubInterface, obj KVObject) (map[string][]byte, error) {
	kvs := make(map[string][]byte)
	typename := getTypeName(obj)
	typeInfo, err := getTypeInfo(stub, typename)
	if err != nil {
		return kvs, fmt.Errorf("获取类型 %s 信息失败: %s", typename, err)
	}

	objID := obj.GetID()
	if len(objID) == 0 {
		if !obj.IsSeqID() {
			return kvs, fmt.Errorf("存取类型 %s 失败，无效ID", typename)
		}
		objID = fmt.Sprintf("%016x", typeInfo.NextSeq)
		obj.SetID(objID)
	}

	objBytes, err := json.Marshal(obj)
	if err != nil {
		return kvs, fmt.Errorf("解析成json失败: %s", err)
	}

	jsonObj := make(map[string]interface{})
	if err := json.Unmarshal(objBytes, &jsonObj); err != nil {
		return kvs, err
	}

	primKey, err := createCompositeKey(typename, []string{objID})
	if err != nil {
		return kvs, err
	}
	kvs[primKey] = objBytes

	for key, keyInfo := range typeInfo.Keys {
		if keyInfo.Key {
			attr := fmt.Sprintf("%v", jsonObj[key])
			subKey, _ := getObjectSubkey(stub, typename, []string{key, attr, objID})
			kvs[subKey] = []byte{0x00}
		}
	}

	return kvs, nil
}

func incrObjectTypeNextSeq(stub shim.ChaincodeStubInterface, typename string) error {
	typeInfo, err := getTypeInfo(stub, typename)
	if err != nil {
		return err
	}
	if !typeInfo.HasSeq {
		return nil
	}

	typeInfo.NextSeq += 1
	return putTypeInfo(stub, typename, typeInfo)
}

func CreateObjectType(stub shim.ChaincodeStubInterface, obj KVObject, keys []string) error {
	typ := reflect.TypeOf(obj)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	} else if typ.Kind() != reflect.Struct {
		return errors.New("对象类型不正确.")
	}

	typename := typ.Name()
	typeInfo := &typeInfoStruct{
		Keys:    make(map[string]typeAttributeInfo),
		Name:    typename,
		HasSeq:  obj.IsSeqID(),
		NextSeq: 1,
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i).Tag.Get("json")
		typeInfo.Keys[field] = typeAttributeInfo{Key: false}
	}

	// check all keys are in type attributes
	for _, key := range keys {
		if _, ok := typeInfo.Keys[key]; !ok {
			return errors.New("Invalid Key: " + key)
		}
		attrInfo := typeInfo.Keys[key]
		attrInfo.Key = true
		typeInfo.Keys[key] = attrInfo
	}

	return putTypeInfo(stub, typename, typeInfo)
}

func IsObjectTypeExist(stub shim.ChaincodeStubInterface, obj KVObject) bool {
	info, err := getTypeInfo(stub, getTypeName(obj))
	if err != nil {
		return false
	}
	return info != nil
}

func GetObject(stub shim.ChaincodeStubInterface, obj KVObject, keys map[string]string) error {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return errors.New("对象类型不正确. 需要ptr类型.")
	}

	if len(keys) == 0 {
		if len(obj.GetID()) == 0 {
			return errors.New("查询不正确.  没有提供键.")
		}
		var k string
		var err error
		var valueBytes []byte
		if k, err = createCompositeKey(getTypeName(obj), []string{obj.GetID()}); err != nil {
			return err
		}
		if valueBytes, err = stub.GetState(k); err != nil {
			return err
		}
		return json.Unmarshal(valueBytes, obj)
	}

	ite, err := QueryObjects(stub, obj, keys)
	if err != nil {
		return err
	}

	defer ite.Close()
	if ite.HasNext() {
		_, err := ite.Next()
		if err != nil {
			return err
		}
		return nil
	}
	return errNotFound
}

/**  GetObjectsIte  use example：
ite, err := GetObjectsIte(stub, salesOrg, keys)
if err != nil {
	fmt.Println("querySalesOrg  err")
	return err
}

for(ite.HasNext() ) {
	_, v, err := ite.Next()
	if err != nil {
		fmt.Println("Next error")
		return err
	}

	salesOrg1 := &SalesOrg{}
	if err = json.Unmarshal(v, salesOrg1); err != nil {
		fmt.Println("Unmarshal error")
		return err
	}
	salesOrgs  = append(salesOrgs, *salesOrg1)
}
*/
func GetObjectsIte(stub shim.ChaincodeStubInterface, obj KVObject, keys map[string]string) (QueryIterator, error) {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return nil, errors.New("对象类型不正确. 需要ptr类型.")
	}

	if len(keys) == 0 {
		k, err := createCompositeKey(getTypeName(obj), []string{})
		if err != nil {
			return nil, errors.New("查询不正确.  创建组合键失败.")
		}
		ite, err := stub.GetStateByRange(k, k+string(maxUnicodeRuneValue))
		if err != nil {
			return nil, fmt.Errorf("查询失败: %s", err)
		}
		return &StateQueryIterator{ite}, nil
	}

	iter, err := QueryObjects(stub, obj, keys)
	if err != nil {
		return nil, err
	}

	return iter, nil
}

func PutObject(stub shim.ChaincodeStubInterface, obj KVObject) error {
	if len(obj.GetID()) == 0 && !obj.IsSeqID() {
		obj.SetID(stub.GetTxID())
	}
	return putObjectInternal(stub, obj)
}

func putObjectInternal(stub shim.ChaincodeStubInterface, obj KVObject) error {
	typename := getTypeName(obj)
	kvs, err := getObjectKVs(stub, obj)
	if err != nil {
		return fmt.Errorf("获取ObjectKVs失败: %s", err)
	}

	for k, v := range kvs {
		err := stub.PutState(k, v)
		if err != nil {
			return fmt.Errorf("putState操作失败: %s", err)
		}
	}
	if err := incrObjectTypeNextSeq(stub, typename); err != nil {
		return fmt.Errorf("更新表 %s 流水号失败: %s", typename, err)
	}
	return nil
}

func UpdateObject(stub shim.ChaincodeStubInterface, obj KVObject) error {
	typename := getTypeName(obj)
	var k string
	var oldValueBytes, newValueBytes []byte
	var err error

	if k, err = createCompositeKey(typename, []string{obj.GetID()}); err != nil {
		return err
	}
	if oldValueBytes, err = stub.GetState(k); err != nil {
		return err
	}
	if newValueBytes, err = json.Marshal(obj); err != nil {
		return err
	}

	// delete old obj
	if err = json.Unmarshal(oldValueBytes, obj); err != nil {
		return err
	}
	if err = DelObject(stub, obj); err != nil {
		return err
	}

	// put new obj
	if err = json.Unmarshal(newValueBytes, obj); err != nil {
		return err
	}
	if err := putObjectInternal(stub, obj); err != nil {
		return err
	}

	return nil
}

func DelObject(stub shim.ChaincodeStubInterface, obj KVObject) error {
	kvs, err := getObjectKVs(stub, obj)
	if err != nil {
		return err
	}

	for k := range kvs {
		err := stub.DelState(k)
		if err != nil {
			return err
		}
	}
	return nil
}

func QueryObjects(stub shim.ChaincodeStubInterface, obj KVObject, keys map[string]string) (QueryIterator, error) {
	typename := getTypeName(obj)
	typeInfo, err := getTypeInfo(stub, typename)
	if err != nil {
		return nil, err
	}

	ite := &ObjectIterator{
		typename:  typename,
		stub:      stub,
		iters:     make([]shim.StateQueryIteratorInterface, 0),
		currObjId: "\x00",
		obj:       obj,
	}
	for key, val := range keys {
		if keyInfo, ok := typeInfo.Keys[key]; ok && keyInfo.Key {
			subKey, _ := getObjectSubkey(stub, typename, []string{key, val})
			subIte, err := stub.GetStateByRange(subKey, subKey+string(maxUnicodeRuneValue))
			if err != nil {
				return nil, err
			}
			ite.iters = append(ite.iters, subIte)
		}
	}

	return ite, nil
}

type StateQueryIterator struct {
	ite shim.StateQueryIteratorInterface
}

func (ite *StateQueryIterator) HasNext() bool {
	return ite.ite.HasNext()
}

func (ite *StateQueryIterator) Next() (*QueryKV, error) {
	kv, err := ite.ite.Next()
	if kv == nil {
		return nil, err
	}
	return &QueryKV{
		Key: kv.Key,
		Value: kv.Value,
	}, err
}

func (ite *StateQueryIterator) Close() error {
	return ite.ite.Close()
}

type ObjectIterator struct {
	typename  string
	stub      shim.ChaincodeStubInterface
	iters     []shim.StateQueryIteratorInterface
	currObjId string
	obj       interface{}
}

func (ite *ObjectIterator) HasNext() bool {
	currObjId := "\x00"
	iteMatched := make([]bool, len(ite.iters))

	for idx := 0; idx < len(ite.iters); idx++ {
		for !iteMatched[idx] {
			if !ite.iters[idx].HasNext() {
				return false
			}
			kv, err := ite.iters[idx].Next()
			if err != nil {
				return false
			}

			_, keys, err := ite.stub.SplitCompositeKey(kv.Key)
			objID := keys[2]
			if objID > currObjId {
				currObjId = objID
				if idx != 0 {
					// loop back
					iteMatched = make([]bool, len(ite.iters))
					iteMatched[idx] = true
					idx = 0
				} else {
					iteMatched[idx] = true
				}
			} else if objID == currObjId {
				iteMatched[idx] = true
			}
		}
	}

	if currObjId == "\x00" {
		return false
	}

	ite.currObjId = currObjId
	return true
}

func (ite *ObjectIterator) Next() (*QueryKV, error) {
	k, err := createCompositeKey(ite.typename, []string{ite.currObjId})
	if err != nil {
		return nil, err
	}
	v, err := ite.stub.GetState(k)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(v, ite.obj); err != nil {
		return nil, err
	}
	return &QueryKV{
		Key:   ite.currObjId,
		Value: v,
	}, err
}

func (ite *ObjectIterator) Close() error {
	for _, it := range ite.iters {
		if err := it.Close(); err != nil {
			return err
		}
	}
	return nil
}
