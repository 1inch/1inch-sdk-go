package fusionorder

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/1inch/1inch-sdk-go/internal/bytesiterator"
)

// AuctionDetails contains the auction configuration for a fusion order
type AuctionDetails struct {
	StartTime       uint32                      `json:"startTime"`
	Duration        uint32                      `json:"duration"`
	InitialRateBump uint32                      `json:"initialRateBump"`
	Points          []AuctionPointClassFixed    `json:"points"`
	GasCost         GasCostConfigClassFixed     `json:"gasCost"`
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
	if gasCost.GasBumpEstimate > Uint24Max || gasCost.GasPriceEstimate > Uint32Max ||
		startTime > Uint32Max || duration > Uint24Max || initialRateBump > Uint24Max {
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

// Encode encodes the auction details to a hex string
func (ad *AuctionDetails) Encode() string {
	bytes := make([]byte, 0)
	bytes = append(bytes, byte(ad.GasCost.GasBumpEstimate>>16), byte(ad.GasCost.GasBumpEstimate>>8), byte(ad.GasCost.GasBumpEstimate))
	bytes = append(bytes, byte(ad.GasCost.GasPriceEstimate>>24), byte(ad.GasCost.GasPriceEstimate>>16), byte(ad.GasCost.GasPriceEstimate>>8), byte(ad.GasCost.GasPriceEstimate))
	bytes = append(bytes, byte(ad.StartTime>>24), byte(ad.StartTime>>16), byte(ad.StartTime>>8), byte(ad.StartTime))
	bytes = append(bytes, byte(ad.Duration>>16), byte(ad.Duration>>8), byte(ad.Duration))
	bytes = append(bytes, byte(ad.InitialRateBump>>16), byte(ad.InitialRateBump>>8), byte(ad.InitialRateBump))
	bytes = append(bytes, byte(len(ad.Points)))

	for _, point := range ad.Points {
		bytes = append(bytes, byte(point.Coefficient>>16), byte(point.Coefficient>>8), byte(point.Coefficient))
		bytes = append(bytes, byte(point.Delay>>8), byte(point.Delay))
	}

	return hex.EncodeToString(bytes)
}

// EncodeWithoutPointCount encodes without the point count byte (used by fusionplus)
func (ad *AuctionDetails) EncodeWithoutPointCount() string {
	bytes := make([]byte, 0)
	bytes = append(bytes, byte(ad.GasCost.GasBumpEstimate>>16), byte(ad.GasCost.GasBumpEstimate>>8), byte(ad.GasCost.GasBumpEstimate))
	bytes = append(bytes, byte(ad.GasCost.GasPriceEstimate>>24), byte(ad.GasCost.GasPriceEstimate>>16), byte(ad.GasCost.GasPriceEstimate>>8), byte(ad.GasCost.GasPriceEstimate))
	bytes = append(bytes, byte(ad.StartTime>>24), byte(ad.StartTime>>16), byte(ad.StartTime>>8), byte(ad.StartTime))
	bytes = append(bytes, byte(ad.Duration>>16), byte(ad.Duration>>8), byte(ad.Duration))
	bytes = append(bytes, byte(ad.InitialRateBump>>16), byte(ad.InitialRateBump>>8), byte(ad.InitialRateBump))

	for _, point := range ad.Points {
		bytes = append(bytes, byte(point.Coefficient>>16), byte(point.Coefficient>>8), byte(point.Coefficient))
		bytes = append(bytes, byte(point.Delay>>8), byte(point.Delay))
	}

	return hex.EncodeToString(bytes)
}

// DecodeLegacyAuctionDetails decodes using the legacy format
func DecodeLegacyAuctionDetails(data string) (*AuctionDetails, error) {
	bytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("invalid hex data: %w", err)
	}

	if len(bytes) < 15 {
		return nil, errors.New("data too short: minimum 15 bytes required")
	}

	gasBumpEstimate := binary.BigEndian.Uint32(append([]byte{0x00}, bytes[0:3]...))
	gasPriceEstimate := binary.BigEndian.Uint32(bytes[3:7])
	startTime := binary.BigEndian.Uint32(bytes[7:11])
	duration := binary.BigEndian.Uint32(append([]byte{0x00}, bytes[11:14]...))
	initialRateBump := binary.BigEndian.Uint32(append([]byte{0x00}, bytes[14:17]...))

	var points []AuctionPointClassFixed
	for i := 18; i+5 <= len(bytes); i += 5 {
		points = append(points, AuctionPointClassFixed{
			Coefficient: binary.BigEndian.Uint32(append([]byte{0x00}, bytes[i:i+3]...)),
			Delay:       binary.BigEndian.Uint16(bytes[i+3 : i+5]),
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
