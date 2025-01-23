package hexidecimal

import (
	"encoding/hex"
	"strings"
)

func IsHexBytes(s string) bool {
	_, err := hex.DecodeString(Trim0x(s))
	return err == nil
}

func Trim0x(s string) string {
	return strings.TrimPrefix(s, "0x")
}
