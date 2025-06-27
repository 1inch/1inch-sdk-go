package fusion

import (
	"fmt"

	"github.com/1inch/1inch-sdk-go/internal/bigint"
)

// Validate validates the custom preset
func (q *CustomPreset) Validate() error {
	if !isValidAmount(q.AuctionStartAmount) {
		return fmt.Errorf("invalid auctionStartAmount: %s", q.AuctionStartAmount)
	}

	if !isValidAmount(q.AuctionEndAmount) {
		return fmt.Errorf("invalid auctionEndAmount: %s", q.AuctionEndAmount)
	}

	err := q.validateAuctionDuration(q.AuctionDuration)
	if err != nil {
		return err
	}

	pointsErr := q.validatePoints(
		q.Points,
		q.AuctionStartAmount,
		q.AuctionEndAmount,
	)
	if pointsErr != nil {
		return pointsErr
	}

	return nil
}

// validateAuctionDuration validates the auction duration
func (q *CustomPreset) validateAuctionDuration(duration int) error {
	if duration <= 0 {
		return fmt.Errorf("auctionDuration should be integer, got: %d", duration)
	}

	return nil
}

// validatePoints validates the points in the custom preset
func (q *CustomPreset) validatePoints(
	points []CustomPresetPoint,
	auctionStartAmount string,
	auctionEndAmount string,
) error {
	if points == nil || len(points) == 0 {
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

		// Check if the token amount is within the range of auction
		if tokenAmount.Cmp(startAmount) > 0 || tokenAmount.Cmp(endAmount) < 0 {
			return fmt.Errorf("points should be in range of auction (between %s and %s), invalid value at index %d: %s", auctionEndAmount, auctionStartAmount, i, point.ToTokenAmount)
		}
	}

	return nil
}

// isValidAmount checks if the amount is a valid big integer
func isValidAmount(amount string) bool {
	_, err := bigint.FromString(amount)
	return err == nil
}
