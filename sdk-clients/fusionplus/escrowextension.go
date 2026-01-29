package fusionplus

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/1inch/1inch-sdk-go/internal/bytesiterator"
	"github.com/1inch/1inch-sdk-go/internal/hexadecimal"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
	"github.com/ethereum/go-ethereum/common"
)

type EscrowExtension struct {
	ExtensionPlus
	HashLock         *HashLock
	DstChainId       float32
	DstToken         common.Address
	SrcSafetyDeposit string
	DstSafetyDeposit string
	TimeLocks        TimeLocks
}

func NewEscrowExtension(escrowParams EscrowExtensionParams) (*EscrowExtension, error) {

	extension, err := NewExtensionPlus(escrowParams.ExtensionParamsPlus)
	if err != nil {
		return nil, err
	}

	escrowExtension := &EscrowExtension{
		ExtensionPlus:  *extension,
		HashLock:         escrowParams.HashLock,
		DstChainId:       escrowParams.DstChainId,
		DstToken:         escrowParams.DstToken,
		SrcSafetyDeposit: escrowParams.SrcSafetyDeposit,
		DstSafetyDeposit: escrowParams.DstSafetyDeposit,
		TimeLocks:        escrowParams.TimeLocks,
	}

	return escrowExtension, nil
}

func (e *EscrowExtension) ConvertToOrderbookExtension() (*orderbook.Extension, error) {

	srcSafetyDepositBig := new(big.Int)
	_, ok := srcSafetyDepositBig.SetString(e.SrcSafetyDeposit, 10)
	if !ok {
		return nil, fmt.Errorf("invalid source safety deposit hex: %s", e.SrcSafetyDeposit)
	}

	dstSafetyDepositBig := new(big.Int)
	_, ok = dstSafetyDepositBig.SetString(e.DstSafetyDeposit, 10)
	if !ok {
		return nil, fmt.Errorf("invalid destination safety deposit hex: %s", e.DstSafetyDeposit)
	}

	extraDataBytes, err := encodeExtraData(&EscrowExtraData{
		HashLock:         e.HashLock,
		DstChainId:       e.DstChainId,
		DstToken:         e.DstToken,
		SrcSafetyDeposit: srcSafetyDepositBig,
		DstSafetyDeposit: dstSafetyDepositBig,
		TimeLocks:        &e.TimeLocks,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to encode extra data: %w", err)
	}

	e.PostInteraction += hexadecimal.Trim0x(fmt.Sprintf("%x", extraDataBytes))

	return &orderbook.Extension{
		MakerAssetSuffix: e.MakerAssetSuffix,
		TakerAssetSuffix: e.TakerAssetSuffix,
		MakingAmountData: e.MakingAmountData,
		TakingAmountData: e.TakingAmountData,
		Predicate:        e.Predicate,
		MakerPermit:      e.MakerPermit,
		PreInteraction:   e.PreInteraction,
		PostInteraction:  e.PostInteraction,
		//hexadecimal.Trim0x(e.CustomData), // TODO Blocking custom data for now because it is breaking the cumsum method. The extension constructor will return with an error if the user provides this field.
	}, nil
}

// DecodeEscrowExtension decodes the input byte slice into an Extension struct using reflection.
func DecodeEscrowExtension(data []byte) (*EscrowExtension, error) {

	const extraDataCharacterLength = 320

	// Create one extension that will be used for the Escrow extension data
	orderbookExtensionTruncated, err := orderbook.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode extension: %w", err)
	}

	// Remove the Fusion Plus Extension data before decoding
	orderbookExtensionTruncated.PostInteraction = orderbookExtensionTruncated.PostInteraction[:len(orderbookExtensionTruncated.PostInteraction)-extraDataCharacterLength]
	extensionPlus, err := FromLimitOrderExtension(orderbookExtensionTruncated)
	if err != nil {
		return &EscrowExtension{}, fmt.Errorf("failed to decode escrow extension: %w", err)
	}

	// Create a second extension to extract extra data
	orderbookExtension, err := orderbook.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode extension: %w", err)
	}
	extraDataRaw := orderbookExtension.PostInteraction[len(orderbookExtension.PostInteraction)-extraDataCharacterLength:]
	extraDataBytes, err := hex.DecodeString(extraDataRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to decode escrow extension extra data: %w", err)
	}

	// Send the final 160 bytes of the postInteraction to decodeExtraData
	extraData, err := decodeExtraData(extraDataBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode escrow extension extra data: %w", err)
	}

	return &EscrowExtension{
		ExtensionPlus:  *extensionPlus,
		HashLock:         extraData.HashLock,
		DstChainId:       extraData.DstChainId,
		DstToken:         extraData.DstToken,
		SrcSafetyDeposit: fmt.Sprintf("%x", extraData.SrcSafetyDeposit),
		DstSafetyDeposit: fmt.Sprintf("%x", extraData.DstSafetyDeposit),
		TimeLocks:        *extraData.TimeLocks,
	}, nil
}

