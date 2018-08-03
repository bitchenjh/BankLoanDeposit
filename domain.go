package main

const (
	Bank_Flag_Loan = 1
	Bank_Flag_Repayment = 2

)

type Bank struct {
	BankName	string	`json:"BankName"`
	Amount		int		`json:"Amount"`
	Flag		int		`json:"Flag"`	// 1: 贷款,  2: 还款
	StartTime	string	`json:"StartTime"`
	EndTime		string	`json:"EndTime"`
}

// 查询到最新状态数据的同时可以获取到历史数据
type Account struct {
	CardNo		string	`json:"CardNo"`
	Aname		string	`json:"Aname"`
	Gender		string	`json:"Gender"`
	Mobile		string	`json:"Mobile"`
	Bank		Bank	`json:"Bank"`

	Historys	[]HistoryItem
}

// History Index => levelDB
type HistoryItem struct {
	TxId	string
	Account Account
}

