package destination

import (
	"github.com/aliarbak/ethereum-connector/model"
	"math/big"
)

const (
	transactionLogChainIdFieldName           = "chainId"
	transactionLogBlockNumberFieldName       = "blockNumber"
	transactionLogBlockHashFieldName         = "blockHash"
	transactionLogBlockTimeFieldName         = "blockTime"
	transactionLogTxContractAddressFieldName = "txContractAddress"
	transactionLogTxToFieldName              = "txTo"
	transactionLogTxIndexFieldName           = "txIndex"
	transactionLogTxHashFieldName            = "txHash"
	transactionLogTxValueFieldName           = "txValue"
)

type blockMessage struct {
	ChainId int64    `json:"chainId"`
	Number  *big.Int `json:"number"`
	Hash    string   `json:"hash"`
	Time    int64    `json:"time"`
}

type transactionMessage struct {
	ChainId         int64    `json:"chainId"`
	BlockNumber     *big.Int `json:"blockNumber"`
	BlockHash       string   `json:"blockHash"`
	BlockTime       int64    `json:"blockTime"`
	ContractAddress string   `json:"contractAddress"`
	To              string   `json:"to"`
	Hash            string   `json:"hash"`
	Index           int64    `json:"index"`
	Value           *big.Int `json:"value"`
}

type transactionLogMessage map[string]interface{}

func newBlockMessage(block model.Block) blockMessage {
	return blockMessage{
		ChainId: block.ChainId,
		Number:  block.Number,
		Hash:    block.Hash,
		Time:    block.Time,
	}
}

func newTransactionMessage(block model.Block, transaction model.Transaction) transactionMessage {
	return transactionMessage{
		ChainId:         block.ChainId,
		BlockNumber:     block.Number,
		BlockHash:       block.Hash,
		BlockTime:       block.Time,
		ContractAddress: transaction.ContractAddress,
		To:              transaction.To,
		Hash:            transaction.Hash,
		Index:           transaction.Index,
		Value:           transaction.Value,
	}
}

func newTransactionLogMessage(block model.Block, transaction model.Transaction, log map[string]interface{}) transactionLogMessage {
	message := transactionLogMessage{}
	for k, v := range log {
		message[k] = v
	}

	message[transactionLogChainIdFieldName] = block.ChainId
	message[transactionLogBlockNumberFieldName] = block.Number
	message[transactionLogBlockHashFieldName] = block.Hash
	message[transactionLogBlockTimeFieldName] = block.Time
	message[transactionLogTxContractAddressFieldName] = transaction.ContractAddress
	message[transactionLogTxToFieldName] = transaction.To
	message[transactionLogTxHashFieldName] = transaction.Hash
	message[transactionLogTxValueFieldName] = transaction.Value
	message[transactionLogTxIndexFieldName] = transaction.Index
	return message
}
