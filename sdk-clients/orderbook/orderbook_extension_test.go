package orderbook

import (
	"fmt"
	"math/big"
	"testing"
)

func TestPackFeeParameter(t *testing.T) {
	tests := []struct {
		name           string
		integratorFee  *IntegratorFee
		resolverFee    *ResolverFee
		expectedPacked uint64
		expectedHex    string
		expectError    bool
	}{
		{
			name: "example from Python - zero integrator fee, resolver fee 50 bps with 50% discount",
			integratorFee: &IntegratorFee{
				Integrator: "0x0000000000000000000000000000000000000000",
				Protocol:   "0x0000000000000000000000000000000000000000",
				Fee:        0,
				Share:      0,
			},
			resolverFee: &ResolverFee{
				Receiver:          "0x90cbe4bdd538d6e9b379bff5fe72c3d67a521de5",
				Fee:               50,
				WhitelistDiscount: 50,
			},
			expectedPacked: 128050, // 0x1f432
			expectedHex:    "0x1f432",
			expectError:    false,
		},
		{
			name:           "nil fees - should return 0",
			integratorFee:  nil,
			resolverFee:    nil,
			expectedPacked: 0,
			expectedHex:    "0x0",
			expectError:    false,
		},
		{
			name: "integrator fee only",
			integratorFee: &IntegratorFee{
				Integrator: "0x1111111111111111111111111111111111111111",
				Protocol:   "0x2222222222222222222222222222222222222222",
				Fee:        10,  // 0.1%
				Share:      500, // 5%
			},
			resolverFee:    nil,
			expectedPacked: (100 << 32) | (5 << 24), // fee*10=100, share/100=5
			expectedHex:    fmt.Sprintf("0x%x", (100<<32)|(5<<24)),
			expectError:    false,
		},
		{
			name:          "resolver fee only",
			integratorFee: nil,
			resolverFee: &ResolverFee{
				Receiver:          "0x3333333333333333333333333333333333333333",
				Fee:               100, // 1%
				WhitelistDiscount: 25,  // 25% discount
			},
			expectedPacked: (1000 << 8) | 75, // fee*10=1000, 100-25=75
			expectedHex:    fmt.Sprintf("0x%x", (1000<<8)|75),
			expectError:    false,
		},
		{
			name: "both fees with various values",
			integratorFee: &IntegratorFee{
				Fee:   1,   // 0.01%
				Share: 500, // 5%
			},
			resolverFee: &ResolverFee{
				Fee:               50, // 0.5%
				WhitelistDiscount: 50, // 50% discount
			},
			expectedPacked: (10 << 32) | (5 << 24) | (500 << 8) | 50,
			expectedHex:    fmt.Sprintf("0x%x", (10<<32)|(5<<24)|(500<<8)|50),
			expectError:    false,
		},
		{
			name: "maximum valid values",
			integratorFee: &IntegratorFee{
				Fee:   6553, // Maximum value that gives 65530 after *10
				Share: 25500,
			},
			resolverFee: &ResolverFee{
				Fee:               6553,
				WhitelistDiscount: 0,
			},
			expectedPacked: (65530 << 32) | (255 << 24) | (65530 << 8) | 100,
			expectedHex:    fmt.Sprintf("0x%x", (65530<<32)|(255<<24)|(65530<<8)|100),
			expectError:    false,
		},
		{
			name: "integrator fee too large",
			integratorFee: &IntegratorFee{
				Fee:   6554, // 6554 * 10 = 65540, exceeds 0xffff
				Share: 0,
			},
			resolverFee: nil,
			expectError: true,
		},
		{
			name: "integrator share too large",
			integratorFee: &IntegratorFee{
				Fee:   0,
				Share: 25600, // 25600 / 100 = 256, exceeds 0xff
			},
			resolverFee: nil,
			expectError: true,
		},
		{
			name:          "resolver fee too large",
			integratorFee: nil,
			resolverFee: &ResolverFee{
				Fee:               6554, // 6554 * 10 = 65540, exceeds 0xffff
				WhitelistDiscount: 0,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packed, err := packFeeParameter(tt.integratorFee, tt.resolverFee)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if packed != tt.expectedPacked {
				t.Errorf("packed value mismatch:\nexpected: %d (0x%x)\ngot:      %d (0x%x)",
					tt.expectedPacked, tt.expectedPacked, packed, packed)
			}

			// Verify the hex representation matches
			gotHex := fmt.Sprintf("0x%x", packed)
			if gotHex != tt.expectedHex {
				t.Errorf("hex representation mismatch:\nexpected: %s\ngot:      %s",
					tt.expectedHex, gotHex)
			}
		})
	}
}

