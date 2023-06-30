package model

import "math/big"

const (
	EventNameLogField       = "eventName"
	EventAliasLogField      = "eventAlias"
	ContractAddressLogField = "contractAddress"
	LogIndexLogField        = "logIndex"
	LogHashLogField         = "logHash"
	RawDataLogField         = "rawData"
)

type Block struct {
	ChainId      int64
	Number       *big.Int
	Time         int64
	Hash         string
	Transactions []Transaction
}

type Transaction struct {
	ContractAddress string
	To              string
	Index           int64
	Hash            string
	Value           *big.Int
	Data            []byte
	Logs            []map[string]interface{}
}

type RawTransactionLog struct {
	Hash            string
	Index           int64
	ChainId         int64
	BlockNumber     *big.Int
	BlockHash       string
	BlockTime       int64
	TxIndex         int64
	TxHash          string
	TxAddress       string
	TxTo            *string
	TxValue         *big.Int
	ContractAddress string
	RawData         []byte
}
