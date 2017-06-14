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

package jscc

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/op/go-logging"

	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"encoding/json"
	"errors"

	"strconv"

	"github.com/robertkrimen/otto"
)

var logger = logging.MustGetLogger("javascriptcc")

type jCCShim struct {
	gostub  shim.ChaincodeStubInterface
	lastErr error
}

func ottoStringValue(s string) otto.Value {
	v, err := otto.ToValue(s)
	if err != nil {
		logger.Warning("ottoStringValue", err.Error())
	}
	return v
}

func ottoJsonValue(x interface{}) (otto.Value, error) {
	bytes, err := json.Marshal(x)
	if err != nil {
		logger.Warning("ottoJSONValue", err.Error())
		return otto.NaNValue(), err
	}
	return ottoStringValue(string(bytes)), nil
}

type jCCQueryIterator struct {
	iter         shim.StateQueryIteratorInterface
	currentKey   string
	currentValue string
	lastErr      error
}

func (jccIte *jCCQueryIterator) HasNext() otto.Value {
	logger.Warning("HasNext")
	if jccIte.iter.HasNext() {
		return otto.TrueValue()
	}
	return otto.FalseValue()
}

func (jccIte *jCCQueryIterator) Next() otto.Value {
	logger.Warning("iter Next")
	k, valBytes, err := jccIte.iter.Next()
	if err != nil {
		jccIte.lastErr = err
		logger.Warning("iter Next", err.Error())
		return otto.FalseValue()
	}

	jccIte.currentKey = k
	var val string
	json.Unmarshal(valBytes, val)
	jccIte.currentValue = val
	jccIte.lastErr = nil
	logger.Warning("iter Next", k, val)

	return otto.TrueValue()
}

func (jccIte *jCCQueryIterator) Close() otto.Value {
	logger.Warning("Close")
	if err := jccIte.iter.Close(); err != nil {
		return ottoStringValue(err.Error())
	}
	jccIte.currentKey = ""
	jccIte.currentValue = ""
	jccIte.lastErr = nil
	return otto.NullValue()
}

func (jccIte *jCCQueryIterator) GetCurrentKey() otto.Value {
	logger.Warning("GetCurrentKey")
	return ottoStringValue(jccIte.currentKey)
}

func (jccIte *jCCQueryIterator) GetCurrentValue() otto.Value {
	logger.Warning("GetCurrentValue")
	return ottoStringValue(jccIte.currentValue)
}

func (jccIte *jCCQueryIterator) GetError() otto.Value {
	logger.Warning("GetError")
	return ottoStringValue(jccIte.lastErr.Error())
}

type JavascriptCC struct {
	codepath string
	vm       *otto.Otto
}

func getShimStub(call otto.FunctionCall) (*jCCShim, error) {
	x, _ := call.This.Export()
	switch v := x.(type) {
	case *jCCShim:
		return v, nil
	}
	return nil, errors.New("Invalid shim type")
}

func toShimResponse(v otto.Value) pb.Response {
	if !v.IsDefined() {
		return shim.Success(nil)
	}

	x, _ := v.Export()
	switch resp := x.(type) {
	case pb.Response:
		return resp
	}
	return shim.Error("Invalid format of response from chaincode: " + v.String())
}

func builtinShim_Success(call otto.FunctionCall) otto.Value {
	logger.Warning("Success, arg#", len(call.ArgumentList))
	var response string
	for _, arg := range call.ArgumentList {
		response += arg.String()
	}
	v, _ := call.Otto.ToValue(shim.Success([]byte(response)))
	return v
}

func builtinShim_Error(call otto.FunctionCall) otto.Value {
	logger.Warning("Error, arg#", len(call.ArgumentList))
	var errResponse string
	for _, arg := range call.ArgumentList {
		errResponse += arg.String()
	}
	v, _ := call.Otto.ToValue(shim.Error(errResponse))
	return v
}

func builtinShim_GetLastError(call otto.FunctionCall) otto.Value {
	logger.Warning("GetError, arg#", len(call.ArgumentList))
	x, _ := call.This.Export()
	switch v := x.(type) {
	case *jCCShim:
		if v.lastErr == nil {
			return otto.NullValue()
		}
		return ottoStringValue(v.lastErr.Error())
	}

	return otto.NaNValue()
}

func builtinShim_GetFunction(call otto.FunctionCall) otto.Value {
	logger.Warning("GetFunction, arg#", len(call.ArgumentList))
	stub, err := getShimStub(call)
	if err != nil {
		return otto.NaNValue()
	}
	funcname, _ := stub.gostub.GetFunctionAndParameters()
	return ottoStringValue(funcname)
}

func builtinShim_GetArguments(call otto.FunctionCall) otto.Value {
	logger.Warning("GetArgs, arg#", len(call.ArgumentList))
	stub, err := getShimStub(call)
	if err != nil {
		return otto.NaNValue()
	}
	_, args := stub.gostub.GetFunctionAndParameters()
	logger.Warning("GetArguments", strconv.Itoa(len(args)))
	ret, err := ottoJsonValue(args)
	if err != nil {
		logger.Warning(err.Error())
		stub.lastErr = err
	}
	return ret
}

func builtinShim_GetState(call otto.FunctionCall) otto.Value {
	logger.Warning("GetState, arg#", len(call.ArgumentList))
	stub, err := getShimStub(call)
	if err != nil {
		return otto.NaNValue()
	}
	res, err := stub.gostub.GetState(call.Argument(0).String())
	if err != nil {
		stub.lastErr = err
		return otto.NullValue()
	}

	var value string
	json.Unmarshal(res, &value)
	logger.Warning("GetState", call.Argument(0).String(), value)
	return ottoStringValue(value)
}

