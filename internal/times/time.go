package times

import "time"

var Now func() int64 = NowImpl

func NowImpl() int64 {
	return time.Now().Unix()
}

var CalculateAuctionStartTime func(uint32, uint32) uint32 = CalculateAuctionStartTimeImpl

func CalculateAuctionStartTimeImpl(startAuctionIn uint32, additionalWaitPeriod uint32) uint32 {
	currentTime := time.Now().Unix()
	return uint32(currentTime) + additionalWaitPeriod + startAuctionIn
}
