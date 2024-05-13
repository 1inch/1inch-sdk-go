package history

type HistoryEventsByAddressParams struct {
	Address string `url:"address"  json:"-"`

	Limit int `url:"limit,omitempty" json:"limit,omitempty"`

	TokenAddress string `url:"tokenAddress,omitempty" json:"tokenAddress,omitempty"`

	ChainId int `url:"chainId,omitempty" json:"chainId,omitempty"`

	ToTimestampMs int `url:"toTimestampMs,omitempty" json:"toTimestampMs,omitempty"`

	FromTimestampMs int `url:"fromTimestampMs,omitempty" json:"fromTimestampMs,omitempty"`
}