func decodeExtraData(data []byte) (*EscrowExtraData, error) {
	iter := bytesiterator.New(data)
	hashlockData, err := iter.NextUint256()
	if err != nil {
		return nil, fmt.Errorf("failed to read hashlock: %w", err)
	}

	dstChainIdData, err := iter.NextUint256()
	if err != nil {
		return nil, fmt.Errorf("failed to read destination chain ID: %w", err)
	}

	addressBig, err := iter.NextUint256()
	if err != nil {
		return nil, fmt.Errorf("failed to read address: %w", err)
	}

	addressHex := strings.ToLower(common.BigToAddress(addressBig).Hex())

	safetyDepositData, err := iter.NextUint256()
	if err != nil {
		return nil, fmt.Errorf("failed to read safety deposit data: %w", err)
	}

	// Define a 128-bit mask (2^128 - 1)
	mask := new(big.Int)
	mask.Exp(big.NewInt(2), big.NewInt(128), nil).Sub(mask, big.NewInt(1))

	srcSafetyDeposit := new(big.Int).And(safetyDepositData, mask)
	dstSafetyDeposit := new(big.Int).Rsh(safetyDepositData, 128)

	timelocksData, err := iter.NextUint256()
	if err != nil {
		return nil, fmt.Errorf("failed to read timelocks data: %w", err)
	}

	timelocks, err := decodeTimeLocks(timelocksData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode timelocks: %w", err)
	}

	return &EscrowExtraData{
		HashLock: &HashLock{
			hashlockData.String(),
		},
		DstChainId:       float32(dstChainIdData.Uint64()),
		DstToken:         common.HexToAddress(addressHex),
		SrcSafetyDeposit: srcSafetyDeposit,
		DstSafetyDeposit: dstSafetyDeposit,
		TimeLocks:        timelocks,
	}, nil
}

// decodeTimeLocks takes a *big.Int containing the raw hex data and returns a TimeLocks struct.
func decodeTimeLocks(value *big.Int) (*TimeLocks, error) {
	tl := &TimeLocks{}

	// Convert big.Int to byte slice
	data := value.Bytes()

	if len(data) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(data):], data)
		data = padded
	}

	//TODO big.Int cannot preserve leading zeroes, so decoding the deploy time is impossible atm

	// tl.DeployTime = float32(binary.BigEndian.Uint32((data[0:4])))
	tl.DstCancellation = float32(binary.BigEndian.Uint32((data[4:8])))
	tl.DstPublicWithdrawal = float32(binary.BigEndian.Uint32((data[8:12])))
	tl.DstWithdrawal = float32(binary.BigEndian.Uint32((data[12:16])))
	tl.SrcPublicCancellation = float32(binary.BigEndian.Uint32((data[16:20])))
	tl.SrcCancellation = float32(binary.BigEndian.Uint32((data[20:24])))
	tl.SrcPublicWithdrawal = float32(binary.BigEndian.Uint32((data[24:28])))
	tl.SrcWithdrawal = float32(binary.BigEndian.Uint32((data[28:32])))

	return tl, nil
}

