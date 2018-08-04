package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"testing"
)

type TxidGen struct {
	txid int
}

func (gen *TxidGen) GetTxID() string {
	gen.txid++
	return fmt.Sprintf("%d", gen.txid)
}

var gTxID = &TxidGen{}

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit(gTxID.GetTxID(), args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke(gTxID.GetTxID(), args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", string(args[0]), string(args[1]), "failed", string(res.Message))
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, args [][]byte) []byte {
	res := stub.MockInvoke(gTxID.GetTxID(), args)
	if res.Status != shim.OK {
		fmt.Println("Query failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil || len(res.Payload) < 5 {
		fmt.Println("Query", string(args[0]), string(args[1]), "failed to get result")
		t.FailNow()
	}

	fmt.Println(string(res.Payload))
	return res.Payload
}

func Test_e2e(t *testing.T) {
	scc := new(ExchangeCC)
	stub := shim.NewMockStub("TestExchange", scc)
	checkInit(t, stub, [][]byte{[]byte("init")})

	checkInvoke(t, stub, [][]byte{[]byte("RegisterUser"), []byte("{\"username\":\"username1\"," +
		"\"password\":\"password1\",\"telephone\":\"telephone1\",\"email\":\"email1\",\"date\":\"20170102-010203\"}")})
	checkInvoke(t, stub, [][]byte{[]byte("RegisterUser"), []byte("{\"username\":\"username2\"," +
		"\"password\":\"password2\",\"telephone\":\"telephone2\",\"email\":\"email2\",\"date\":\"20170102-020304\"}")})
	checkQuery(t, stub, [][]byte{[]byte("UserLogin"), []byte("{\"username\":\"username1\"," +
		"\"password\":\"password1\",\"date\":\"20170103-010203\"}")})
	checkQuery(t, stub, [][]byte{[]byte("UserLogin"), []byte("{\"username\":\"username2\"," +
		"\"password\":\"password2\",\"date\":\"20170104-070809\"}")})

	checkQuery(t, stub, [][]byte{[]byte("GetUserInfo"), []byte("{\"username\":\"username1\"}")})
	checkQuery(t, stub, [][]byte{[]byte("GetUserInfo"), []byte("{\"username\":\"username2\"}")})

	checkInvoke(t, stub, [][]byte{[]byte("CreateProduct"), []byte("{\"product_code\":\"product_code1\"," +
		"\"product_name\":\"product_name1\",\"risk_rate\":\"risk_rate1\",\"period\":\"period1\"," +
		"\"exp_annual_rate\":100,\"issuer\":\"issuer1\",\"issue_scale\":10000000," +
		"\"username\":\"username1\"}")})

	checkQuery(t, stub, [][]byte{[]byte("GetAccountBalance"), []byte("{\"username\":\"username1\"}")})

	checkInvoke(t, stub, [][]byte{[]byte("ReserveProduct"), []byte("{\"username\":\"username1\"," +
		"\"product_code\":\"product_code1\",\"reserve_amount\":10000,\"reserve_date\":\"20170523010203\"}")})
	checkInvoke(t, stub, [][]byte{[]byte("ReserveProduct"), []byte("{\"username\":\"username2\"," +
		"\"product_code\":\"product_code1\",\"reserve_amount\":20000,\"reserve_date\":\"20170523010203\"}")})

	checkQuery(t, stub, [][]byte{[]byte("GetProducts"), []byte("")})
	checkQuery(t, stub, [][]byte{[]byte("GetAccountBalance"), []byte("{\"username\":\"username1\"}")})
	checkQuery(t, stub, [][]byte{[]byte("GetAccountPositions"), []byte("{\"username\":\"username1\"}")})

	checkInvoke(t, stub, [][]byte{[]byte("RepurProduct"), []byte("{\"username\":\"username1\"," +
		"\"product_code\":\"product_code1\",\"repur_amount\":5000,\"repur_date\":\"20170523010203\"}")})

	checkQuery(t, stub, [][]byte{[]byte("GetAccountPositions"), []byte("{\"username\":\"username1\"}")})

	checkInvoke(t, stub, [][]byte{[]byte("Settle"), []byte("{\"username\":\"username1\",\"date\":\"20170623010203\"}")})
	checkQuery(t, stub, [][]byte{[]byte("GetAccountPositions"), []byte("{\"username\":\"username1\"}")})
	checkQuery(t, stub, [][]byte{[]byte("GetAccountBalance"), []byte("{\"username\":\"username1\"}")})
	checkQuery(t, stub, [][]byte{[]byte("GetAccountPositions"), []byte("{\"username\":\"username2\"}")})
}
