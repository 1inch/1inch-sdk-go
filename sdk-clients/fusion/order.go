package fusion

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

func CreateOrder(orderParams OrderParams, quote GetQuoteOutputFixed, paramsData FusionOrderParamsData, additionalParams AdditionalParams) PreparedOrder {

	fusionParams := NewFusionOrderParams(paramsData)

	preset := getPreset(quote, fusionParams.Preset)

	auctionDetails := CreateAuctionDetails(preset, paramsData.DelayAuctionStartTimeBy)

	allowPartialFills := false
	allowMultipleFills := false
	isNonceRequired := !allowPartialFills || !allowMultipleFills

	var nonce *big.Int
	if isNonceRequired {
		if paramsData.Nonce != nil {
			nonce = paramsData.Nonce
		} else {
			nonce = randBigInt(UINT_40_MAX)
		}
	} else {
		nonce = paramsData.Nonce
	}

	takerAsset := common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE") // TODO this needs to be part of the request somehow

	var ok bool
	if takerAsset.Hex() == nativeToken {
		takerAsset, ok = CHAIN_TO_WRAPPER[NetworkEnum(paramsData.NetworkId)]
		if !ok {
			panic("Unknown network")
		}
	}

	orderInfo := FusionOrderV4{
		Maker:        additionalParams.FromAddress,
		MakerAsset:   orderParams.FromTokenAddress,
		MakerTraits:  "",
		MakingAmount: orderParams.Amount,
		Receiver:     orderParams.Receiver,
		Salt:         "",
		TakerAsset:   orderParams.ToTokenAddress,
		TakingAmount: preset.AuctionEndAmount,
	}

	//TODO this should be parsed as a big.int
	bankFee, err := BigIntFromString(preset.BankFee)
	if err != nil {
		panic("Error parsing bank fee")
	}

	fees := Fees{
		IntFee: IntegratorFee{
			Ratio:    big.NewInt(0),                                                     // TODO correctly accept ratio
			Receiver: common.HexToAddress("0x0000000000000000000000000000000000000000"), // TODO accept this from the user
		},
		BankFee: bankFee,
	}

	whitelistAddresses := make([]AuctionWhitelistItem, len(quote.Whitelist))
	for _, address := range quote.Whitelist {
		whitelistAddresses = append(whitelistAddresses, AuctionWhitelistItem{
			Address:   common.HexToAddress(address),
			AllowFrom: big.NewInt(0), // TODO generating the correct list here is more involved
		})
	}

	details := Details{
		Auction:   auctionDetails,
		Fees:      fees,
		Whitelist: whitelistAddresses,
	}

	extraParams := ExtraParams{
		Nonce:                nonce,
		Permit:               "",
		AllowPartialFills:    allowPartialFills,
		AllowMultipleFills:   allowMultipleFills,
		OrderExpirationDelay: paramsData.OrderExpirationDelay,
		EnablePermit2:        false,
		Source:               "", // TODO unsure what this is
	}

	fusionOrder, makerTraits, extension, err := CreateFusionOrder(quote.SettlementAddress, orderInfo, details, extraParams)
	if err != nil {
		panic("Error creating fusion order")
	}

	fmt.Printf("fusionOrder: %v\n", fusionOrder)

	// Add a decode makertraits function to avoid the extra return values

	limitOrder, err := orderbook.CreateLimitOrderMessage(orderbook.CreateOrderParams{
		MakerTraits:  makerTraits,
		Extension:    *extension,
		PrivateKey:   "0x7a28c1b1478581b9e1293fc1c20449e2ed3efec93c68c169d32d768c57e2d99d",
		Maker:        fusionOrder.OrderInfo.Maker,
		MakerAsset:   fusionOrder.OrderInfo.MakerAsset,
		TakerAsset:   fusionOrder.OrderInfo.TakerAsset,
		TakingAmount: fusionOrder.OrderInfo.TakingAmount,
		MakingAmount: fusionOrder.OrderInfo.MakingAmount,
		Taker:        fusionOrder.OrderInfo.Receiver, // TODO unsure if this is right
	}, 1)
	if err != nil {
		panic("Error creating limit order message")
	}

	return PreparedOrder{
		Order:   *fusionOrder,
		Hash:    limitOrder.OrderHash,
		QuoteId: quote.QuoteId,
	}
}

