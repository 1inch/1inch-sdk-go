package fusionplus

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/1inch/1inch-sdk-go/internal/hexadecimal"
	"github.com/1inch/1inch-sdk-go/sdk-clients/fusion"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
	"github.com/ethereum/go-ethereum/common"
)

type EscrowExtension struct {
	fusion.Extension
	HashLock         *HashLock
	DstChainId       float32
	DstToken         common.Address
	SrcSafetyDeposit string
	DstSafetyDeposit string
	TimeLocks        TimeLocks
}

func NewEscrowExtension(escrowParams EscrowExtensionParams) (*EscrowExtension, error) {

	extension, err := fusion.NewExtension(escrowParams.ExtensionParams)
	if err != nil {
		return nil, err
	}

	escrowExtension := &EscrowExtension{
		Extension:        *extension,
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
		return nil, fmt.Errorf("invalid hexadecimal string for source safety deposit: %v", e.SrcSafetyDeposit)
	}

	dstSafetyDepositBig := new(big.Int)
	_, ok = dstSafetyDepositBig.SetString(e.DstSafetyDeposit, 10)
	if !ok {
		return nil, fmt.Errorf("invalid hexadecimal string for destination safety deposit: %v", e.DstSafetyDeposit)
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
		return nil, fmt.Errorf("failed to encode extra data: %v", err) // TODO handle
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
		return nil, fmt.Errorf("invalid HashLock value: %s", data.HashLock.Value)
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
		return fmt.Errorf("value too large to fit in uint256")
	}
	// Pad with leading zeros to make it 32 bytes
	padded := make([]byte, 32)
	copy(padded[32-len(bytes):], bytes)
	_, err := buffer.Write(padded)
	return err
}
