package orderbook

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strings"
)

const nullAddress = "0x0000000000000000000000000000000000000000"

// packFeeParameter encodes integrator and resolver fee details into a single 48-bit value with error checks on input ranges.
// Returns the packed 48-bit value and an error if any input values exceed their allowed ranges.
func packFeeParameter(integratorFee *IntegratorFee, resolverFee *ResolverFee) (uint64, error) {
	var integratorFeeValue uint64 = 0
	var integratorShare uint64 = 0
	var resolverFeeValue uint64 = 0
	var resolverDiscount uint64 = 0

	if integratorFee != nil {
		integratorFeeValue = uint64(integratorFee.Fee * 10) // Convert to basis points * 10
		integratorShare = uint64(integratorFee.Share / 100) // Convert percentage to 0-100 range
	}

	if resolverFee != nil {
		resolverFeeValue = uint64(resolverFee.Fee * 10)                // Convert to basis points * 10
		resolverDiscount = uint64(100 - resolverFee.WhitelistDiscount) // Invert discount (100 - discount%)
	}

	// Range checks to ensure values fit in their allocated bit space
	if integratorFeeValue > 0xffff {
		return 0, fmt.Errorf("integrator fee value must be between 0 and 65535, got %d", integratorFeeValue)
	}
	if integratorShare > 0xff {
		return 0, fmt.Errorf("integrator share must be between 0 and 255, got %d", integratorShare)
	}
	if resolverFeeValue > 0xffff {
		return 0, fmt.Errorf("resolver fee value must be between 0 and 65535, got %d", resolverFeeValue)
	}
	if resolverDiscount > 0xff {
		return 0, fmt.Errorf("resolver discount must be between 0 and 255, got %d", resolverDiscount)
	}

	packed := (integratorFeeValue << 32) | // bits 47-32 (16 bits)
		(integratorShare << 24) | // bits 31-24 (8 bits)
		(resolverFeeValue << 8) | // bits 23-8 (16 bits)
		resolverDiscount // bits 7-0 (8 bits)

	return packed, nil
}

// encodeWhitelist encodes a list of Ethereum addresses into a single *big.Int value with a specific bit layout.
// Returns an error if any address is invalid, has more than 255 addresses, or encounters issues during encoding.
func encodeWhitelist(whitelist []string) (*big.Int, error) {
	if len(whitelist) == 0 {
		return big.NewInt(0), nil
	}

	if len(whitelist) > 255 {
		return nil, fmt.Errorf("whitelist can have at most 255 addresses, got %d", len(whitelist))
	}

	encoded := big.NewInt(int64(len(whitelist)))

	mask80 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 80), big.NewInt(1))

	for _, address := range whitelist {
		address = strings.TrimPrefix(address, "0x")
		address = strings.TrimPrefix(address, "0X") // Check for capital x too

		addressInt := new(big.Int)
		_, success := addressInt.SetString(address, 16)
		if !success {
			return nil, fmt.Errorf("invalid hex address: %s", address)
		}

		// Mask to get only lower 80 bits
		addressLower80 := new(big.Int).And(addressInt, mask80)

		// Shift encoded left by 80 bits and OR with the address
		encoded.Lsh(encoded, 80)
		encoded.Or(encoded, addressLower80)
	}

	return encoded, nil
}

// concatFeeAndWhitelist combines fee parameters and whitelist into a single encoded value
// Returns the combined value and the total bit length
func concatFeeAndWhitelist(whitelist []string, integratorFee *IntegratorFee, resolverFee *ResolverFee) (*big.Int, int, error) {
	feeParameter, err := packFeeParameter(integratorFee, resolverFee)
	if err != nil {
		return nil, 0, err
	}

	whitelistEncoded, err := encodeWhitelist(whitelist)
	if err != nil {
		return nil, 0, err
	}

	// Calculate whitelist bit length
	// If empty: 0 bits
	// Otherwise: 8 bits (length) + (number of addresses * 80 bits each)
	var whitelistBitLength int
	if len(whitelist) == 0 {
		whitelistBitLength = 0
	} else {
		whitelistBitLength = 8 + (len(whitelist) * 80)
	}

	// fee_parameter << whitelist_bit_length | whitelist_encoded
	feeAndWhitelist := new(big.Int).SetUint64(feeParameter)
	feeAndWhitelist.Lsh(feeAndWhitelist, uint(whitelistBitLength))
	feeAndWhitelist.Or(feeAndWhitelist, whitelistEncoded)

	totalBitLength := 48 + whitelistBitLength

	return feeAndWhitelist, totalBitLength, nil
}

