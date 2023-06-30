package model

import "math/big"

type Chain struct {
	Id                        int64    `json:"id"`
	BlockNumber               *big.Int `json:"blockNumber"`
	Paused                    bool     `json:"paused"`
	UpdatedAt                 int64    `json:"updatedAt"`
	RpcUrl                    string   `json:"rpcUrl"`
	BlockNumberWaitPeriodInMS int64    `json:"blockNumberWaitPeriodInMS"`
	BlockNumberReadPeriodInMS int64    `json:"blockNumberReadPeriodInMS"`
	ConcurrentBlockReadLimit  int      `json:"concurrentBlockReadLimit"`
}
