package traces

type GetBlockTraceByNumberParam int

type GetTxTraceByNumberAndHashParams struct {
	BlockNumber     int    `json:"blockNumber" uri:"blockNumber"`
	TransactionHash string `json:"txHash" uri:"txHash"`
}

type GetTxTraceByNumberAndOffsetParams struct {
	BlockNumber int `json:"blockNumber" uri:"blockNumber"`
	Offset      int `json:"offset" uri:"offset"`
}
