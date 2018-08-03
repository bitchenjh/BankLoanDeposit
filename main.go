package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

type TraceChaincode struct {

}


// 实例化/升级链码时被调用
// 在实例化链码时生成相应的测试数据
func (t *TraceChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response  {
	// 初始化测试数据
	return initTest(stub)
	//return shim.Success(nil)
}

/**
 * loan: 贷款
 * repayemnt: 还款
 * queryAccountByCardNo: 根据账户身份证号码查询相应数据(包含历史数据)
 */
func (t *TraceChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response  {
	fun, args := stub.GetFunctionAndParameters()
	if fun == "loan" {
		return loan(stub, args)
	}else if fun == "repayment" {
		return repayment(stub, args)
	}else if fun == "queryAccountByCardNo" {
		return queryAccountByCardNo(stub, args)
	}
	return shim.Error("指定的操作为非法操作, 无法完成")
}

// 创建初始测试数据
func initTest(stub shim.ChaincodeStubInterface) peer.Response {
	bank := Bank{"ccic", 60000, 1, "2010-01-10", "2010-07-09"}
	bank2 := Bank{"ccic", 10000, 2, "2010-02-10", "2010-07-09"}

	acc := Account{"4243645", "jack", "男", "7633276", bank, nil}
	acc2 := Account{"4243645", "jack", "男", "7633276", bank2, nil}

	// 将对象进行存储或在网络上进行传输必须要将其序列化
	accBytes, err := json.Marshal(acc)
	if err != nil {
		return shim.Error("序列化acc对象时发生错误")
	}
	err = stub.PutState(acc.CardNo, accBytes)
	if err != nil {
		return shim.Error("保存acc对象时发生错误")
	}

	accBytes2, err := json.Marshal(acc2)
	if err != nil {
		return shim.Error("序列化acc2对象时发生错误")
	}
	err = stub.PutState(acc2.CardNo, accBytes2)
	if err != nil {
		return shim.Error("保存acc2对象时发生错误")
	}
	return shim.Success(nil)



}



// 将链码中业务方法与链码源文件分离
func main() {
	err := shim.Start(new(TraceChaincode))
	if err != nil {
		fmt.Printf("启动链码时发生错误: %s", err)
	}
}


