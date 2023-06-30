package datasource

import (
	"github.com/aliarbak/ethereum-connector/model"
	"math/big"
)

type ChainDocument struct {
	Id                        int64    `json:"id"`
	BlockNumber               *big.Int `json:"blockNumber"`
	Paused                    bool     `json:"paused"`
	RpcUrl                    string   `json:"rpcUrl"`
	BlockNumberWaitPeriodInMS int64    `json:"blockNumberWaitPeriodInMS"`
	BlockNumberReadPeriodInMS int64    `json:"blockNumberReadPeriodInMS"`
	ConcurrentBlockReadLimit  int      `json:"concurrentBlockReadLimit"`
	UpdatedAt                 int64    `json:"updatedAt"`
}

type SchemaDocument struct {
	Id     string                `json:"id"`
	Name   string                `json:"name"`
	Abi    []model.Abi           `json:"abi"`
	Events []SchemaEventDocument `json:"events"`
}

type SchemaEventDocument struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
}
