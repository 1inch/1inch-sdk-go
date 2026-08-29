package fusionorder

import "testing"

// These decoders parse untrusted data (API responses and on-chain calldata).
// The fuzz targets assert they never panic; errors are the expected failure
// mode for malformed input.

func FuzzDecodeAuctionDetails(f *testing.F) {
	f.Add("0x00000000000000000000000000000000000000000000000000")
	f.Add("64b0c00000b400000000c800012c01")
	f.Add("")
	f.Add("zz")
	f.Fuzz(func(t *testing.T, data string) {
		_, _ = DecodeAuctionDetails(data)
	})
}

func FuzzDecodeLegacyAuctionDetails(f *testing.F) {
	f.Add("64b0c00000b400000000c800012c")
	f.Add("")
	f.Fuzz(func(t *testing.T, data string) {
		_, _ = DecodeLegacyAuctionDetails(data)
	})
}

func FuzzDecodeInteraction(f *testing.F) {
	f.Add("0x111111125421ca6dc452d289314280a0f8842a65deadbeef")
	f.Add("0x")
	f.Add("nothex")
	f.Fuzz(func(t *testing.T, data string) {
		_, _ = DecodeInteraction(data)
	})
}