func builtinShim_PutState(call otto.FunctionCall) otto.Value {
	logger.Warning("PutState, arg#", len(call.ArgumentList))
	stub, err := getShimStub(call)
	if err != nil {
		return otto.NaNValue()
	}
	if len(call.ArgumentList) < 2 {
		stub.lastErr = errors.New("Invalid arg number for PutState")
		return otto.NullValue()
	}

	logger.Warning("PutState", call.Argument(0).String(), call.Argument(1).String())
	bytes, _ := json.Marshal(call.Argument(1).String())
	if err := stub.gostub.PutState(call.Argument(0).String(), bytes); err != nil {
		stub.lastErr = err
	}

	return otto.NullValue()
}

func builtinShim_DelState(call otto.FunctionCall) otto.Value {
	logger.Warning("DelState, arg#", len(call.ArgumentList))
	stub, err := getShimStub(call)
	if err != nil {
		return otto.NaNValue()
	}
	if err := stub.gostub.DelState(call.Argument(0).String()); err != nil {
		stub.lastErr = err
	}
	return otto.NullValue()
}

func builtinShim_GetStateByRange(call otto.FunctionCall) otto.Value {
	logger.Warning("GetStateByRange, arg#", len(call.ArgumentList))
	stub, err := getShimStub(call)
	if err != nil {
		return otto.NaNValue()
	}
	iter, err := stub.gostub.GetStateByRange(call.Argument(0).String(), call.Argument(1).String())
	if err != nil {
		stub.lastErr = err
		return otto.NullValue()
	}

	v, err := call.Otto.ToValue(&jCCQueryIterator{iter, "", "", nil})
	if err != nil {
		stub.lastErr = err
		logger.Warning(err.Error())
	}
	return v
}

func NewJavascriptCC(ccid *pb.ChaincodeID) (*JavascriptCC, error) {
	fpath := filepath.Join(ccid.Path, ccid.Name+"."+ccid.Version)
	return &JavascriptCC{filepath.Join(fpath, "src"), nil}, nil
}

func (jscc *JavascriptCC) loadSourceCode(stub shim.ChaincodeStubInterface) error {
	fis, err := ioutil.ReadDir(jscc.codepath)
	if err != nil {
		return fmt.Errorf("ReadDir failed %s\n", err)
	}

	sourcecode := make([]byte, 0)
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		if !strings.HasSuffix(fi.Name(), ".js") {
			continue
		}
		name := filepath.Join(jscc.codepath, fi.Name())
		logger.Warningf("javascript file %s", name)

		src, err := ioutil.ReadFile(name)
		if err != nil {
			return err
		}
		sourcecode = append(append(sourcecode, []byte("\n")...), src...)
	}

	if len(sourcecode) == 0 {
		return errors.New("failed to found js source in " + jscc.codepath)
	}

	jscc.vm = otto.New()
	jscc.initShimBuiltin(stub)
	_, err = jscc.vm.Run(sourcecode)
	return err
}

func (jscc *JavascriptCC) initShimBuiltin(stub shim.ChaincodeStubInterface) error {
	v, err := jscc.vm.ToValue(&jCCShim{stub, nil})
	if err != nil {
		fmt.Println("failed toValue of shim")
		return err
	}
	jscc.vm.Set("shim", v)
	jscc.vm.SetObjectFunction("shim", "GetState", builtinShim_GetState)
	jscc.vm.SetObjectFunction("shim", "PutState", builtinShim_PutState)
	jscc.vm.SetObjectFunction("shim", "DelState", builtinShim_DelState)
	jscc.vm.SetObjectFunction("shim", "GetStateByRange", builtinShim_GetStateByRange)
	jscc.vm.SetObjectFunction("shim", "GetFunction", builtinShim_GetFunction)
	jscc.vm.SetObjectFunction("shim", "GetArguments", builtinShim_GetArguments)
	jscc.vm.SetObjectFunction("shim", "GetLastError", builtinShim_GetLastError)
	jscc.vm.SetObjectFunction("shim", "Success", builtinShim_Success)
	jscc.vm.SetObjectFunction("shim", "Error", builtinShim_Error)
	return nil
}

func (jscc *JavascriptCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	err := jscc.loadSourceCode(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	cc, err := jscc.vm.Get("chaincode")
	if err != nil {
		return shim.Error(err.Error())
	}
	if cc.Object() == nil {
		return shim.Error("invalid chaincode obj" + cc.String())
	}
	init_func, err := cc.Object().Get("init")
	if err != nil {
		return shim.Error(err.Error())
	}

	rsp, err := init_func.Call(cc)
	if err != nil {
		logger.Warning("init", err.Error())
	}
	return toShimResponse(rsp)
}

func (jscc *JavascriptCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	if err := jscc.initShimBuiltin(stub); err != nil {
		return shim.Error(err.Error())
	}

	cc, err := jscc.vm.Get("chaincode")
	if err != nil {
		return shim.Error(err.Error())
	}
	if cc.Object() == nil {
		return shim.Error("invalid chaincode obj" + cc.String())
	}
	invoke_func, err := cc.Object().Get("invoke")
	if err != nil {
		return shim.Error(err.Error())
	}

	rsp, err := invoke_func.Call(cc)
	if err != nil {
		logger.Warning("invoke", err.Error())
	}
	return toShimResponse(rsp)
}
