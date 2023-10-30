package helpers

import "regexp"

// IsEthereumAddress checks if the provided string is a valid Ethereum address.
func IsEthereumAddress(address string) bool {
	// Ethereum address starts with '0x' followed by 40 hexadecimal characters.
	re := regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
	return re.MatchString(address)
}
