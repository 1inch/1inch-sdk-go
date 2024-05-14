package traces

type GetBlockTraceByNumberParam int

type GetTxTraceByNumberAndHashParams struct {
	ChainId         int    `json:"chain" uri:"chain"`
	BlockNumber     int    `json:"blockNumber" uri:"blockNumber"`
	TransactionHash string `json:"txHash" uri:"txHash"`
}

type GetTxTraceByNumberAndOffsetParams struct {
	ChainId     int `json:"chain" uri:"chain"`
	BlockNumber int `json:"blockNumber" uri:"blockNumber"`
	Offset      int `json:"offset" uri:"offset"`
}