func TestPackFeeParameterBitLayout(t *testing.T) {
	// Test to verify the bit layout is correct
	integratorFee := &IntegratorFee{
		Fee:   1,   // becomes 10 after *10
		Share: 200, // becomes 2 after /100
	}
	resolverFee := &ResolverFee{
		Fee:               3,  // becomes 30 after *10
		WhitelistDiscount: 25, // becomes 75 after (100-25)
	}

	packed, err := packFeeParameter(integratorFee, resolverFee)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Extract each component and verify
	extractedIntegratorFee := (packed >> 32) & 0xffff
	extractedIntegratorShare := (packed >> 24) & 0xff
	extractedResolverFee := (packed >> 8) & 0xffff
	extractedResolverDiscount := packed & 0xff

	if extractedIntegratorFee != 10 {
		t.Errorf("integrator fee: expected 10, got %d", extractedIntegratorFee)
	}
	if extractedIntegratorShare != 2 {
		t.Errorf("integrator share: expected 2, got %d", extractedIntegratorShare)
	}
	if extractedResolverFee != 30 {
		t.Errorf("resolver fee: expected 30, got %d", extractedResolverFee)
	}
	if extractedResolverDiscount != 75 {
		t.Errorf("resolver discount: expected 75, got %d", extractedResolverDiscount)
	}
}

func TestEncodeWhitelist(t *testing.T) {
	tests := []struct {
		name        string
		whitelist   []string
		expectedHex string
		expectedInt string
		expectError bool
	}{
		{
			name:        "empty whitelist",
			whitelist:   []string{},
			expectedHex: "0x0",
			expectedInt: "0",
			expectError: false,
		},
		{
			name: "single address",
			whitelist: []string{
				"0x0b8a49d816cc709b6eadb09498030ae3416b66dc",
			},
			// Length: 1 (0x01)
			// Address: 0x0b8a49d816cc709b6eadb09498030ae3416b66dc
			// Lower 80 bits of address: 0xb09498030ae3416b66dc (not 0x9498... - the 'b0' is included!)
			// Result: 0x01 << 80 | 0xb09498030ae3416b66dc = 0x01b09498030ae3416b66dc
			expectedHex: "0x1b09498030ae3416b66dc",
			expectedInt: "2042803392333285612545756",
			expectError: false,
		},
		{
			name: "two addresses",
			whitelist: []string{
				"0x0b8a49d816cc709b6eadb09498030ae3416b66dc",
				"0xad3b67bca8935cb510c8d18bd45f0b94f54a968f",
			},
			// Length: 2 (0x02)
			// Addr1: 0x0b8a49d816cc709b6eadb09498030ae3416b66dc -> lower 80 bits: 0xb09498030ae3416b66dc
			// Addr2: 0xad3b67bca8935cb510c8d18bd45f0b94f54a968f -> lower 80 bits: 0xd18bd45f0b94f54a968f
			// Result: 0x02 << 160 | addr1 << 80 | addr2
			expectedHex: "0x2b09498030ae3416b66dcd18bd45f0b94f54a968f",
			expectedInt: "3931099402718965111428264566673874676661292275343",
			expectError: false,
		},
		{
			name: "six addresses from Python example",
			whitelist: []string{
				"0x0b8a49d816cc709b6eadb09498030ae3416b66dc",
				"0xad3b67bca8935cb510c8d18bd45f0b94f54a968f",
				"0xf81377c3f03996fde219c90ed87a54c23dc480b3",
				"0xbeef02961503351625926ea9a11ae13b29f5c555",
				"0x00000688768803bbd44095770895ad27ad6b0d95",
				"0xf0a12fefa78181a3749474db31d09524fa87b1f7",
			},
			expectedHex: "0x6b09498030ae3416b66dcd18bd45f0b94f54a968fc90ed87a54c23dc480b36ea9a11ae13b29f5c55595770895ad27ad6b0d9574db31d09524fa87b1f7",
			expectedInt: "20883771562388111227015435143042649718086364789015763900152491422545023853078207873061945043557603652822441846333414299315624750241895600466932215",
			expectError: false,
		},
		{
			name: "address without 0x prefix",
			whitelist: []string{
				"0b8a49d816cc709b6eadb09498030ae3416b66dc",
			},
			// Same address, just without 0x prefix
			// Lower 80 bits: 0xb09498030ae3416b66dc
			expectedHex: "0x1b09498030ae3416b66dc",
			expectedInt: "2042803392333285612545756",
			expectError: false,
		},
		{
			name: "too many addresses",
			whitelist: func() []string {
				addresses := make([]string, 256)
				for i := 0; i < 256; i++ {
					addresses[i] = "0x0b8a49d816cc709b6eadb09498030ae3416b66dc"
				}
				return addresses
			}(),
			expectError: true,
		},
		{
			name: "invalid hex address",
			whitelist: []string{
				"0xGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := encodeWhitelist(tt.whitelist)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Check hex representation
			gotHex := "0x" + encoded.Text(16)
			if gotHex != tt.expectedHex {
				t.Errorf("hex mismatch:\nexpected: %s\ngot:      %s", tt.expectedHex, gotHex)
			}

			// Check integer representation
			gotInt := encoded.String()
			if gotInt != tt.expectedInt {
				t.Errorf("int mismatch:\nexpected: %s\ngot:      %s", tt.expectedInt, gotInt)
			}

			// Check bit length for the Python example
			if tt.name == "six addresses from Python example" {
				bitLength := encoded.BitLen()
				expectedBitLength := 483
				if bitLength != expectedBitLength {
					t.Errorf("bit length mismatch:\nexpected: %d\ngot:      %d", expectedBitLength, bitLength)
				}
			}
		})
	}
}