func BigIntFromString(s string) (*big.Int, error) {
	bigInt, ok := new(big.Int).SetString(s, 10) // base 10 for decimal
	if !ok {
		return nil, fmt.Errorf("failed to convert string (%v) to big.Int", s)
	}
	return bigInt, nil
}

// TODO fix panics
func getPreset(quote GetQuoteOutputFixed, presetType GetQuoteOutputRecommendedPreset) PresetClass {
	switch presetType {
	case Custom:
		if quote.Presets.Custom == nil {
			panic("Custom preset is not available")
		}
		return *quote.Presets.Custom
	case Fast:
		return quote.Presets.Fast
	case Medium:
		return quote.Presets.Medium
	case Slow:
		return quote.Presets.Slow
	}
	panic("Unknown preset type")
}

func CalcAuctionStartTime(startAuctionIn uint32, additionalWaitPeriod uint32) uint32 {
	currentTime := time.Now().Unix() / 1000
	return uint32(currentTime) + additionalWaitPeriod + startAuctionIn
}

func CreateAuctionDetails(preset PresetClass, additionalWaitPeriod float32) AuctionDetails {
	pointsFixed := make([]AuctionPointClassFixed, len(preset.Points))
	for _, point := range preset.Points {
		pointsFixed = append(pointsFixed, AuctionPointClassFixed{
			Coefficient: uint32(point.Coefficient),
			Delay:       uint16(point.Delay),
		})
	}

	gasPriceEstimateFixed, err := strconv.ParseUint(preset.GasCost.GasPriceEstimate, 10, 32)
	if err != nil {
		panic("Error parsing gas price estimate") //TODO get rid of panics
	}

	gasCostFixed := GasCostConfigClassFixed{
		GasBumpEstimate:  uint32(preset.GasCost.GasBumpEstimate),
		GasPriceEstimate: uint32(gasPriceEstimateFixed),
	}

	return AuctionDetails{
		StartTime:       CalcAuctionStartTime(uint32(preset.StartAuctionIn), uint32(additionalWaitPeriod)),
		Duration:        uint32(preset.AuctionDuration),
		InitialRateBump: uint32(preset.InitialRateBump),
		Points:          pointsFixed,
		GasCost:         gasCostFixed,
	}
}

const UINT_40_MAX = (1 << 40) - 1

func randBigInt(max int64) *big.Int {
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		fmt.Println("Error generating random number:", err)
		return big.NewInt(0)
	}
	return n
}

const nativeToken = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"

type NetworkEnum int

const (
	ETHEREUM  NetworkEnum = 1
	POLYGON   NetworkEnum = 137
	BINANCE   NetworkEnum = 56
	ARBITRUM  NetworkEnum = 42161
	AVALANCHE NetworkEnum = 43114
	OPTIMISM  NetworkEnum = 10
	FANTOM    NetworkEnum = 250
	GNOSIS    NetworkEnum = 100
	COINBASE  NetworkEnum = 8453
)

var CHAIN_TO_WRAPPER = map[NetworkEnum]common.Address{
	ETHEREUM:  common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
	BINANCE:   common.HexToAddress("0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c"),
	POLYGON:   common.HexToAddress("0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270"),
	ARBITRUM:  common.HexToAddress("0x82af49447d8a07e3bd95bd0d56f35241523fbab1"),
	AVALANCHE: common.HexToAddress("0xb31f66aa3c1e785363f0875a1b74e27b85fd66c7"),
	GNOSIS:    common.HexToAddress("0xe91d153e0b41518a2ce8dd3d7944fa863463a97d"),
	COINBASE:  common.HexToAddress("0x4200000000000000000000000000000000000006"),
	OPTIMISM:  common.HexToAddress("0x4200000000000000000000000000000000000006"),
	FANTOM:    common.HexToAddress("0x21be370d5312f44cb42ce377bc9b8a0cef1a4c83"),
}
