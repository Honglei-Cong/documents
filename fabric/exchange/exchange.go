package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type ExchangeCC struct {
}

func (t *ExchangeCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	err := initUserTables(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = initProductTables(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *ExchangeCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	var resp []byte = nil
	var err error

	function, args := stub.GetFunctionAndParameters()
	if function == "RegisterUser" {
		if len(args) != 1 {
			return shim.Error("输入参数不正确")
		}
		err = registerUser(stub, args)

	} else if function == "GetUserInfo" {
		if len(args) != 1 {
			return shim.Error("输入参数不正确")
		}
		resp, err = getUserInfo(stub, args)

	} else if function == "CreateProduct" {
		if len(args) != 1 {
			return shim.Error("输入参数不正确")
		}
		err = createProduct(stub, args)

	} else if function == "ReserveProduct" {
		if len(args) != 1 {
			return shim.Error("输入参数不正确")
		}
		err = reserveProduct(stub, args)

	} else if function == "RepurProduct" {
		if len(args) != 1 {
			return shim.Error("输入参数不正确")
		}
		err = repurProduct(stub, args)

	} else if function == "Settle" {
		err = settle(stub, args)

	} else if function == "UserLogin" {
		resp, err = userLogin(stub, args)

	} else if function == "GetProducts" {
		resp, err = getProducts(stub, args)

	} else if function == "GetAccountPositions" {
		if len(args) != 1 {
			return shim.Error("输入参数不正确")
		}
		resp, err = getAccountPositions(stub, args)

	} else if function == "GetAccountBalance" {
		if len(args) != 1 {
			return shim.Error("输入参数不正确")
		}
		resp, err = getAccountBalance(stub, args)

	} else {
		err = errors.New("Unsupported invoke: " + function)
	}

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(resp)
}

func main() {
	err := shim.Start(new(ExchangeCC))
	if err != nil {
		fmt.Printf("Error starting TianJS chaincode: %s", err)
	}
}
