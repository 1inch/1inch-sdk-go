package fusion

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math"
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
	bytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, errors.New("invalid hex data")
	}

	if len(bytes) < 15 {
		return nil, errors.New("data too short")
	}

	gasBumpEstimate := binary.BigEndian.Uint32(append([]byte{0x00}, bytes[0:3]...))
	gasPriceEstimate := binary.BigEndian.Uint32(bytes[3:7])
	startTime := binary.BigEndian.Uint32(bytes[7:11])
	duration := binary.BigEndian.Uint32(append([]byte{0x00}, bytes[11:14]...))
	initialRateBump := binary.BigEndian.Uint32(append([]byte{0x00}, bytes[14:17]...))

	var points []AuctionPointClassFixed
	for i := 17; i+5 <= len(bytes); i += 5 {
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

func (ad AuctionDetails) Encode() string {
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
