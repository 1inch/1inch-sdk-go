package times

import "time"

var timeNow func() int64 = GetCurrentTime

func GetCurrentTime() int64 {
	return time.Now().Unix()
}

var CalcAuctionStartTimeFunc func(uint32, uint32) uint32 = CalcAuctionStartTime

func CalcAuctionStartTime(startAuctionIn uint32, additionalWaitPeriod uint32) uint32 {
	currentTime := time.Now().Unix()
	return uint32(currentTime) + additionalWaitPeriod + startAuctionIn
}