func TestEncodeWhitelistBitLayout(t *testing.T) {
	// Test to verify addresses are properly masked to 80 bits
	whitelist := []string{
		"0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", // Full address with all bits set
	}

	encoded, err := encodeWhitelist(whitelist)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The encoding should be: length (1) << 80 | (address & 80-bit mask)
	// 80-bit mask of 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF = 0xFFFFFFFFFFFFFFFFFFFF
	expected := new(big.Int)
	expected.SetString("1ffffffffffffffffffff", 16) // 0x01 << 80 | 0xFFFFFFFFFFFFFFFFFFFF

	if encoded.Cmp(expected) != 0 {
		t.Errorf("masking failed:\nexpected: 0x%x\ngot:      0x%x", expected, encoded)
	}

	// Verify the address was properly truncated to 80 bits (10 bytes)
	// Extract the address part (lower 80 bits)
	addressPart := new(big.Int).Set(encoded)
	mask80 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 80), big.NewInt(1))
	addressPart.And(addressPart, mask80)

	expectedAddr := new(big.Int)
	expectedAddr.SetString("ffffffffffffffffffff", 16)

	if addressPart.Cmp(expectedAddr) != 0 {
		t.Errorf("address masking failed:\nexpected: 0x%x\ngot:      0x%x", expectedAddr, addressPart)
	}
}

func TestEncodeWhitelistOrderPreservation(t *testing.T) {
	// Test that the order of addresses is preserved in encoding
	whitelist := []string{
		"0x1111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222",
		"0x3333333333333333333333333333333333333333",
	}

	encoded, err := encodeWhitelist(whitelist)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Manually decode and verify order
	// Expected format: [length=3][addr1_lower80][addr2_lower80][addr3_lower80]

	// Extract addresses by shifting
	mask80 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 80), big.NewInt(1))

	temp := new(big.Int).Set(encoded)

	// Extract third address (rightmost 80 bits)
	addr3 := new(big.Int).And(temp, mask80)
	temp.Rsh(temp, 80)

	// Extract second address
	addr2 := new(big.Int).And(temp, mask80)
	temp.Rsh(temp, 80)

	// Extract first address
	addr1 := new(big.Int).And(temp, mask80)

	// Verify addresses match (lower 80 bits)
	expected1 := new(big.Int)
	expected1.SetString("1111111111111111111111111111111111111111", 16)
	expected1.And(expected1, mask80)

	expected2 := new(big.Int)
	expected2.SetString("2222222222222222222222222222222222222222", 16)
	expected2.And(expected2, mask80)

	expected3 := new(big.Int)
	expected3.SetString("3333333333333333333333333333333333333333", 16)
	expected3.And(expected3, mask80)

	if addr1.Cmp(expected1) != 0 {
		t.Errorf("first address mismatch:\nexpected: 0x%x\ngot:      0x%x", expected1, addr1)
	}
	if addr2.Cmp(expected2) != 0 {
		t.Errorf("second address mismatch:\nexpected: 0x%x\ngot:      0x%x", expected2, addr2)
	}
	if addr3.Cmp(expected3) != 0 {
		t.Errorf("third address mismatch:\nexpected: 0x%x\ngot:      0x%x", expected3, addr3)
	}
}

