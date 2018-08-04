package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type ExchangeTableInfo struct {
	name    string
	obj     KVObject
	subkeys []string
}

type Product struct {
	ProductID        string `json:"product_id"`
	ProductCode      string `json:"product_code"`
	ProductName      string `json:"product_name"`
	RiskRate         string `json:"risk_rate"`
	Period           string `json:"period"`
	ExpAnnualRate    uint64 `json:"exp_annual_rate"`
	CurrencyCategory string `json:"currency_category"`
	Issuer           string `json:"issuer"`
	IssueScale       uint64 `json:"issue_scale"`

	PurchasedScale uint64 `json:"purchased_scale"`

	CreatedDate string `json:"created_date"`
	CreatorName string `json:"creator_name"`
}

func (product *Product) GetID() string   { return product.ProductID }
func (product *Product) SetID(id string) { product.ProductID = id }
func (Product *Product) IsSeqID() bool   { return false }

var productTable = ExchangeTableInfo{
	name:    "ProductTable",
	obj:     &Product{},
	subkeys: []string{"product_code"},
}

var ProductTables = [...]ExchangeTableInfo{
	productTable,
	userProductPositionTable,
	userBalanceTable,
}

func isTableExisted(stub shim.ChaincodeStubInterface, table ExchangeTableInfo) bool {
	return IsObjectTypeExist(stub, table.obj)
}

func initProductTables(stub shim.ChaincodeStubInterface) error {

	for _, table := range ProductTables {
		if !isTableExisted(stub, table) {
			err := CreateObjectType(stub, table.obj, table.subkeys)
			if err != nil {
				return fmt.Errorf("Failed to init %s: %s", table.name, err)
			}
		}
	}

	return nil
}

type CreateProductParam struct {
	ProductCode      string `json:"product_code"`
	ProductName      string `json:"product_name"`
	RiskRate         string `json:"risk_rate"`
	Period           string `json:"period"`
	ExpAnnualRate    uint64 `json:"exp_annual_rate"`
	CurrencyCategory string `json:"currency_category"`
	Issuer           string `json:"issuer"`
	IssueScale       uint64 `json:"issue_scale"`
	Date             string `json:"date"`
	Username         string `json:"username"`
}

func createProduct(stub shim.ChaincodeStubInterface, params []string) error {
	param := &CreateProductParam{}
	err := json.Unmarshal([]byte(params[0]), param)
	if err != nil {
		return errors.New("创建产品解析输入参数失败: " + err.Error())
	}

	_, err = getUserByName(stub, param.Username)
	if err != nil {
		return fmt.Errorf("用户 %s 不存在", param.Username)
	}

	product := &Product{}
	product.ProductCode = param.ProductCode
	product.ProductName = param.ProductName
	product.RiskRate = param.RiskRate
	product.Period = param.Period
	product.ExpAnnualRate = param.ExpAnnualRate
	product.CurrencyCategory = param.CurrencyCategory
	product.Issuer = param.Issuer
	product.IssueScale = param.IssueScale

	product.CreatedDate = param.Date
	product.CreatorName = param.Username
	product.PurchasedScale = 0
	err = PutObject(stub, product)
	if err != nil {
		return fmt.Errorf("发布产品失败: %s", err)
	}

	return nil
}

func getProductByCode(stub shim.ChaincodeStubInterface, productCode string) (*Product, error) {
	productInfo := &Product{}
	keys := map[string]string{"product_code": productCode}
	err := GetObject(stub, productInfo, keys)
	if err != nil {
		return nil, err
	}

	return productInfo, nil
}

type ReserveProductParam struct {
	Username      string `json:"username"`
	ProductCode   string `json:"product_code"`
	ReserveAmount uint64 `json:"reserve_amount"`
	ReserveDate   string `json:"reserve_date"`
}

func reserveProduct(stub shim.ChaincodeStubInterface, params []string) error {
	param := &ReserveProductParam{}
	err := json.Unmarshal([]byte(params[0]), param)
	if err != nil {
		return errors.New("预约产品失败: " + err.Error())
	}

	if len(param.ProductCode) < 1 {
		return fmt.Errorf("无效的产品:%s", param.ProductCode)
	}

	// check if username has been in table,
	_, err = getUserByName(stub, param.Username)
	if err != nil {
		return fmt.Errorf("用户不存在或信息错误 %s ", err)
	}

	product, err := getProductByCode(stub, param.ProductCode)
	if err != nil {
		return fmt.Errorf("产品不存在或信息错误 %s ", err)
	}

	if product.IssueScale < product.PurchasedScale+param.ReserveAmount {
		return fmt.Errorf("预约份额超出产品 %s 剩余份额", product.ProductName)
	}

	balance, err := getUserBalance(stub, param.Username)
	if err != nil {
		return err
	}

	if balance.Balance < param.ReserveAmount {
		return fmt.Errorf("用户 %s 余额不足", param.Username)
	}

	userPosition := &UserProductionPositionInfo{
		Username:       param.Username,
		ProductCode:    param.ProductCode,
		Position:       param.ReserveAmount,
		PurTxnDate:     param.ReserveDate,
		LastSettleDate: param.ReserveDate,
	}

	err = PutObject(stub, userPosition)
	if err != nil {
		return fmt.Errorf("预约产品保存失败: %s", err)
	}

	product.PurchasedScale += param.ReserveAmount
	err = PutObject(stub, product)
	if err != nil {
		return fmt.Errorf("更新产品可购买份额失败: %s", err)
	}

	balance.Balance -= param.ReserveAmount
	err = PutObject(stub, balance)
	if err != nil {
		return fmt.Errorf("更新用户账户失败: %s", err)
	}

	return nil
}

