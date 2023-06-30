package chainmodel

type CreateChainInput struct {
	Id                        int64  `json:"id" validate:"required"`
	BlockNumber               *int64 `json:"blockNumber"`
	Paused                    bool   `json:"paused"`
	RpcUrl                    string `json:"rpcUrl" validate:"required"`
	BlockNumberWaitPeriodInMS *int64 `json:"blockNumberWaitPeriodInMs"`
	BlockNumberReadPeriodInMS *int64 `json:"blockNumberReadPeriodInMs"`
	ConcurrentBlockReadLimit  *int   `json:"concurrentBlockReadLimit"`
}