func TestEncodeWhitelistMaxLength(t *testing.T) {
	// Test with exactly 255 addresses (maximum allowed)
	whitelist := make([]string, 255)
	for i := 0; i < 255; i++ {
		whitelist[i] = "0x0b8a49d816cc709b6eadb09498030ae3416b66dc"
	}

	encoded, err := encodeWhitelist(whitelist)
	if err != nil {
		t.Fatalf("unexpected error with 255 addresses: %v", err)
	}

	// Verify length is 255
	// Extract the highest byte(s) which contain the length
	temp := new(big.Int).Set(encoded)
	temp.Rsh(temp, 255*80) // Shift past all addresses

	if temp.Int64() != 255 {
		t.Errorf("length mismatch:\nexpected: 255\ngot:      %d", temp.Int64())
	}
}

func TestConcatFeeAndWhitelist(t *testing.T) {
	tests := []struct {
		name              string
		whitelist         []string
		integratorFee     *IntegratorFee
		resolverFee       *ResolverFee
		expectedHex       string
		expectedInt       string
		expectedBitLength int
		expectError       bool
	}{
		{
			name:      "empty whitelist with fees",
			whitelist: []string{},
			integratorFee: &IntegratorFee{
				Fee:   0,
				Share: 0,
			},
			resolverFee: &ResolverFee{
				Fee:               50,
				WhitelistDiscount: 50,
			},
			// fee_parameter = 0x1f432 (128050)
			// whitelist is empty (0 bits)
			// Result: 0x1f432 << 0 | 0 = 0x1f432
			expectedHex:       "0x1f432",
			expectedInt:       "128050",
			expectedBitLength: 48,
			expectError:       false,
		},
		{
			name: "six addresses from Python example",
			whitelist: []string{
				"0x0b8a49d816cc709b6eadb09498030ae3416b66dc",
				"0xad3b67bca8935cb510c8d18bd45f0b94f54a968f",
				"0xf81377c3f03996fde219c90ed87a54c23dc480b3",
				"0xbeef02961503351625926ea9a11ae13b29f5c555",
				"0x00000688768803bbd44095770895ad27ad6b0d95",
				"0xf0a12fefa78181a3749474db31d09524fa87b1f7",
			},
			integratorFee: &IntegratorFee{
				Integrator: "0x0000000000000000000000000000000000000000",
				Protocol:   "0x0000000000000000000000000000000000000000",
				Fee:        0,
				Share:      0,
			},
			resolverFee: &ResolverFee{
				Receiver:          "0x90cbe4bdd538d6e9b379bff5fe72c3d67a521de5",
				Fee:               50,
				WhitelistDiscount: 50,
			},
			// fee_parameter = 0x1f432 (128050)
			// whitelist_encoded = 0x6b09498030ae3416b66dcd18bd45f0b94f54a968fc90ed87a54c23dc480b36ea9a11ae13b29f5c55595770895ad27ad6b0d9574db31d09524fa87b1f7
			// whitelist_bit_length = 8 + 6*80 = 488 bits
			// Result: 0x1f432 << 488 | whitelist_encoded
			expectedHex:       "0x1f43206b09498030ae3416b66dcd18bd45f0b94f54a968fc90ed87a54c23dc480b36ea9a11ae13b29f5c55595770895ad27ad6b0d9574db31d09524fa87b1f7",
			expectedInt:       "102333435761970040526585089485838969078133364081436675317847752614578808440926825918247603365248761165799707486146452557225400316772363664864037468353015",
			expectedBitLength: 536,
			expectError:       false,
		},
		{
			name: "single address",
			whitelist: []string{
				"0x0b8a49d816cc709b6eadb09498030ae3416b66dc",
			},
			integratorFee: &IntegratorFee{
				Fee:   10,
				Share: 500,
			},
			resolverFee: &ResolverFee{
				Fee:               100,
				WhitelistDiscount: 25,
			},
			// fee_parameter: integrator_fee=100, integrator_share=5, resolver_fee=1000, resolver_discount=75
			// = (100 << 32) | (5 << 24) | (1000 << 8) | 75
			// = 0x640503e84b
			// whitelist_encoded = 0x01b09498030ae3416b66dc (note: 0x0b8a... truncated to lower 80 bits = 0xb09498...)
			// whitelist_bit_length = 8 + 1*80 = 88 bits
			// Result: 0x640503e84b << 88 | 0x01b09498030ae3416b66dc
			expectedHex:       "0x640503e84b01b09498030ae3416b66dc",
			expectedInt:       "132948840314160194232854165493602019036",
			expectedBitLength: 136,
			expectError:       false,
		},
		{
			name: "nil fees with addresses",
			whitelist: []string{
				"0x1111111111111111111111111111111111111111",
			},
			integratorFee: nil,
			resolverFee:   nil,
			// fee_parameter = 0
			// Address 0x1111111111111111111111111111111111111111 masked to lower 80 bits = 0x1111111111111111111111
			// whitelist_encoded = 0x011111111111111111111111 (length=1, then 80-bit address)
			// whitelist_bit_length = 88 bits
			// Result: 0 << 88 | 0x011111111111111111111111
			// Note: big.Int.Text(16) drops leading zeros, so 0x011... becomes 0x11...
			expectedHex:       "0x111111111111111111111",
			expectedInt:       "1289520874255604453019921",
			expectedBitLength: 136,
			expectError:       false,
		},
		{
			name: "two addresses with fees",
			whitelist: []string{
				"0x1111111111111111111111111111111111111111",
				"0x2222222222222222222222222222222222222222",
			},
			integratorFee: &IntegratorFee{
				Fee:   1,
				Share: 1000,
			},
			resolverFee: &ResolverFee{
				Fee:               1,
				WhitelistDiscount: 0,
			},
			// fee_parameter: (10 << 32) | (10 << 24) | (10 << 8) | 100
			// whitelist_bit_length = 8 + 2*80 = 168 bits
			expectedBitLength: 216,
			expectError:       false,
		},
		{
			name: "invalid whitelist causes error",
			whitelist: []string{
				"0xGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG",
			},
			integratorFee: nil,
			resolverFee:   nil,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, bitLength, err := concatFeeAndWhitelist(tt.whitelist, tt.integratorFee, tt.resolverFee)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Check bit length
			if bitLength != tt.expectedBitLength {
				t.Errorf("bit length mismatch:\nexpected: %d\ngot:      %d", tt.expectedBitLength, bitLength)
			}

			// Check hex representation if provided
			if tt.expectedHex != "" {
				gotHex := "0x" + result.Text(16)
				if gotHex != tt.expectedHex {
					t.Errorf("hex mismatch:\nexpected: %s\ngot:      %s", tt.expectedHex, gotHex)
				}
			}

			// Check integer representation if provided
			if tt.expectedInt != "" {
				gotInt := result.String()
				if gotInt != tt.expectedInt {
					t.Errorf("int mismatch:\nexpected: %s\ngot:      %s", tt.expectedInt, gotInt)
				}
			}
		})
	}
}