type RepurProductParam struct {
	Username    string `json:"username"`
	ProductCode string `json:"product_code"`
	RepurAmount uint64 `json:"repur_amount"`
	RepurDate   string `json:"repur_date"`
}

func repurProduct(stub shim.ChaincodeStubInterface, params []string) error {
	param := &RepurProductParam{}
	err := json.Unmarshal([]byte(params[0]), param)
	if err != nil {
		return errors.New("赎回产品解析输入参数失败: " + err.Error())
	}
	if param.RepurAmount == 0 {
		return nil
	}

	// check if username has been in table,
	user, err := getUserByName(stub, param.Username)
	if err != nil {
		return fmt.Errorf("用户 %s 没有权限赎回", param.Username)
	}
	if user == nil {
		return errors.New("用户 " + param.Username + " 没有注册")
	}

	product, err := getProductByCode(stub, param.ProductCode)
	if err != nil {
		return fmt.Errorf("产品不存在或信息错误 %s ", err)
	}
	if product.PurchasedScale < param.RepurAmount {
		return errors.New("赎回份额超出可赎回持仓.")
	}

	userPosition := &UserProductionPositionInfo{}
	userPositionList := make([]*UserProductionPositionInfo, 0)

	keys := map[string]string{"username": param.Username, "product_code": param.ProductCode}
	ite, err := GetObjectsIte(stub, userPosition, keys)
	if err != nil {
		return fmt.Errorf("查询用户持仓失败 (%s).", err.Error())
	}

	leftAmount := param.RepurAmount
	for ite.HasNext() {
		kv, err := ite.Next()
		if err != nil {
			return fmt.Errorf("查询用户持仓迭代失败: %s", err)
		}

		up := &UserProductionPositionInfo{}
		if err = json.Unmarshal(kv.Value, up); err != nil {
			return fmt.Errorf("查询用户持仓解析成json失败: %s", err)
		}

		userPositionList = append(userPositionList, up)
		if up.Position > leftAmount {
			leftAmount = 0
			break
		}
		leftAmount -= up.Position
	}
	ite.Close()

	if leftAmount > 0 {
		return errors.New("没有足够的持仓.")
	}

	leftAmount = param.RepurAmount
	interest := uint64(0)
	for _, up := range userPositionList {
		if up.Position > leftAmount {
			up.Position -= leftAmount
			if err := UpdateObject(stub, up); err != nil {
				return fmt.Errorf("更新用户持仓失败: %s", err)
			}
			break
		} else {
			leftAmount -= up.Position
			interest += up.Interest
			if err := DelObject(stub, up); err != nil {
				return fmt.Errorf("更新用户持仓记录失败: %s", err)
			}
		}
	}

	product.PurchasedScale -= param.RepurAmount
	err = PutObject(stub, product)
	if err != nil {
		return fmt.Errorf("更新产品购买份额失败: %s", err)
	}

	balance, err := getUserBalance(stub, param.Username)
	if err != nil {
		return err
	}
	balance.Balance += param.RepurAmount + interest
	err = PutObject(stub, balance)
	if err != nil {
		return fmt.Errorf("更新用户余额失败: %s", err)
	}

	return nil
}

func getProducts(stub shim.ChaincodeStubInterface, params []string) ([]byte, error) {

	maxCount := 1000
	products := make([]*Product, 0, maxCount)

	ite, err := GetObjectsIte(stub, &Product{}, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("查询产品获取ite失败: %s", err)
	}

	for ite.HasNext() {
		kv, err := ite.Next()
		if err != nil {
			return nil, fmt.Errorf("查询产品获取下一条记录失败: %s", err)
		}

		p := &Product{}
		if err = json.Unmarshal(kv.Value, p); err != nil {
			return nil, fmt.Errorf("查询产品解析json失败: %s", err)
		}
		products = append(products, p)
	}

	return json.Marshal(products)
}
