package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"fmt"
	"encoding/json"
	"bytes"
	"time"
	"strconv"
)

type AfcChainCode struct {
}

type Afc struct {
	//Id           string `json:"Id"`
	IdentityCard string `json:"IdentityCard"`
	Desc         string `json:"Desc"`
}

func (t *AfcChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *AfcChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	if fn == "userInfo" {
		return t.userInfo(stub, args)
	} else if fn == "userList" {
		return t.userList(stub, args)
	} else if fn == "userAddUp" {
		return t.userAddUp(stub, args)
	} else if fn == "userHistory" {
		return t.userHistory(stub, args)
	}
	jsonResp := "{\"error_code\":\"" + string(20000) + "\",\"error_msg\":\"" + "function is not set" + "\"}"
	return shim.Error(jsonResp)

}
func (t *AfcChainCode) userInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("start userInfo")
	//afc := Afc{}
	identityCard := args[0]
	userinfo, err := stub.GetState(identityCard)
	if err != nil {
		jsonResp := "{\"error_code\":\"" + string(20000) + "\",\"error_msg\":\"" + "userInfo is empty" + "\"}"
		return shim.Error(jsonResp)
	}
	if userinfo != nil {
		//err = json.Unmarshal(userinfo, &afc)
		//if err != nil{
		//	return shim.Error(err.Error())
		//}
		return shim.Success(userinfo)

	} else {
		jsonResp := "{\"error_code\":\"" + string(20000) + "\",\"error_msg\":\"" + "userInfo is empty" + "\"}"
		return shim.Error(jsonResp)
	}

}

func (t *AfcChainCode) userList(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	fmt.Println("start userList")
	// 获取所有用户的票数
	resultIterator, err := stub.GetStateByRange("", "")
	if err != nil {
		return shim.Error("get user message error！")
	}
	defer resultIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	isWritten := false

	for resultIterator.HasNext() {
		queryResult, err := resultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if isWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString(string(queryResult.Value))
		isWritten = true
	}

	buffer.WriteString("]")

	fmt.Printf("user result：\n%s\n", buffer.String())
	fmt.Println("end getUserVote")
	return shim.Success(buffer.Bytes())
}

func (t *AfcChainCode) userHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("start History")
	identityCard := args[0]
	res, err := stub.GetHistoryForKey(identityCard)

	if err != nil {
		return shim.Error("get user message error！")
	}

	resIt, err := resultIterator(res)
	if err != nil {
		return shim.Error("get user message error！")
	}
	return shim.Success(resIt)
}

func resultIterator(resultIterator shim.HistoryQueryIteratorInterface) ([]byte, error) {
	defer resultIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	isWritten := false

	for resultIterator.HasNext() {
		response, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}

		if isWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		isWritten = true
	}

	buffer.WriteString("]")

	fmt.Printf("result：\n%s\n", buffer.String())
	fmt.Println("end getUserVote")
	return buffer.Bytes(), nil
}

func (t *AfcChainCode) userAddUp(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	IdentityCard := args[0]
	desc := args[1]
	fmt.Println("identityCard:" + IdentityCard)
	fmt.Println("desc")
	afc := Afc{IdentityCard, desc}
	afcJsonBytes, err := json.Marshal(afc)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(IdentityCard, afcJsonBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("saved student!"))
}

func main() {
	err := shim.Start(new(AfcChainCode))
	if err != nil {
		fmt.Println("vote chaincode start err")
	}
}