func TestConcatFeeAndWhitelistBitLayout(t *testing.T) {
	// Test to verify the bit layout: [fee_parameter][whitelist_encoded]
	whitelist := []string{
		"0x1111111111111111111111111111111111111111",
	}
	integratorFee := &IntegratorFee{
		Fee:   1,
		Share: 100,
	}
	resolverFee := &ResolverFee{
		Fee:               2,
		WhitelistDiscount: 50,
	}

	result, bitLength, err := concatFeeAndWhitelist(whitelist, integratorFee, resolverFee)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Expected bit length: 48 (fee) + 88 (1 address with length byte) = 136
	expectedBitLength := 136
	if bitLength != expectedBitLength {
		t.Errorf("bit length mismatch:\nexpected: %d\ngot:      %d", expectedBitLength, bitLength)
	}

	// Extract fee parameter (upper 48 bits)
	temp := new(big.Int).Set(result)
	temp.Rsh(temp, 88) // Shift right by whitelist bit length
	feeParam := temp.Uint64()

	// Calculate expected fee parameter
	expectedFeeParam, _ := packFeeParameter(integratorFee, resolverFee)
	if feeParam != expectedFeeParam {
		t.Errorf("fee parameter mismatch:\nexpected: 0x%x\ngot:      0x%x", expectedFeeParam, feeParam)
	}

	// Extract whitelist (lower 88 bits)
	mask88 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 88), big.NewInt(1))
	whitelistPart := new(big.Int).And(result, mask88)

	// Calculate expected whitelist
	expectedWhitelist, _ := encodeWhitelist(whitelist)
	if whitelistPart.Cmp(expectedWhitelist) != 0 {
		t.Errorf("whitelist mismatch:\nexpected: 0x%x\ngot:      0x%x", expectedWhitelist, whitelistPart)
	}
}