// buildFeePostInteractionData encodes interaction data combining receiver, fee, whitelist, and optional interaction inputs into a byte array.
func buildFeePostInteractionData(params *buildFeePostInteractionDataParams) ([]byte, error) {

	if params.Whitelist == nil {
		params.Whitelist = []string{}
	}

	if params.CustomReceiverAddress == "" {
		params.CustomReceiverAddress = nullAddress
	}

	if params.ExtraInteractionTarget == "" {
		params.ExtraInteractionTarget = nullAddress
	}

	postInteraction := big.NewInt(0)

	// Add integrator address (20 bytes = 160 bits)
	integratorAddress := nullAddress
	if params.IntegratorFee != nil && params.IntegratorFee.Integrator != "" {
		integratorAddress = params.IntegratorFee.Integrator
	}
	integratorInt := new(big.Int)
	integratorInt.SetString(strings.TrimPrefix(integratorAddress, "0x"), 16)
	postInteraction.Or(postInteraction, integratorInt)

	// Add resolver/protocol fee receiver address (20 bytes = 160 bits)
	resolverAddress := nullAddress
	if params.ResolverFee != nil && params.ResolverFee.Receiver != "" {
		resolverAddress = params.ResolverFee.Receiver
	}
	resolverInt := new(big.Int)
	resolverInt.SetString(strings.TrimPrefix(resolverAddress, "0x"), 16)
	postInteraction.Lsh(postInteraction, 160)
	postInteraction.Or(postInteraction, resolverInt)

	// Add custom receiver if present (20 bytes = 160 bits)
	if params.CustomReceiver && params.CustomReceiverAddress != nullAddress {
		receiverInt := new(big.Int)
		receiverInt.SetString(strings.TrimPrefix(params.CustomReceiverAddress, "0x"), 16)
		postInteraction.Lsh(postInteraction, 160)
		postInteraction.Or(postInteraction, receiverInt)
	}

	// Add fee and whitelist data
	feeAndWhitelist, feeAndWhitelistLength, err := concatFeeAndWhitelist(params.Whitelist, params.IntegratorFee, params.ResolverFee)
	if err != nil {
		return nil, err
	}
	postInteraction.Lsh(postInteraction, uint(feeAndWhitelistLength))
	postInteraction.Or(postInteraction, feeAndWhitelist)

	// Add extra interaction if present
	if params.ExtraInteractionTarget != nullAddress && len(params.ExtraInteractionData) > 0 {
		targetInt := new(big.Int)
		targetInt.SetString(strings.TrimPrefix(params.ExtraInteractionTarget, "0x"), 16)
		postInteraction.Lsh(postInteraction, 160)
		postInteraction.Or(postInteraction, targetInt)

		// Add extra interaction data (preserve all bits including leading zeros)
		extraDataInt := new(big.Int).SetBytes(params.ExtraInteractionData)
		postInteraction.Lsh(postInteraction, uint(len(params.ExtraInteractionData)*8))
		postInteraction.Or(postInteraction, extraDataInt)
	}

	// Calculate expected byte length
	expectedLength := 1 + 20 + 20 // flag + integrator + resolver
	if params.CustomReceiver && params.CustomReceiverAddress != nullAddress {
		expectedLength += 20 // custom receiver
	}
	if len(params.Whitelist) == 0 {
		expectedLength += 6 // just fee_parameter (48 bits = 6 bytes)
	} else {
		expectedLength += (48 + 8 + len(params.Whitelist)*80) / 8 // fee + whitelist
	}
	if params.ExtraInteractionTarget != nullAddress && len(params.ExtraInteractionData) > 0 {
		expectedLength += 20 + len(params.ExtraInteractionData)
	}

	postInteractionBytes := postInteraction.Bytes()

	result := make([]byte, expectedLength)
	if params.CustomReceiver {
		result[0] = 0x01
	} else {
		result[0] = 0x00
	}

	copy(result[expectedLength-len(postInteractionBytes):], postInteractionBytes)

	return result, nil
}