type EscrowExtraData struct {
	HashLock         *HashLock
	DstChainId       float32
	DstToken         common.Address
	SrcSafetyDeposit *big.Int
	DstSafetyDeposit *big.Int
	TimeLocks        *TimeLocks
}

// encodeExtraData takes an EscrowExtraData struct and encodes it into a byte slice.
func encodeExtraData(data *EscrowExtraData) ([]byte, error) {
	var buffer bytes.Buffer

	// 1. Encode HashLock.Value
	hashlockData := new(big.Int)
	_, ok := hashlockData.SetString(hexadecimal.Trim0x(data.HashLock.Value), 16)
	if !ok {
		return nil, fmt.Errorf("invalid hashlock value: %s", data.HashLock.Value)
	}
	err := writeBigIntAsUint256(&buffer, hashlockData)
	if err != nil {
		return nil, err
	}

	// 2. Encode DstChainId
	dstChainIdBigInt := new(big.Int).SetUint64(uint64(data.DstChainId))
	err = writeBigIntAsUint256(&buffer, dstChainIdBigInt)
	if err != nil {
		return nil, err
	}

	// 3. Encode DstToken
	addressBig := new(big.Int).SetBytes(data.DstToken.Bytes())
	err = writeBigIntAsUint256(&buffer, addressBig)
	if err != nil {
		return nil, err
	}

	// 4. Encode SafetyDeposits
	safetyDepositData := new(big.Int)
	srcShifted := new(big.Int).Lsh(data.SrcSafetyDeposit, 128)
	safetyDepositData.Add(srcShifted, data.DstSafetyDeposit)
	err = writeBigIntAsUint256(&buffer, safetyDepositData)
	if err != nil {
		return nil, err
	}

	// 5. Encode TimeLocks
	timeLocksData, err := encodeTimeLocks(data.TimeLocks)
	if err != nil {
		return nil, err
	}
	err = writeBigIntAsUint256(&buffer, timeLocksData)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// encodeTimeLocks packs a TimeLocks struct into a *big.Int.
func encodeTimeLocks(tl *TimeLocks) (*big.Int, error) {
	data := make([]byte, 32)

	//TODO statically putting a timeDeployed value of 0 at the beginning of the encoded data for now. The data is missing from the generated struct.
	// https://github.com/1inch/cross-chain-sdk/blob/532f6ae6dc401ddaf8fe3ad040305f2500156710/src/cross-chain-order/time-locks/time-locks.ts#L33-L33
	// https://github.com/1inch/cross-chain-sdk/blob/532f6ae6dc401ddaf8fe3ad040305f2500156710/src/cross-chain-order/time-locks/time-locks.ts#L188-L188
	binary.BigEndian.PutUint32(data[0:4], uint32(0))
	binary.BigEndian.PutUint32(data[4:8], uint32(tl.DstCancellation))
	binary.BigEndian.PutUint32(data[8:12], uint32(tl.DstPublicWithdrawal))
	binary.BigEndian.PutUint32(data[12:16], uint32(tl.DstWithdrawal))
	binary.BigEndian.PutUint32(data[16:20], uint32(tl.SrcPublicCancellation))
	binary.BigEndian.PutUint32(data[20:24], uint32(tl.SrcCancellation))
	binary.BigEndian.PutUint32(data[24:28], uint32(tl.SrcPublicWithdrawal))
	binary.BigEndian.PutUint32(data[28:32], uint32(tl.SrcWithdrawal))
	timeLocksData := new(big.Int).SetBytes(data)
	return timeLocksData, nil
}

// writeBigIntAsUint256 writes a *big.Int as a 32-byte big-endian uint256 to a buffer.
func writeBigIntAsUint256(buffer *bytes.Buffer, value *big.Int) error {
	bytes := value.Bytes()
	if len(bytes) > 32 {
		return fmt.Errorf("value exceeds uint256 maximum")
	}
	// Pad with leading zeros to make it 32 bytes
	padded := make([]byte, 32)
	copy(padded[32-len(bytes):], bytes)
	_, err := buffer.Write(padded)
	return err
}
