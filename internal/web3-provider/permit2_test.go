package web3_provider

//func TestHashPermitTransferFrom_MaxValues(t *testing.T) {
//	permit := PermitTransferFrom{
//		Permitted: TokenPermissions{
//			Token:  "0x0000000000000000000000000000000000000000",
//			Amount: MaxSignatureTransferAmount,
//		},
//		Spender:  "0x0000000000000000000000000000000000000000",
//		Nonce:    MaxUnorderedNonce,
//		Deadline: MaxSigDeadline,
//	}
//	permit2Address := "0x0000000000000000000000000000000000000000"
//	chainId := 1
//
//	// Call the hash function
//	hash := hashPermitTransferFrom(permit, permit2Address, chainId, nil)
//
//	// Define the expected hash value
//	expectedHash := "0x99e8cd5cd187c1dcb3c9cb41664cb12c1a3a76143d21b16f7880f4839d2b2ad4"
//
//	// Convert expectedHash to bytes for comparison
//	expectedHashBytes := common.Hex2Bytes(expectedHash)
//
//	// Convert hash to bytes for comparison
//	hashBytes := common.Hex2Bytes(hash)
//
//	// Compare the hashes
//	if !bytes.Equal(hashBytes, expectedHashBytes) {
//		t.Errorf("Expected hash %s, but got %s", expectedHash, hash)
//	}
//}
