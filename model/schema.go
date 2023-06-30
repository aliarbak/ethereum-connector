package model

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"strings"
)

type Schema struct {
	Id                        string
	Name                      string
	Abi                       []Abi
	Events                    []SchemaEvent
	ContractFilteringEnabled  bool
	FilteredContractAddresses map[int64]map[string]bool
	abi                       *abi.ABI
}

type SchemaEvent struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

type SchemaFilter struct {
	Id              string            `json:"id"`
	ChainId         int64             `json:"chainId"`
	SchemaId        string            `json:"schemaId"`
	ContractAddress string            `json:"contractAddress"`
	Sync            *SchemaFilterSync `json:"sync"`
}

type SchemaFilterSync struct {
	FromBlockNumber    int64 `json:"fromBlockNumber"`
	ToBlockNumber      int64 `json:"toBlockNumber"`
	CurrentBlockNumber int64 `json:"currentBlockNumber"`
	Completed          bool  `json:"completed"`
}

func AddFiltersToSchemas(schemas []Schema, filters []SchemaFilter) {
	for _, filter := range filters {
		for i, schema := range schemas {
			if filter.SchemaId == schema.Id {
				schemas[i].AddContractFilter(filter)
				break
			}
		}
	}
}

func (s *Schema) AddContractFilter(filter SchemaFilter) {
	if !s.ContractFilteringEnabled {
		s.ContractFilteringEnabled = true
		s.FilteredContractAddresses = map[int64]map[string]bool{}
	}

	if s.FilteredContractAddresses[filter.ChainId] == nil {
		s.FilteredContractAddresses[filter.ChainId] = map[string]bool{}
	}

	s.FilteredContractAddresses[filter.ChainId][filter.ContractAddress] = true
}

func (s Schema) IsContractEligible(chainId int64, contractAddress string) bool {
	if !s.ContractFilteringEnabled {
		return true
	}

	if s.FilteredContractAddresses[chainId] == nil {
		return true
	}

	return s.FilteredContractAddresses[chainId][contractAddress]
}

func (s Schema) GetEventAlias(name string) interface{} {
	for _, event := range s.Events {
		if event.Name == name {
			return event.Alias
		}
	}

	return ""
}

func (s *Schema) GetABI() (*abi.ABI, error) {
	if s.abi != nil {
		return s.abi, nil
	}

	abiBytes, err := json.Marshal(s.Abi)
	if err != nil {
		return nil, err
	}

	schemaABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		return nil, err
	}

	s.abi = &schemaABI
	return s.abi, err
}

func (s *Schema) GetLogABI(log *types.Log) (*abi.Event, error) {
	schemaAbi, err := s.GetABI()
	if err != nil {
		return nil, err
	}

	eventAbi, err := schemaAbi.EventByID(log.Topics[0])
	if eventAbi == nil || err != nil {
		return nil, nil
	}

	indexedInputCount := 0
	for _, input := range eventAbi.Inputs {
		if input.Indexed {
			indexedInputCount++
		}
	}

	if indexedInputCount+1 != len(log.Topics) {
		return nil, nil
	}

	return eventAbi, nil
}
