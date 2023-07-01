package protomessage

import (
	"encoding/json"
	"github.com/aliarbak/ethereum-connector/errors"
	"github.com/aliarbak/ethereum-connector/model"
	"github.com/aliarbak/ethereum-connector/utils"
)

var excludedLogFields = []string{model.EventNameLogField, model.EventAliasLogField, model.ContractAddressLogField, model.LogHashLogField, model.LogIndexLogField, model.RawDataLogField}

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
