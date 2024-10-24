package aggregation

type TokensResponse struct {
	Tokens map[string]TokenInfo `json:"tokens"`
}
