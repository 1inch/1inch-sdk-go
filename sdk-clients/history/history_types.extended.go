package history

type EventsByAddressParams struct {
	Address string `url:"address"  json:"-"`

	Limit int `url:"limit,omitempty" json:"limit,omitempty"`

	TokenAddress string `url:"tokenAddress,omitempty" json:"tokenAddress,omitempty"`

	ChainId int `url:"chainId,omitempty" json:"chainId,omitempty"`

	ToTimestampMs int `url:"toTimestampMs,omitempty" json:"toTimestampMs,omitempty"`

	FromTimestampMs int `url:"fromTimestampMs,omitempty" json:"fromTimestampMs,omitempty"`
}

type EventsByAddressResponse struct {
	Items        []Item `json:"items"`
	CacheCounter int    `json:"cache_counter"`
}

type Item struct {
	TimeMs                  int64   `json:"timeMs"`
	Address                 string  `json:"address"`
	Type                    int     `json:"type"`
	Rating                  string  `json:"rating"`
	Details                 Details `json:"details"`
	ID                      string  `json:"id"`
	EventOrderInTransaction int     `json:"eventOrderInTransaction"`
}

type Details struct {
	TxHash       string        `json:"txHash"`
	ChainID      int           `json:"chainId"`
	BlockNumber  int           `json:"blockNumber"`
	BlockTimeSec int64         `json:"blockTimeSec"`
	Status       string        `json:"status"`
	Type         string        `json:"type"`
	TokenActions []TokenAction `json:"tokenActions"`
	FromAddress  string        `json:"fromAddress"`
	ToAddress    string        `json:"toAddress"`
	OrderInBlock int           `json:"orderInBlock"`
	Nonce        int           `json:"nonce"`
	FeeInWei     string        `json:"feeInWei"`
}

type TokenAction struct {
	Address     string `json:"address"`
	Standard    string `json:"standard"`
	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress"`
	Amount      string `json:"amount"`
	Direction   string `json:"direction"`
}
