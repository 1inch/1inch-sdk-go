package fusionorder

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/internal/bytesbuilder"
	"github.com/1inch/1inch-sdk-go/internal/bytesiterator"
)

// AuctionDetails contains the auction configuration for a fusion order
type AuctionDetails struct {
	StartTime       uint32                   `json:"startTime"`
	Duration        uint32                   `json:"duration"`
	InitialRateBump uint32                   `json:"initialRateBump"`
	Points          []AuctionPointClassFixed `json:"points"`
	GasCost         GasCostConfigClassFixed  `json:"gasCost"`
}

// AuctionPointClassFixed represents a point in the auction curve
type AuctionPointClassFixed struct {
	Coefficient uint32 `json:"coefficient"`
	Delay       uint16 `json:"delay"`
}

// GasCostConfigClassFixed contains gas cost estimation parameters
type GasCostConfigClassFixed struct {
	GasBumpEstimate  uint32 `json:"gasBumpEstimate"`
	GasPriceEstimate uint32 `json:"gasPriceEstimate"`
}

// NewAuctionDetails creates validated AuctionDetails
func NewAuctionDetails(startTime, duration, initialRateBump uint32, points []AuctionPointClassFixed, gasCost GasCostConfigClassFixed) (*AuctionDetails, error) {
	if gasCost.GasBumpEstimate > constants.Uint24Max || gasCost.GasPriceEstimate > constants.Uint32Max ||
		startTime > constants.Uint32Max || duration > constants.Uint24Max || initialRateBump > constants.Uint24Max {
		return nil, errors.New("values exceed limits")
	}

	return &AuctionDetails{
		StartTime:       startTime,
		Duration:        duration,
		InitialRateBump: initialRateBump,
		Points:          points,
		GasCost:         gasCost,
	}, nil
}

// CalcAuctionStartTime calculates the auction start time based on current time and delays
func CalcAuctionStartTime(startAuctionIn uint32, additionalWaitPeriod uint32) uint32 {
	currentTime := time.Now().Unix()
	return uint32(currentTime) + additionalWaitPeriod + startAuctionIn
}

// DecodeAuctionDetails decodes hex-encoded auction details
func DecodeAuctionDetails(data string) (*AuctionDetails, error) {
	rawBytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("invalid hex data: %w", err)
	}

	if len(rawBytes) < 17 {
		return nil, errors.New("data too short: minimum 17 bytes required")
	}

	iter := bytesiterator.New(rawBytes)

	gasBumpEstimate, err := iter.NextUint24()
	if err != nil {
		return nil, fmt.Errorf("failed to read gas bump estimate: %w", err)
	}

	gasPriceEstimateBI, err := iter.NextUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read gas price estimate: %w", err)
	}
	gasPriceEstimate := uint32(gasPriceEstimateBI.Uint64())

	startTimeBI, err := iter.NextUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read start time: %w", err)
	}
	startTime := uint32(startTimeBI.Uint64())

	duration, err := iter.NextUint24()
	if err != nil {
		return nil, fmt.Errorf("failed to read duration: %w", err)
	}

	initialRateBump, err := iter.NextUint24()
	if err != nil {
		return nil, fmt.Errorf("failed to read initial rate bump: %w", err)
	}

	var points []AuctionPointClassFixed
	for !iter.IsEmpty() {
		if iter.BytesLeft() < 5 {
			return nil, errors.New("insufficient bytes to read next auction point")
		}

		coeff, err := iter.NextUint24()
		if err != nil {
			return nil, fmt.Errorf("failed to read coefficient in points: %w", err)
		}

		delayBI, err := iter.NextUint16()
		if err != nil {
			return nil, fmt.Errorf("failed to read delay in points: %w", err)
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

// encodeHeader writes the common auction header fields to the builder.
func (ad *AuctionDetails) encodeHeader(b *bytesbuilder.BytesBuilder) {
	b.AddNativeUint24(ad.GasCost.GasBumpEstimate)
	b.AddNativeUint32(ad.GasCost.GasPriceEstimate)
	b.AddNativeUint32(ad.StartTime)
	b.AddNativeUint24(ad.Duration)
	b.AddNativeUint24(ad.InitialRateBump)
}

// encodePoints writes auction curve points to the builder.
func (ad *AuctionDetails) encodePoints(b *bytesbuilder.BytesBuilder) {
	for _, point := range ad.Points {
		b.AddNativeUint24(point.Coefficient)
		b.AddNativeUint16(point.Delay)
	}
}

// Encode encodes the auction details to a hex string
func (ad *AuctionDetails) Encode() string {
	b := bytesbuilder.New()
	ad.encodeHeader(b)
	b.AddUint8(uint8(len(ad.Points)))
	ad.encodePoints(b)
	return b.AsHex()
}

// EncodeWithoutPointCount encodes without the point count byte (used by fusionplus)
func (ad *AuctionDetails) EncodeWithoutPointCount() string {
	b := bytesbuilder.New()
	ad.encodeHeader(b)
	ad.encodePoints(b)
	return b.AsHex()
}

// DecodeLegacyAuctionDetails decodes using the legacy format (includes point count byte)
func DecodeLegacyAuctionDetails(data string) (*AuctionDetails, error) {
	rawBytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("invalid hex data: %w", err)
	}

	if len(rawBytes) < 18 {
		return nil, errors.New("data too short: minimum 18 bytes required")
	}

	iter := bytesiterator.New(rawBytes)

	gasBumpEstimate, err := iter.NextUint24()
	if err != nil {
		return nil, fmt.Errorf("failed to read gas bump estimate: %w", err)
	}

	gasPriceEstimateBI, err := iter.NextUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read gas price estimate: %w", err)
	}
	gasPriceEstimate := uint32(gasPriceEstimateBI.Uint64())

	startTimeBI, err := iter.NextUint32()
	if err != nil {
		return nil, fmt.Errorf("failed to read start time: %w", err)
	}
	startTime := uint32(startTimeBI.Uint64())

	duration, err := iter.NextUint24()
	if err != nil {
		return nil, fmt.Errorf("failed to read duration: %w", err)
	}

	initialRateBump, err := iter.NextUint24()
	if err != nil {
		return nil, fmt.Errorf("failed to read initial rate bump: %w", err)
	}

	// Legacy format includes a point count byte
	pointCount, err := iter.NextByte()
	if err != nil {
		return nil, fmt.Errorf("failed to read point count: %w", err)
	}

	var points []AuctionPointClassFixed
	for i := 0; i < int(pointCount); i++ {
		coeff, err := iter.NextUint24()
		if err != nil {
			return nil, fmt.Errorf("failed to read coefficient in points: %w", err)
		}

		delayBI, err := iter.NextUint16()
		if err != nil {
			return nil, fmt.Errorf("failed to read delay in points: %w", err)
		}
		delay := uint16(delayBI.Uint64())

		points = append(points, AuctionPointClassFixed{
			Coefficient: coeff,
			Delay:       delay,
		})
	}

	return NewAuctionDetails(startTime, duration, initialRateBump, points, GasCostConfigClassFixed{
		GasBumpEstimate:  gasBumpEstimate,
		GasPriceEstimate: gasPriceEstimate,
	})
}