func TestConcatFeeAndWhitelistEmptyWhitelist(t *testing.T) {
	// Test that empty whitelist results in just the fee parameter
	integratorFee := &IntegratorFee{
		Fee:   50,
		Share: 2500,
	}
	resolverFee := &ResolverFee{
		Fee:               100,
		WhitelistDiscount: 30,
	}

	result, bitLength, err := concatFeeAndWhitelist([]string{}, integratorFee, resolverFee)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// With empty whitelist, result should just be the fee parameter
	expectedFeeParam, _ := packFeeParameter(integratorFee, resolverFee)

	if result.Uint64() != expectedFeeParam {
		t.Errorf("with empty whitelist, result should equal fee parameter:\nexpected: 0x%x\ngot:      0x%x",
			expectedFeeParam, result.Uint64())
	}

	// Bit length should be exactly 48
	if bitLength != 48 {
		t.Errorf("empty whitelist bit length should be 48, got: %d", bitLength)
	}
}

func TestBuildFeePostInteractionData(t *testing.T) {
	tests := []struct {
		name                   string
		customReceiver         bool
		customReceiverAddress  string
		integratorFee          *IntegratorFee
		resolverFee            *ResolverFee
		whitelist              []string
		extraInteractionTarget string
		extraInteractionData   []byte
		expectedHex            string
		expectedLength         int
		expectError            bool
	}{
		{
			name:                  "Python example - no custom receiver, 6 addresses in whitelist",
			customReceiver:        false,
			customReceiverAddress: "",
			integratorFee: &IntegratorFee{
				Integrator: "0x0000000000000000000000000000000000000000",
				Protocol:   "0x0000000000000000000000000000000000000000",
				Fee:        0,
				Share:      0,
			},
			resolverFee: &ResolverFee{
				Receiver:          "0x90cbe4bdd538d6e9b379bff5fe72c3d67a521de5",
				Fee:               50,
				WhitelistDiscount: 50,
			},
			whitelist: []string{
				"0x0b8a49d816cc709b6eadb09498030ae3416b66dc",
				"0xad3b67bca8935cb510c8d18bd45f0b94f54a968f",
				"0xf81377c3f03996fde219c90ed87a54c23dc480b3",
				"0xbeef02961503351625926ea9a11ae13b29f5c555",
				"0x00000688768803bbd44095770895ad27ad6b0d95",
				"0xf0a12fefa78181a3749474db31d09524fa87b1f7",
			},
			extraInteractionTarget: "0x0000000000000000000000000000000000000000",
			extraInteractionData:   []byte{},
			expectedHex:            "00000000000000000000000000000000000000000090cbe4bdd538d6e9b379bff5fe72c3d67a521de500000001f43206b09498030ae3416b66dcd18bd45f0b94f54a968fc90ed87a54c23dc480b36ea9a11ae13b29f5c55595770895ad27ad6b0d9574db31d09524fa87b1f7",
			expectedLength:         108,
			expectError:            false,
		},
		{
			name:                   "empty whitelist, no fees",
			customReceiver:         false,
			customReceiverAddress:  "",
			integratorFee:          nil,
			resolverFee:            nil,
			whitelist:              []string{},
			extraInteractionTarget: "",
			extraInteractionData:   nil,
			// Structure: [1 byte flag: 0x00][20 bytes zeros][20 bytes zeros][6 bytes fee_parameter: 0x000000000000]
			// Total: 1 + 20 + 20 + 6 = 47 bytes = 94 hex chars
			expectedHex:    "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			expectedLength: 47,
			expectError:    false,
		},
		{
			name:                  "with custom receiver",
			customReceiver:        true,
			customReceiverAddress: "0x1111111111111111111111111111111111111111",
			integratorFee: &IntegratorFee{
				Integrator: "0x2222222222222222222222222222222222222222",
				Protocol:   "0x3333333333333333333333333333333333333333",
				Fee:        10,
				Share:      500,
			},
			resolverFee: &ResolverFee{
				Receiver:          "0x4444444444444444444444444444444444444444",
				Fee:               20,
				WhitelistDiscount: 30,
			},
			whitelist:              []string{},
			extraInteractionTarget: "",
			extraInteractionData:   nil,
			// First byte should be 0x01 for custom receiver
			expectedLength: 67, // 1 + 20 + 20 + 20 + 6
			expectError:    false,
		},
		{
			name:                   "with extra interaction",
			customReceiver:         false,
			integratorFee:          nil,
			resolverFee:            nil,
			whitelist:              []string{},
			extraInteractionTarget: "0x5555555555555555555555555555555555555555",
			extraInteractionData:   []byte{0xde, 0xad, 0xbe, 0xef},
			expectedLength:         71, // 1 + 20 + 20 + 6 + 20 + 4
			expectError:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := buildFeePostInteractionData(&buildFeePostInteractionDataParams{
				CustomReceiver:         tt.customReceiver,
				CustomReceiverAddress:  tt.customReceiverAddress,
				IntegratorFee:          tt.integratorFee,
				ResolverFee:            tt.resolverFee,
				Whitelist:              tt.whitelist,
				ExtraInteractionTarget: tt.extraInteractionTarget,
				ExtraInteractionData:   tt.extraInteractionData,
			})

			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Check length
			if len(result) != tt.expectedLength {
				t.Errorf("length mismatch:\nexpected: %d\ngot:      %d", tt.expectedLength, len(result))
			}

			// Check hex if provided
			if tt.expectedHex != "" {
				gotHex := fmt.Sprintf("%x", result)
				if gotHex != tt.expectedHex {
					t.Errorf("hex mismatch:\nexpected: %s\ngot:      %s", tt.expectedHex, gotHex)
				}
			}

			// Verify first byte (flag)
			if tt.customReceiver {
				if result[0] != 0x01 {
					t.Errorf("expected first byte to be 0x01 for custom receiver, got: 0x%02x", result[0])
				}
			} else {
				if result[0] != 0x00 {
					t.Errorf("expected first byte to be 0x00 for no custom receiver, got: 0x%02x", result[0])
				}
			}
		})
	}
}

