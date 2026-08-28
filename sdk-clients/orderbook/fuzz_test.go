package orderbook

import "testing"

// These decoders parse untrusted data (on-chain orders and API responses).
// The fuzz targets assert they never panic; errors are the expected failure
// mode for malformed input.

func FuzzDecodeMakerTraits(f *testing.F) {
	f.Add("0x4000000000000000000000000000000000000000000000000000000000000000")
	f.Add("62419058400000000000000000000000000000000000000000000000000000")
	f.Add("")
	f.Add("nothex")
	f.Fuzz(func(t *testing.T, encoded string) {
		_, _ = DecodeMakerTraits(encoded)
	})
}

func FuzzDecodeExtension(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0x00, 0x00, 0x00, 0x01, 0xff})
	f.Add(make([]byte, 32))
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = Decode(data)
	})
}
