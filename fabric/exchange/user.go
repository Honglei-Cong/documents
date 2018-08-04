package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
)

type User struct {
	UserID        string `json:"user_id"`
	Username      string `json:"username"`
	Org           string `json:"org"`
	Password      string `json:"password"`
	Telephone     string `json:"telephone"`
	Email         string `json:"email"`
	RegisterDate  string `json:"register_date"`
	LastLoginTime string `json:"last_login_time"`
}

func (user *User) GetID() string   { return user.UserID }
func (user *User) SetID(id string) { user.UserID = id }
func (user *User) IsSeqID() bool   { return false }

var userTable = ExchangeTableInfo{
	name:    "UserTable",
	obj:     &User{},
	subkeys: []string{"username"},
}

var UserTables = [...]ExchangeTableInfo{
	userTable,
}

func initUserTables(stub shim.ChaincodeStubInterface) error {
	for _, table := range UserTables {
		if !isTableExisted(stub, table) {
			err := CreateObjectType(stub, table.obj, table.subkeys)
			if err != nil {
				return fmt.Errorf("表初始化失败 %s: %s", table.name, err)
			}
		}
	}
	return nil
}

type RegisterUserParam struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
	Date      string `json:"date"`
}

func registerUser(stub shim.ChaincodeStubInterface, params []string) error {
	param := &RegisterUserParam{}
	err := json.Unmarshal([]byte(params[0]), param)
	if err != nil {
		return errors.New("用户注册参数不正确. " + err.Error())
	}

	if len(param.Username) == 0 ||
		len(param.Telephone) == 0 ||
		len(param.Email) == 0 {
		return errors.New("用户注册参数不完整.")
	}

	_, err = getUserByName(stub, param.Username)
	if err == nil {
		return errors.New("用户 " + param.Username + " 已经注册.")
	}

	creatorBytes, err := stub.GetCreator()
	if err != nil {
		return errors.New("用户 " + param.Username + " 请求为无效请求.")
	}
	sid := &msp.SerializedIdentity{}
	if err := proto.Unmarshal(creatorBytes, sid); err != nil {
		return fmt.Errorf("用户 %s 请求错误 %s.", param.Username, err)
	}
	if len(sid.Mspid) == 0 {
		return errors.New("用户 " + param.Username + " 请求包含无效信息.")
	}

	user, err := createUser(stub, param.Username, sid.Mspid, param.Password, param.Telephone, param.Email, param.Date)
	if err != nil {
		return fmt.Errorf("创建用户失败 %s.", err)
	}

	if err = createUserBalanceAccount(stub, user.Username, 100000000); err != nil {
		return fmt.Errorf("创建用户资金账户失败: %s", err)
	}

	return nil
}

type getUserInfoParam struct {
	Username string `json:"username"`
}

type getUserInfoResp struct {
	UserID        string `json:"user_id"`
	Username      string `json:"username"`
	Org           string `json:"org"`
	Telephone     string `json:"telephone"`
	Email         string `json:"email"`
	RegisterDate  string `json:"register_date"`
	LastLoginTime string `json:"last_login_time"`
}

func getUserInfo(stub shim.ChaincodeStubInterface, params []string) ([]byte, error) {
	param := &getUserInfoParam{}
	err := json.Unmarshal([]byte(params[0]), param)
	if err != nil {
		return nil, errors.New("用户查询参数不正确. " + err.Error())
	}
	if len(param.Username) == 0 {
		return nil, errors.New("用户查询参数不完整.")
	}

	creatorBytes, err := stub.GetCreator()
	if err != nil {
		return nil, errors.New("用户 " + param.Username + " 请求为无效请求.")
	}
	sid := &msp.SerializedIdentity{}
	if err := proto.Unmarshal(creatorBytes, sid); err != nil {
		return nil, fmt.Errorf("用户 %s 请求错误 %s.", param.Username, err)
	}
	if len(sid.Mspid) == 0 {
		return nil, errors.New("用户 " + param.Username + " 请求包含无效信息.")
	}

	user, err := getUserByName(stub, param.Username)
	if err != nil || user.Org != sid.Mspid {
		return nil, errors.New("用户 " + param.Username + " 查询证书错误.")
	}

	return json.Marshal(getUserInfoResp{
		UserID:        user.UserID,
		Username:      user.Username,
		Org:           user.Org,
		Telephone:     user.Telephone,
		Email:         user.Email,
		RegisterDate:  user.RegisterDate,
		LastLoginTime: user.LastLoginTime,
	})
}

type LoginParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Date     string `json:"date"`
}

type LoginResp struct {
	Username string `json:"username"`
}

func userLogin(stub shim.ChaincodeStubInterface, params []string) ([]byte, error) {
	param := &LoginParam{}
	err := json.Unmarshal([]byte(params[0]), param)
	if err != nil {
		return nil, errors.New("用户登录参数不正确. " + err.Error())
	}

	if len(param.Username) == 0 ||
		len(param.Password) == 0 {
		return nil, errors.New("用户登录参数不完整.")
	}

	creatorBytes, err := stub.GetCreator()
	if err != nil {
		return nil, errors.New("用户 " + param.Username + " 请求为无效请求.")
	}
	sid := &msp.SerializedIdentity{}
	if err := proto.Unmarshal(creatorBytes, sid); err != nil {
		return nil, fmt.Errorf("用户 %s 请求错误 %s.", param.Username, err)
	}
	if len(sid.Mspid) == 0 {
		return nil, errors.New("用户 " + param.Username + " 请求包含无效信息.")
	}

	user, err := getUserByName(stub, param.Username)
	if err != nil || user.Password != param.Password || user.Org != sid.Mspid {
		return nil, errors.New("用户 " + param.Username + " 不存在或密码错误.")
	}

	user.LastLoginTime = param.Date
	if err := PutObject(stub, user); err != nil {
		return nil, errors.New("用户 " + param.Username + " 登录更新失败.")
	}

	return json.Marshal(LoginResp{
		Username: user.Username,
	})
}

func getUserByName(stub shim.ChaincodeStubInterface, username string) (*User, error) {
	user := &User{}
	keys := map[string]string{"username": username}
	err := GetObject(stub, user, keys)
	return user, err
}

func createUser(stub shim.ChaincodeStubInterface, username string, org, password, telephone, email, date string) (*User, error) {
	user := &User{
		Username:     username,
		Org:          org,
		Password:     password,
		Telephone:    telephone,
		Email:        email,
		RegisterDate: date,
	}

	err := PutObject(stub, user)
	return user, err
}