func TestBuildFeePostInteractionDataStructure(t *testing.T) {
	// Test to verify the byte structure is correct
	integratorFee := &IntegratorFee{
		Integrator: "0x1111111111111111111111111111111111111111",
		Protocol:   "0x2222222222222222222222222222222222222222",
		Fee:        5,
		Share:      1000,
	}
	resolverFee := &ResolverFee{
		Receiver:          "0x3333333333333333333333333333333333333333",
		Fee:               10,
		WhitelistDiscount: 20,
	}

	result, err := buildFeePostInteractionData(&buildFeePostInteractionDataParams{
		CustomReceiver:         false,
		CustomReceiverAddress:  "",
		IntegratorFee:          integratorFee,
		ResolverFee:            resolverFee,
		Whitelist:              []string{},
		ExtraInteractionTarget: "",
		ExtraInteractionData:   nil,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify structure:
	// [1 byte flag][20 bytes integrator][20 bytes resolver][6 bytes fee_parameter]
	if len(result) != 47 {
		t.Errorf("expected length 47, got %d", len(result))
	}

	// Check flag byte
	if result[0] != 0x00 {
		t.Errorf("expected flag byte 0x00, got 0x%02x", result[0])
	}

	// Check integrator address (bytes 1-20)
	expectedIntegrator := "1111111111111111111111111111111111111111"
	gotIntegrator := fmt.Sprintf("%x", result[1:21])
	if gotIntegrator != expectedIntegrator {
		t.Errorf("integrator address mismatch:\nexpected: %s\ngot:      %s", expectedIntegrator, gotIntegrator)
	}

	// Check resolver address (bytes 21-40)
	expectedResolver := "3333333333333333333333333333333333333333"
	gotResolver := fmt.Sprintf("%x", result[21:41])
	if gotResolver != expectedResolver {
		t.Errorf("resolver address mismatch:\nexpected: %s\ngot:      %s", expectedResolver, gotResolver)
	}
}

func TestBuildOrderExtension(t *testing.T) {
	extensionTarget := "0xc0dfdb9e7a392c3dbbe7c6fbe8fbc1789c9fe05e"
	integratorFee := &IntegratorFee{
		Integrator: "0x0000000000000000000000000000000000000000",
		Protocol:   "0x0000000000000000000000000000000000000000",
		Fee:        0,
		Share:      0,
	}
	resolverFee := &ResolverFee{
		Receiver:          "0x90cbe4bdd538d6e9b379bff5fe72c3d67a521de5",
		Fee:               50,
		WhitelistDiscount: 50,
	}
	whitelist := map[string]string{
		"0x9b8a49d816cc709b6eadb09498030ae3416b66dc": "0x0b8a49d816cc709b6eadb09498030ae3416b66dc",
		"0x9d3b67bca8935cb510c8d18bd45f0b94f54a968f": "0xad3b67bca8935cb510c8d18bd45f0b94f54a968f",
		"0x981377c3f03996fde219c90ed87a54c23dc480b3": "0xf81377c3f03996fde219c90ed87a54c23dc480b3",
		"0x9eef02961503351625926ea9a11ae13b29f5c555": "0xbeef02961503351625926ea9a11ae13b29f5c555",
		"0x90000688768803bbd44095770895ad27ad6b0d95": "0x00000688768803bbd44095770895ad27ad6b0d95",
		"0x90a12fefa78181a3749474db31d09524fa87b1f7": "0xf0a12fefa78181a3749474db31d09524fa87b1f7",
	}

	extension, err := BuildOrderExtensionBytes(&BuildOrderExtensionBytesParams{
		ExtensionTarget:  extensionTarget,
		IntegratorFee:    integratorFee,
		ResolverFee:      resolverFee,
		Whitelist:        whitelist,
		MakerPermit:      nil,      // makerPermit
		CustomReceiver:   "",       // customReceiver
		ExtraInteraction: []byte{}, // extraInteraction
		CustomData:       nil,      // customData
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedExtension := "0x0000012e000000ae000000ae000000ae000000ae000000570000000000000000c0dfdb9e7a392c3dbbe7c6fbe8fbc1789c9fe05e00000001f4320695770895ad27ad6b0d95b09498030ae3416b66dcd18bd45f0b94f54a968f6ea9a11ae13b29f5c55574db31d09524fa87b1f7c90ed87a54c23dc480b3c0dfdb9e7a392c3dbbe7c6fbe8fbc1789c9fe05e00000001f4320695770895ad27ad6b0d95b09498030ae3416b66dcd18bd45f0b94f54a968f6ea9a11ae13b29f5c55574db31d09524fa87b1f7c90ed87a54c23dc480b3c0dfdb9e7a392c3dbbe7c6fbe8fbc1789c9fe05e00000000000000000000000000000000000000000090cbe4bdd538d6e9b379bff5fe72c3d67a521de500000001f4320695770895ad27ad6b0d95b09498030ae3416b66dcd18bd45f0b94f54a968f6ea9a11ae13b29f5c55574db31d09524fa87b1f7c90ed87a54c23dc480b3"

	if extension != expectedExtension {
		t.Errorf("extension mismatch:\nexpected: %s\ngot:      %s", expectedExtension, extension)
	}

	// Verify length (670 chars as per Python output)
	if len(extension) != 670 {
		t.Errorf("extension length mismatch:\nexpected: 670\ngot:      %d", len(extension))
	}
}
