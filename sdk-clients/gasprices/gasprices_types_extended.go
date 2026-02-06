package gasprices

// GetGasPriceLegacyResponse is used instead of codegen struct to right now as params for API handle
type GetGasPriceLegacyResponse struct {
	Standard string `json:"standard"`
	Fast     string `json:"fast"`
	Instant  string `json:"instant"`
}
