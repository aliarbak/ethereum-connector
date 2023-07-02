package destmessage

import (
	"encoding/json"
	"github.com/aliarbak/ethereum-connector/errors"
	"github.com/aliarbak/ethereum-connector/model"
	"github.com/aliarbak/ethereum-connector/utils"
	"math/big"
)

var excludedLogFields = []string{model.EventNameLogField, model.EventAliasLogField, model.ContractAddressLogField, model.LogHashLogField, model.LogIndexLogField, model.RawDataLogField}

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

type BlockMessage struct {
	ChainId int64    `json:"chainId"`
	Number  *big.Int `json:"number"`
	Hash    string   `json:"hash"`
	Time    int64    `json:"time"`
}

type TransactionMessage struct {
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

type TransactionLogMessage map[string]interface{}

func NewBlockMessage(block model.Block) BlockMessage {
	return BlockMessage{
		ChainId: block.ChainId,
		Number:  block.Number,
		Hash:    block.Hash,
		Time:    block.Time,
	}
}

func NewTransactionMessage(block model.Block, transaction model.Transaction) TransactionMessage {
	return TransactionMessage{
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

func NewTransactionLogMessage(block model.Block, transaction model.Transaction, log map[string]interface{}) TransactionLogMessage {
	message := TransactionLogMessage{}
	for k, v := range log {
		if k == model.RawDataLogField {
			continue
		}
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

func NewProtoBlockMessage(block model.Block) Block {
	return Block{
		Chain_Id:   block.ChainId,
		Number:     block.Number.String(),
		Block_Hash: block.Hash,
		Time:       block.Time,
	}
}

func NewProtoTransactionMessage(block model.Block, transaction model.Transaction) Transaction {
	return Transaction{
		Chain_Id:         block.ChainId,
		Block_Number:     block.Number.String(),
		Block_Hash:       block.Hash,
		Block_Time:       block.Time,
		Contract_Address: transaction.ContractAddress,
		Tx_To:            transaction.To,
		Tx_Hash:          transaction.Hash,
		Index:            transaction.Index,
		Value:            transaction.Value.String(),
	}
}

func NewProtoTransactionLogMessage(block model.Block, transaction model.Transaction, log map[string]interface{}) (transactionLog TransactionLog, err error) {
	data := getLogDataMap(log)
	dataJson, err := json.Marshal(data)
	if err != nil {
		return transactionLog, errors.From(err, "transaction log marshall failed, data: %+v", log)
	}

	return TransactionLog{
		Chain_Id:            block.ChainId,
		Block_Number:        block.Number.String(),
		Block_Hash:          block.Hash,
		Block_Time:          block.Time,
		Tx_Contract_Address: transaction.ContractAddress,
		Tx_To:               transaction.To,
		Tx_Hash:             transaction.Hash,
		Tx_Value:            transaction.Value.String(),
		Tx_Index:            transaction.Index,
		Log_Hash:            log[model.LogHashLogField].(string),
		Event_Name:          log[model.EventNameLogField].(string),
		Event_Alias:         log[model.EventAliasLogField].(string),
		Index:               log[model.LogIndexLogField].(int64),
		Contract_Address:    log[model.ContractAddressLogField].(string),
		Data:                string(dataJson),
	}, nil
}

func NewProtoRawTransactionLogMessage(block model.Block, transaction model.Transaction, log map[string]interface{}) (rawTransactionLog RawTransactionLog) {
	return RawTransactionLog{
		Chain_Id:            block.ChainId,
		Block_Number:        block.Number.String(),
		Block_Hash:          block.Hash,
		Block_Time:          block.Time,
		Tx_Contract_Address: transaction.ContractAddress,
		Tx_To:               transaction.To,
		Tx_Hash:             transaction.Hash,
		Tx_Value:            transaction.Value.String(),
		Tx_Index:            transaction.Index,
		Log_Hash:            log[model.LogHashLogField].(string),
		Index:               log[model.LogIndexLogField].(int64),
		Contract_Address:    log[model.ContractAddressLogField].(string),
		Data:                log[model.RawDataLogField].(string),
	}
}

func getLogDataMap(originalMap map[string]interface{}) map[string]interface{} {
	data := map[string]interface{}{}
	for k, v := range originalMap {
		if utils.StringArrayContains(excludedLogFields, k) {
			continue
		}

		data[k] = v
	}

	return data
}