// BuildOrderExtensionBytes builds the complete order extension
// Returns the encoded extension as a hex string
func BuildOrderExtensionBytes(params *BuildOrderExtensionBytesParams) (string, error) {

	var whiteListResolvers []string
	for _, value := range params.Whitelist {
		whiteListResolvers = append(whiteListResolvers, value)
		sort.Strings(whiteListResolvers) // Sorting the final list to ensure a deterministic order
	}

	feePostInteraction, err := buildFeePostInteractionData(&buildFeePostInteractionDataParams{
		CustomReceiver:         params.CustomReceiver != "" && params.CustomReceiver != nullAddress,
		CustomReceiverAddress:  params.CustomReceiver,
		IntegratorFee:          params.IntegratorFee,
		ResolverFee:            params.ResolverFee,
		Whitelist:              whiteListResolvers,
		ExtraInteractionTarget: params.ExtensionTarget,
		ExtraInteractionData:   params.ExtraInteraction,
	})
	if err != nil {
		return "", err
	}

	makingTakingAmountData, makingTakingAmountDataLength, err := concatFeeAndWhitelist(whiteListResolvers, params.IntegratorFee, params.ResolverFee)
	if err != nil {
		return "", err
	}

	extensionTargetBytes := make([]byte, 20)
	targetInt := new(big.Int)
	targetInt.SetString(strings.TrimPrefix(params.ExtensionTarget, "0x"), 16)
	targetBytes := targetInt.Bytes()
	copy(extensionTargetBytes[20-len(targetBytes):], targetBytes)

	makingTakingBytes := makingTakingAmountData.Bytes()
	expectedMakingTakingLen := (makingTakingAmountDataLength + 7) / 8 // Round up to bytes
	if len(makingTakingBytes) < expectedMakingTakingLen {
		padded := make([]byte, expectedMakingTakingLen)
		copy(padded[expectedMakingTakingLen-len(makingTakingBytes):], makingTakingBytes)
		makingTakingBytes = padded
	}

	makingTakingAmountDataBytes := append(extensionTargetBytes, makingTakingBytes...)

	// Prepend extension target to fee post-interaction
	feePostInteractionWithTarget := append(extensionTargetBytes, feePostInteraction...)

	interactions := [][]byte{
		{},                           // MakerAssetSuffix (empty)
		{},                           // TakerAssetSuffix (empty)
		makingTakingAmountDataBytes,  // MakingAmountData
		makingTakingAmountDataBytes,  // TakingAmountData (same as making)
		{},                           // Predicate (empty)
		params.MakerPermit,           // MakerPermit
		{},                           // PreInteractionData (empty)
		feePostInteractionWithTarget, // PostInteractionData
	}

	// Add customData if present
	if len(params.CustomData) > 0 {
		interactions = append(interactions, params.CustomData)
	}

	extension := buildExtensionFromBytes(interactions)

	return extension, nil
}

// buildExtensionFromBytes builds an extension hex string from byte slices
func buildExtensionFromBytes(interactions [][]byte) string {
	var byteCounts []int
	var dataBytes []byte

	// Process first 8 interactions (these are used in the cumulative sum calculation)
	for i := 0; i < len(interactions) && i < 8; i++ {
		byteCounts = append(byteCounts, len(interactions[i]))
		dataBytes = append(dataBytes, interactions[i]...)
	}

	// Add customData if present (no offset data)
	if len(interactions) > 8 {
		dataBytes = append(dataBytes, interactions[8]...)
	}

	// Calculate cumulative offsets
	cumulativeSum := 0
	var offsets []byte
	for i := 0; i < len(byteCounts); i++ {
		cumulativeSum += byteCounts[i]
		offsetBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(offsetBytes, uint32(cumulativeSum))
		offsets = append(offsetBytes, offsets...)
	}

	// If no data, return an empty extension
	if len(dataBytes) == 0 {
		return "0x"
	}

	offsetsHex := hex.EncodeToString(offsets)
	dataHex := hex.EncodeToString(dataBytes)

	return "0x" + offsetsHex + dataHex
}
