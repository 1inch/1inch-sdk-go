package fusionplus

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/1inch/1inch-sdk-go/internal/bytesbuilder"
	"github.com/1inch/1inch-sdk-go/internal/bytesiterator"
)

const (
	uint24Max = (1 << 24) - 1
	uint32Max = math.MaxUint32
)

func NewAuctionDetails(startTime, duration, initialRateBump uint32, points []AuctionPointClassFixed, gasCost GasCostConfigClassFixed) (*AuctionDetails, error) {

	if gasCost.GasBumpEstimate > uint24Max || gasCost.GasPriceEstimate > uint32Max ||
		startTime > uint32Max || duration > uint24Max || initialRateBump > uint24Max {
		return nil, errors.New("values exceed their respective limits")
	}

	return &AuctionDetails{
		StartTime:       startTime,
		Duration:        duration,
		InitialRateBump: initialRateBump,
		Points:          points,
		GasCost:         gasCost,
	}, nil
}

func DecodeAuctionDetails(data string) (*AuctionDetails, error) {
	// Decode the hex string to bytes
	rawBytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, errors.New("invalid hex data")
	}

	// We expect at least 3 + 4 + 4 + 3 + 3 = 17 bytes for the mandatory fields
	if len(rawBytes) < 17 {
		return nil, errors.New("data too short for mandatory fields")
	}

	iter := bytesiterator.New(rawBytes)

	// 1) GasBumpEstimate (3 bytes)
	gasBumpEstimate, err := iter.NextUint24()
	if err != nil {
		return nil, fmt.Errorf("failed reading gasBumpEstimate: %w", err)
	}

	// 2) GasPriceEstimate (4 bytes)
	gasPriceEstimateBI, err := iter.NextUint32()
	if err != nil {
		return nil, fmt.Errorf("failed reading gasPriceEstimate: %w", err)
	}
	gasPriceEstimate := uint32(gasPriceEstimateBI.Uint64())

	// 3) StartTime (4 bytes)
	startTimeBI, err := iter.NextUint32()
	if err != nil {
		return nil, fmt.Errorf("failed reading startTime: %w", err)
	}
	startTime := uint32(startTimeBI.Uint64())

	// 4) Duration (3 bytes)
	duration, err := iter.NextUint24()
	if err != nil {
		return nil, fmt.Errorf("failed reading duration: %w", err)
	}

	// 5) InitialRateBump (3 bytes)
	initialRateBump, err := iter.NextUint24()
	if err != nil {
		return nil, fmt.Errorf("failed reading initialRateBump: %w", err)
	}

	// Now read points (each point is 5 bytes: 3 for Coefficient, 2 for Delay)
	var points []AuctionPointClassFixed
	for !iter.IsEmpty() {
		// Ensure we have at least 5 bytes left for the next point
		if iter.BytesLeft() < 5 {
			return nil, errors.New("insufficient bytes to read next auction point")
		}

		// Coefficient (3 bytes => uint24)
		coeff, err := iter.NextUint24()
		if err != nil {
			return nil, fmt.Errorf("failed reading Coefficient in points: %w", err)
		}

		// Delay (2 bytes => *big.Int => convert to uint16)
		delayBI, err := iter.NextUint16()
		if err != nil {
			return nil, fmt.Errorf("failed reading Delay in points: %w", err)
		}
		delay := uint16(delayBI.Uint64())

		points = append(points, AuctionPointClassFixed{
			Coefficient: coeff,
			Delay:       delay,
		})
	}

	return NewAuctionDetails(
		startTime,
		duration,
		initialRateBump,
		points,
		GasCostConfigClassFixed{
			GasBumpEstimate:  gasBumpEstimate,
			GasPriceEstimate: gasPriceEstimate,
		},
	)
}

func (ad AuctionDetails) Encode() string {
	// Create a new bytes builder
	bb := bytesbuilder.New()

	// GasBumpEstimate -> 3 bytes
	bb.AddUint24(big.NewInt(int64(ad.GasCost.GasBumpEstimate)))

	// GasPriceEstimate -> 4 bytes
	bb.AddUint32(big.NewInt(int64(ad.GasCost.GasPriceEstimate)))

	// StartTime -> 4 bytes
	bb.AddUint32(big.NewInt(int64(ad.StartTime)))

	// Duration -> 3 bytes
	bb.AddUint24(big.NewInt(int64(ad.Duration)))

	// InitialRateBump -> 3 bytes
	bb.AddUint24(big.NewInt(int64(ad.InitialRateBump)))

	// Encode each point: 3 bytes for Coefficient, 2 bytes for Delay
	for _, point := range ad.Points {
		bb.AddUint24(big.NewInt(int64(point.Coefficient)))
		bb.AddUint16(big.NewInt(int64(point.Delay)))
	}

	// Return the final hex string
	return bb.AsHex()
}
