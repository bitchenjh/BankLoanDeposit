package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"strconv"
)

func putAccount(stub shim.ChaincodeStubInterface, account Account) bool {
	accBytes, err := json.Marshal(account)
	if err != nil {
		return false
	}

	err = stub.PutState(account.CardNo, accBytes)
	if err != nil {
		return false
	}
	return true
}

func getAccount(stub shim.ChaincodeStubInterface, cardNo string) (Account, bool) {
	var account Account
	b, err := stub.GetState(cardNo)
	if err != nil {
		return account, false
	}
	if b == nil {
		return account, false
	}
	json.Unmarshal(b, &account)
	return account, true
}

// 贷款
// '{"Args":["loan", "身份证号码", "银行名称", "贷款金额"]}'
func loan(stub shim.ChaincodeStubInterface, args []string) peer.Response  {
	if len(args) != 3 {
		return shim.Error("指定的所需参数错误")
	}
	am, err := strconv.Atoi(args[2])
	if err != nil{
		return shim.Error("指定的贷款金额错误")
	}
	bank := Bank{
		BankName:args[1],
		Amount:am,
		Flag:Bank_Flag_Loan,
		StartTime:"2011-09-10",
		EndTime:"2012-07-09",
	}
	acc := Account{
		CardNo:args[0],
		Aname:"jack",
		Gender:"男",
		Mobile:"7633276",
		Bank:bank,
		// 历史数据无需考虑, 由Fabric自行维护
	}

	bl := putAccount(stub, acc)
	if !bl{
		return shim.Error("保存贷款数据失败")
	}
	return shim.Success([]byte("保存贷款数据成功"))
}

// 还款
// '{"Args":["repayment", "身份证号码", "银行名称", "还款金额"]}'
func repayment(stub shim.ChaincodeStubInterface, args []string) peer.Response{
	if len(args) != 3 {
		return shim.Error("指定的所需参数错误")
	}
	am, err := strconv.Atoi(args[2])
	if err != nil{
		return shim.Error("指定的还款金额错误")
	}
	bank := Bank{
		BankName:args[1],
		Amount:am,
		Flag:Bank_Flag_Repayment,
		StartTime:"2011-09-10",
		EndTime:"2012-07-09",
	}
	acc := Account{
		CardNo:args[0],
		Aname:"jack",
		Gender:"男",
		Mobile:"7633276",
		Bank:bank,
		// 历史数据无需考虑, 由Fabric自行维护
	}

	bl := putAccount(stub, acc)
	if !bl{
		return shim.Error("保存还款数据失败")
	}
	return shim.Success([]byte("保存还款数据成功"))
}

// 根据身份证号码查询最新的状态数据及相应的历史数据
// '{"Args":["queryAccountByCardNo","身份证号码"]}'
func queryAccountByCardNo(stub shim.ChaincodeStubInterface, args []string) peer.Response  {
	if len(args) != 1 {
		return shim.Error("指定的所需参数错误")
	}

	account, bl := getAccount(stub, args[0])
	if !bl {
		return shim.Error("根据指定的身份证号码查询对应数据时错误")
	}

	// 查询历史数据: GetHistoryForKey
	accIterator, err := stub.GetHistoryForKey(account.CardNo)
	if err != nil {
		return shim.Error("获取历史数据时发生错误")
	}
	defer accIterator.Close()

	var historys []HistoryItem
	var acc Account
	for accIterator.HasNext() {
		hisData, err := accIterator.Next()
		if err != nil {
			return shim.Error("处理迭代器数据时发生错误")
		}

		var hisItem HistoryItem
		// 获取当前交易的交易编号
		hisItem.TxId = hisData.TxId
		// 获取当前交易的历史信息
		err = json.Unmarshal(hisData.Value, &acc)
		if err != nil {
			return shim.Error("反序列化历史数据时发生错误")
		}
		if hisData.Value == nil {
			// 如果没有对应的历史记录, 则赋值为一个空对象
			var empty Account
			hisItem.Account = empty
		}else {
			hisItem.Account = acc
		}

		historys = append(historys, hisItem)

	}

	account.Historys = historys

	accountBytes, err := json.Marshal(account)
	if err != nil {
		return shim.Error("序列化数据时发生错误")
	}

	return shim.Success(accountBytes)
}