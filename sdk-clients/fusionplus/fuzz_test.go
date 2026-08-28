package fusionplus

import "testing"

// These decoders parse untrusted data (on-chain extensions and API
// responses). The fuzz targets assert they never panic; errors are the
// expected failure mode for malformed input.

func FuzzDecodeEscrowExtension(f *testing.F) {
	f.Add([]byte{})
	f.Add(make([]byte, 160))
	f.Add([]byte{0x00, 0x00, 0x00, 0x01, 0xff})
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = DecodeEscrowExtension(data)
	})
}

func FuzzDecodeExtensionPlus(f *testing.F) {
	f.Add([]byte{})
	f.Add(make([]byte, 64))
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = DecodeExtension(data)
	})
}

func FuzzDecodeSettlementPostInteractionData(f *testing.F) {
	f.Add("")
	f.Add("00")
	f.Add("0000000000000000000000000000000000000000000000000000000000000000")
	f.Fuzz(func(t *testing.T, data string) {
		_, _ = DecodeSettlementPostInteractionData(data)
	})
}
