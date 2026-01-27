package times

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNowImpl(t *testing.T) {
	before := time.Now().Unix()
	result := NowImpl()
	after := time.Now().Unix()

	// Result should be between before and after
	assert.GreaterOrEqual(t, result, before)
	assert.LessOrEqual(t, result, after)
}

func TestNow_CanBeMocked(t *testing.T) {
	// Save the original function
	originalNow := Now
	defer func() { Now = originalNow }()

	// Mock the function
	mockedTime := int64(1234567890)
	Now = func() int64 {
		return mockedTime
	}

	result := Now()
	assert.Equal(t, mockedTime, result)
}

func TestCalculateAuctionStartTimeImpl(t *testing.T) {
	tests := []struct {
		name                 string
		startAuctionIn       uint32
		additionalWaitPeriod uint32
	}{
		{
			name:                 "Zero values",
			startAuctionIn:       0,
			additionalWaitPeriod: 0,
		},
		{
			name:                 "Non-zero startAuctionIn",
			startAuctionIn:       60,
			additionalWaitPeriod: 0,
		},
		{
			name:                 "Non-zero additionalWaitPeriod",
			startAuctionIn:       0,
			additionalWaitPeriod: 30,
		},
		{
			name:                 "Both non-zero",
			startAuctionIn:       60,
			additionalWaitPeriod: 30,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			before := uint32(time.Now().Unix())
			result := CalculateAuctionStartTimeImpl(tc.startAuctionIn, tc.additionalWaitPeriod)
			after := uint32(time.Now().Unix())

			expectedMin := before + tc.startAuctionIn + tc.additionalWaitPeriod
			expectedMax := after + tc.startAuctionIn + tc.additionalWaitPeriod

			assert.GreaterOrEqual(t, result, expectedMin)
			assert.LessOrEqual(t, result, expectedMax)
		})
	}
}

func TestCalculateAuctionStartTime_CanBeMocked(t *testing.T) {
	// Save the original function
	originalFunc := CalculateAuctionStartTime
	defer func() { CalculateAuctionStartTime = originalFunc }()

	// Mock the function
	mockedResult := uint32(9999999)
	CalculateAuctionStartTime = func(startAuctionIn, additionalWaitPeriod uint32) uint32 {
		return mockedResult
	}

	result := CalculateAuctionStartTime(10, 20)
	assert.Equal(t, mockedResult, result)
}
