package fusionorder

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/1inch/1inch-sdk-go/internal/bigint"
)

// CustomPreset allows users to specify custom auction parameters
type CustomPreset struct {
	AuctionDuration    int                 `json:"auctionDuration"`
	AuctionStartAmount string              `json:"auctionStartAmount"`
	AuctionEndAmount   string              `json:"auctionEndAmount"`
	Points             []CustomPresetPoint `json:"points,omitempty"`
}

// CustomPresetPoint represents a point in a custom auction curve
type CustomPresetPoint struct {
	ToTokenAmount string `json:"toTokenAmount"`
	Delay         int    `json:"delay"`
}

// Validate validates the custom preset
func (q *CustomPreset) Validate() error {
	if !isValidAmount(q.AuctionStartAmount) {
		return fmt.Errorf("invalid auctionStartAmount: %s", q.AuctionStartAmount)
	}

	if !isValidAmount(q.AuctionEndAmount) {
		return fmt.Errorf("invalid auctionEndAmount: %s", q.AuctionEndAmount)
	}

	if err := q.validateAuctionDuration(q.AuctionDuration); err != nil {
		return err
	}

	if err := q.validatePoints(q.Points, q.AuctionStartAmount, q.AuctionEndAmount); err != nil {
		return err
	}

	return nil
}

func (q *CustomPreset) validateAuctionDuration(duration int) error {
	if duration <= 0 {
		return fmt.Errorf("auctionDuration should be integer, got: %d", duration)
	}
	return nil
}

func (q *CustomPreset) validatePoints(points []CustomPresetPoint, auctionStartAmount, auctionEndAmount string) error {
	if len(points) == 0 {
		return nil
	}

	startAmount, err := bigint.FromString(auctionStartAmount)
	if err != nil {
		return fmt.Errorf("invalid auctionStartAmount: %s", auctionStartAmount)
	}

	endAmount, err := bigint.FromString(auctionEndAmount)
	if err != nil {
		return fmt.Errorf("invalid auctionEndAmount: %s", auctionEndAmount)
	}

	for i, point := range points {
		tokenAmount, err := bigint.FromString(point.ToTokenAmount)
		if err != nil {
			return fmt.Errorf("points should be an array of valid amounts, invalid value at index %d: %s", i, point.ToTokenAmount)
		}

		if tokenAmount.Cmp(startAmount) > 0 || tokenAmount.Cmp(endAmount) < 0 {
			return fmt.Errorf("points should be in range of auction (between %s and %s), invalid value at index %d: %s", auctionEndAmount, auctionStartAmount, i, point.ToTokenAmount)
		}
	}

	return nil
}

func isValidAmount(amount string) bool {
	_, err := bigint.FromString(amount)
	return err == nil
}

// PresetType represents the type of preset to use
type PresetType string

const (
	PresetCustom PresetType = "custom"
	PresetFast   PresetType = "fast"
	PresetMedium PresetType = "medium"
	PresetSlow   PresetType = "slow"
)

// PresetData contains parsed preset information for auction creation
type PresetData struct {
	AuctionDuration    float32
	StartAuctionIn     float32
	InitialRateBump    float32
	AuctionStartAmount string
	AuctionEndAmount   string
	Points             []AuctionPointClassFixed
	GasCost            GasCostConfigClassFixed
	AllowPartialFills  bool
	AllowMultipleFills bool
	BankFee            string
	TokenFee           string
	ExclusiveResolver  string
}

// CreateAuctionDetailsFromPreset creates AuctionDetails from preset data
func CreateAuctionDetailsFromPreset(preset *PresetData, additionalWaitPeriod float32) (*AuctionDetails, error) {
	points := make([]AuctionPointClassFixed, len(preset.Points))
	for i, p := range preset.Points {
		points[i] = AuctionPointClassFixed{
			Coefficient: p.Coefficient,
			Delay:       p.Delay,
		}
	}

	return &AuctionDetails{
		StartTime:       CalcAuctionStartTime(uint32(preset.StartAuctionIn), uint32(additionalWaitPeriod)),
		Duration:        uint32(preset.AuctionDuration),
		InitialRateBump: uint32(preset.InitialRateBump),
		Points:          points,
		GasCost:         preset.GasCost,
	}, nil
}

// ParseGasPriceEstimate parses a gas price estimate string to uint32
func ParseGasPriceEstimate(estimate string) (uint32, error) {
	if estimate == "" {
		return 0, nil
	}
	val, err := strconv.ParseUint(estimate, 10, 32)
	if err != nil {
		return 0, errors.New("invalid gas price estimate")
	}
	return uint32(val), nil
}
