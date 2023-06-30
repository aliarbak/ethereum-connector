package chainmodel

type PatchChainInput struct {
	Paused                    *bool   `json:"paused"`
	RpcUrl                    *string `json:"rpcUrl"`
	BlockNumberWaitPeriodInMS *int64  `json:"blockNumberWaitPeriodInMs"`
	BlockNumberReadPeriodInMS *int64  `json:"blockNumberReadPeriodInMs"`
	ConcurrentBlockReadLimit  *int    `json:"concurrentBlockReadLimit"`
}
