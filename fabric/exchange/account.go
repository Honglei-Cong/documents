package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type UserProductionPositionInfo struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	ProductCode    string `json:"product_code"`
	Position       uint64 `json:"position"`
	Interest       uint64 `json:"interest"`
	PurTxnDate     string `json:"pur_txn_date"`
	LastSettleDate string `json:"last_settle_date"`
}

func (userPosition *UserProductionPositionInfo) GetID() string   { return userPosition.ID }
func (userPosition *UserProductionPositionInfo) SetID(id string) { userPosition.ID = id }
func (userPosition *UserProductionPositionInfo) IsSeqID() bool   { return false }

var userProductPositionTable = ExchangeTableInfo{
	name:    "UserProductPositionTable",
	obj:     &UserProductionPositionInfo{},
	subkeys: []string{"username", "product_code"},
}

type UserBalanceInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Balance  uint64 `json:"balance"`
}

func (userBalance *UserBalanceInfo) GetID() string   { return userBalance.ID }
func (userBalance *UserBalanceInfo) SetID(id string) { userBalance.ID = id }
func (userBalance *UserBalanceInfo) IsSeqID() bool   { return false }

var userBalanceTable = ExchangeTableInfo{
	name:    "UserBalanceTable",
	obj:     &UserBalanceInfo{},
	subkeys: []string{"username"},
}

type GetAccountPositionParam struct {
	Username string `json:"username"`
}

type GetAccountPositionResp struct {
	Username       string `json:"username"`
	ProductCode    string `json:"product_code"`
	ProductName    string `json:"product_name"`
	Position       uint64 `json:"position"`
	Interest       uint64 `json:"interest"`
	PurDate        string `json:"pur_date"`
	LastSettleDate string `json:"last_settle_date"`
}

func getAccountPositions(stub shim.ChaincodeStubInterface, params []string) ([]byte, error) {
	param := &GetAccountPositionParam{}
	err := json.Unmarshal([]byte(params[0]), param)
	if err != nil {
		return nil, errors.New("查询账户持仓解析输入参数失败: " + err.Error())
	}

	positions := make([]*GetAccountPositionResp, 0, 400)
	keys := map[string]string{"username": param.Username}
	ite, err := GetObjectsIte(stub, &UserProductionPositionInfo{}, keys)
	if err != nil {
		return nil, fmt.Errorf("查询过程中获取ite失败: %s", err)
	}

	for ite.HasNext() {
		kv, err := ite.Next()
		if err != nil {
			return nil, fmt.Errorf("查询过程中获取下一条纪录失败: %s", err)
		}

		up := &UserProductionPositionInfo{}
		if err = json.Unmarshal(kv.Value, up); err != nil {
			return nil, fmt.Errorf("查询过程中解析json失败 %s", err)
		}
		positions = append(positions, &GetAccountPositionResp{
			Username:       up.Username,
			ProductCode:    up.ProductCode,
			Position:       up.Position,
			Interest:       up.Interest,
			PurDate:        up.PurTxnDate,
			LastSettleDate: up.LastSettleDate,
		})
	}

	for _, p := range positions {
		productInfo, _ := getProductByCode(stub, p.ProductCode)
		p.ProductName = productInfo.ProductName
	}

	return json.Marshal(positions)
}

type SettleParam struct {
	Username string `json:"username"`
	Date     string `json:"date"`
}

func settle(stub shim.ChaincodeStubInterface, params []string) error {
	param := &SettleParam{}
	err := json.Unmarshal([]byte(params[0]), param)
	if err != nil {
		return errors.New("赎回产品解析输入参数失败: " + err.Error())
	}

	// check if username has been in table,
	user, err := getUserByName(stub, param.Username)
	if err != nil {
		return fmt.Errorf("用户 %s 没有权限结算", param.Username)
	}
	if user == nil {
		return errors.New("用户 " + param.Username + " 没有注册")
	}

	userPosition := &UserProductionPositionInfo{}
	ite, err := GetObjectsIte(stub, userPosition, map[string]string{})
	if err != nil {
		return fmt.Errorf("查询用户持仓失败 (%s).", err.Error())
	}

	for ite.HasNext() {
		kv, err := ite.Next()
		if err != nil {
			return fmt.Errorf("查询用户持仓迭代失败: %s", err)
		}

		up := &UserProductionPositionInfo{}
		if err = json.Unmarshal(kv.Value, up); err != nil {
			return fmt.Errorf("查询用户持仓解析成json失败: %s", err)
		}

		up.Interest += 100
		up.LastSettleDate = param.Date
		if err := PutObject(stub, up); err != nil {
			return fmt.Errorf("更新用户利息失败: %s", err)
		}

	}
	ite.Close()

	return nil
}

type GetAccountBalanceParam struct {
	Username string `json:"username"`
}

type GetAccountBalanceResp struct {
	Username string `json:"username"`
	Balance  uint64 `json:"balance"`
}

func getAccountBalance(stub shim.ChaincodeStubInterface, params []string) ([]byte, error) {
	param := &GetAccountBalanceParam{}
	err := json.Unmarshal([]byte(params[0]), param)
	if err != nil {
		return nil, errors.New("查询账户余额解析输入参数失败: " + err.Error())
	}

	balance, err := getUserBalance(stub, param.Username)
	if err != nil {
		return nil, err
	}

	return json.Marshal(GetAccountBalanceResp{
		Username: param.Username,
		Balance:  balance.Balance,
	})
}

func createUserBalanceAccount(stub shim.ChaincodeStubInterface, username string, balance uint64) error {
	userBalance := &UserBalanceInfo{
		Username: username,
		Balance:  balance,
	}
	return PutObject(stub, userBalance)
}

func getUserBalance(stub shim.ChaincodeStubInterface, username string) (*UserBalanceInfo, error) {
	balance := &UserBalanceInfo{}
	keys := map[string]string{"username": username}
	if err := GetObject(stub, balance, keys); err != nil {
		return nil, fmt.Errorf("查询用户 %s 账户失败: %s", username, err)
	}

	return balance, nil
}
