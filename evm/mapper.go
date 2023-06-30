package evm

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"github.com/aliarbak/ethereum-connector/errors"
	"github.com/aliarbak/ethereum-connector/model"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"strconv"
)

func mapLog(chainId int64, log *types.Log, schemas []model.Schema) (map[string]interface{}, error) {
	mapping := map[string]interface{}{}
	rawData, err := log.MarshalJSON()
	if err != nil {
		return mapping, errors.From(err, "rawData marshal failed for txHash: %s, logIndex: %d", log.TxHash.String(), log.Index)
	}

	mapping[model.LogHashLogField] = logHash(log.TxHash.String(), log.Index)
	mapping[model.ContractAddressLogField] = log.Address.String()
	mapping[model.LogIndexLogField] = int64(log.Index)
	mapping[model.RawDataLogField] = string(rawData)

	for _, schema := range schemas {
		eventAbi, err := schema.GetLogABI(log)
		if err != nil {
			return mapping, err
		}

		if eventAbi == nil {
			continue
		}

		if !schema.IsContractEligible(chainId, log.Address.String()) {
			break
		}

		schemaAbi, err := schema.GetABI()
		if err != nil {
			return mapping, err
		}

		err = schemaAbi.UnpackIntoMap(mapping, eventAbi.Name, log.Data)
		if err != nil {
			return mapping, err
		}

		if len(log.Topics) > 1 {
			err = unpackTopicParametersIntoMap(mapping, eventAbi, log.Topics[1:])
			if err != nil {
				return mapping, err
			}
		}

		mapping[model.EventNameLogField] = eventAbi.Name
		mapping[model.EventAliasLogField] = schema.GetEventAlias(eventAbi.Name)
		break
	}

	return mapping, nil
}

func mapRawLog(rawData []byte, schema model.Schema) (map[string]interface{}, error) {
	log := &types.Log{}
	err := log.UnmarshalJSON(rawData)
	if err != nil {
		return nil, errors.From(err, "rawData unmarshall failed")
	}

	mapping := map[string]interface{}{}
	mapping[model.LogHashLogField] = logHash(log.TxHash.String(), log.Index)
	mapping[model.ContractAddressLogField] = log.Address.String()
	mapping[model.LogIndexLogField] = int64(log.Index)
	mapping[model.RawDataLogField] = string(rawData)

	eventAbi, err := schema.GetLogABI(log)
	if err != nil {
		return mapping, err
	}

	if eventAbi == nil {
		return nil, nil
	}

	schemaAbi, err := schema.GetABI()
	if err != nil {
		return mapping, err
	}

	err = schemaAbi.UnpackIntoMap(mapping, eventAbi.Name, log.Data)
	if err != nil {
		return mapping, err
	}

	if len(log.Topics) > 1 {
		err = unpackTopicParametersIntoMap(mapping, eventAbi, log.Topics[1:])
		if err != nil {
			return mapping, err
		}
	}

	mapping[model.EventNameLogField] = eventAbi.Name
	mapping[model.EventAliasLogField] = schema.GetEventAlias(eventAbi.Name)

	return mapping, nil
}

func unpackTopicParametersIntoMap(mapping map[string]interface{}, event *abi.Event, topics []common.Hash) error {
	indexedParamIndex := 0
	for _, param := range event.Inputs {
		if !param.Indexed {
			continue
		}

		val, err := getParamVal(param.Type.String(), topics[indexedParamIndex])
		if err != nil {
			return err
		}

		mapping[param.Name] = val

		indexedParamIndex++
		if indexedParamIndex >= len(topics) {
			break
		}
	}

	return nil
}

func logHash(txHash string, logIndex uint) string {
	s := txHash + strconv.FormatInt(int64(logIndex), 10)
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func getParamVal(paramType string, value common.Hash) (val interface{}, err error) {
	switch paramType {
	case "address":
		val = common.HexToAddress(value.Hex()).String()
	case "uint8":
		val = hashToBigInt(value)
		break
	case "uint16":
		val = hashToBigInt(value)
		break
	case "uint32":
		val = hashToBigInt(value)
		break
	case "uint64":
		val = hashToBigInt(value)
		break
	case "uint128":
		val = hashToBigInt(value)
		break
	case "uint256":
		val = hashToBigInt(value)
		break
	case "int8":
		val = hashToBigInt(value)
		break
	case "int16":
		val = hashToBigInt(value)
		break
	case "int32":
		val = hashToBigInt(value)
		break
	case "int64":
		val = hashToBigInt(value)
		break
	case "int128":
		val = hashToBigInt(value)
		break
	case "int256":
		val = hashToBigInt(value)
		break
	case "bytes":
		val = value.Bytes()
		break
	case "bytes32":
		val = value.Bytes()
		break
	case "bool":
		val, err = hashToBool(value)
		break
	case "boolean":
		val, err = hashToBool(value)
		break
	default:
		err = errors.InvalidInput("unknown param type: %s", paramType)
	}

	return val, err
}

func hashToBigInt(hash common.Hash) *big.Int {
	intVal := new(big.Int)
	intVal.SetBytes(hash.Bytes())
	return intVal
}

func hashToBool(hash common.Hash) (bool, error) {
	var b bool
	err := binary.Read(bytes.NewReader(hash.Bytes()), binary.LittleEndian, &b)
	return b, err
}
