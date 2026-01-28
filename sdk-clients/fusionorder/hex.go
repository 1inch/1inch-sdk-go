package fusionorder

// Prefix0x ensures a hex string has the 0x prefix
func Prefix0x(value string) string {
	if len(value) >= 2 && value[:2] == "0x" {
		return value
	}
	return "0x" + value
}