// AuctionPointInput represents an auction point from the API quote response
type AuctionPointInput struct {
	Coefficient float32
	Delay       float32
}

// GasCostInput represents gas cost configuration from the API quote response
type GasCostInput struct {
	GasBumpEstimate  float32
	GasPriceEstimate string
}

// CreateAuctionDetailsParams contains parameters to create AuctionDetails from API response
type CreateAuctionDetailsParams struct {
	StartAuctionIn       float32
	AdditionalWaitPeriod float32
	AuctionDuration      float32
	InitialRateBump      float32
	Points               []AuctionPointInput
	GasCost              GasCostInput
}

// CalcAuctionStartTimeFunc allows overriding the auction start time calculation for testing
var CalcAuctionStartTimeFunc func(uint32, uint32) uint32 = CalcAuctionStartTime

// CreateAuctionDetailsFromParams creates AuctionDetails from API response parameters
// This is shared between fusion and fusionplus packages
func CreateAuctionDetailsFromParams(params CreateAuctionDetailsParams) (*AuctionDetails, error) {
	pointsFixed := make([]AuctionPointClassFixed, 0)
	for _, point := range params.Points {
		pointsFixed = append(pointsFixed, AuctionPointClassFixed{
			Coefficient: uint32(point.Coefficient),
			Delay:       uint16(point.Delay),
		})
	}

	gasPriceEstimateFixed, err := strconv.ParseUint(params.GasCost.GasPriceEstimate, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gas price estimate: %w", err)
	}

	gasCostFixed := GasCostConfigClassFixed{
		GasBumpEstimate:  uint32(params.GasCost.GasBumpEstimate),
		GasPriceEstimate: uint32(gasPriceEstimateFixed),
	}

	return &AuctionDetails{
		StartTime:       CalcAuctionStartTimeFunc(uint32(params.StartAuctionIn), uint32(params.AdditionalWaitPeriod)),
		Duration:        uint32(params.AuctionDuration),
		InitialRateBump: uint32(params.InitialRateBump),
		Points:          pointsFixed,
		GasCost:         gasCostFixed,
	}, nil
}
