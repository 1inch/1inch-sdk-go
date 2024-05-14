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

type TransactionTraceResponse struct {
	TransactionTrace TransactionTrace `json:"transactionTrace"`
	Type             string           `json:"type"`
}

type TransactionTrace struct {
	TxHash               string        `json:"txHash"`
	Nonce                string        `json:"nonce"`
	GasPrice             string        `json:"gasPrice"`
	Type                 string        `json:"type"`
	From                 string        `json:"from"`
	To                   string        `json:"to"`
	GasLimit             int           `json:"gasLimit"`
	GasActual            int           `json:"gasActual"`
	GasHex               string        `json:"gasHex"`
	GasUsed              int           `json:"gasUsed"`
	IntrinsicGas         int           `json:"intrinsicGas"`
	GasRefund            int           `json:"gasRefund"`
	Input                string        `json:"input"`
	Calls                []Call        `json:"calls"`
	Logs                 []Log         `json:"logs"`
	Status               string        `json:"status"`
	Storage              []StorageItem `json:"storage"`
	Value                string        `json:"value"`
	MaxFeePerGas         string        `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string        `json:"maxPriorityFeePerGas"`
	Depth                int           `json:"depth"`
}

type Call struct {
	Type         string        `json:"type"`
	From         string        `json:"from"`
	To           string        `json:"to"`
	GasLimit     int           `json:"gasLimit"`
	GasUsed      int           `json:"gasUsed"`
	PrevGasLimit int           `json:"prevGasLimit"`
	Gas          string        `json:"gas"`
	GasCost      int           `json:"gasCost"`
	Input        string        `json:"input"`
	Calls        []Call        `json:"calls"`
	Logs         []Log         `json:"logs"`
	Status       string        `json:"status"`
	Storage      []StorageItem `json:"storage"`
	Success      int           `json:"success"`
	Res          string        `json:"res"`
	Depth        int           `json:"depth"`
	Value        string        `json:"value"`
}

type Log struct {
	Topics   []string `json:"topics"`
	Contract string   `json:"contract"`
	Data     string   `json:"data"`
}

type StorageItem struct {
	Type  string `json:"type"`
	Key   string `json:"key"`
	Value string `json:"value"`
}
