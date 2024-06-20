package fusion

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"

	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

var uint40Max = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 40), big.NewInt(1))

func CreateFusionOrderData(quote GetQuoteOutputFixed, orderParams OrderParams, additionalParams AdditionalParams) (*PreparedOrder, *orderbook.Order, error) {

	preset, err := getPreset(quote, orderParams.Preset)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting preset: %v", err)
	}

	auctionDetails, err := CreateAuctionDetails(preset, orderParams.DelayAuctionStartTimeBy)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating auction details: %v", err)
	}

	fmt.Printf("Auction start time: %v\n", auctionDetails.StartTime)

	allowPartialFills := orderParams.AllowPartialFills
	allowMultipleFills := orderParams.AllowMultipleFills
	isNonceRequired := !allowPartialFills || !allowMultipleFills

	var nonce *big.Int
	if isNonceRequired {
		if orderParams.Nonce != nil {
			nonce = orderParams.Nonce
		} else {
			nonce, err = random_number_generation.BigIntMaxFunc(uint40Max)
			if err != nil {
				return nil, nil, fmt.Errorf("error generating nonce: %v\n", err)
			}
		}
	} else {
		nonce = orderParams.Nonce
	}

	fmt.Printf("Nonce: %v\n", nonce)

	takerAsset := orderParams.ToTokenAddress
	if takerAsset == nativeToken {
		takerAssetWrapped, ok := chainToWrapper[NetworkEnum(additionalParams.NetworkId)]
		if !ok {
			return nil, nil, fmt.Errorf("unable to get address for taker asset's wrapped token. unrecognized network: %v", additionalParams.NetworkId)
		}
		takerAsset = takerAssetWrapped.Hex()
	}

	//TODO this should be parsed as a big.int after the generated struct types are fixed
	bankFee, err := BigIntFromString(preset.BankFee)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing bank fee: %v", err)
	}

	fees := Fees{
		IntFee: IntegratorFee{
			Ratio:    bpsToRatioFormat(orderParams.Fee.TakingFeeBps),
			Receiver: orderParams.Fee.TakingFeeReceiver,
		},
		BankFee: bankFee,
	}

	whitelistAddresses := make([]AuctionWhitelistItem, 0)
	for _, address := range quote.Whitelist {
		whitelistAddresses = append(whitelistAddresses, AuctionWhitelistItem{
			Address:   common.HexToAddress(address),
			AllowFrom: big.NewInt(0), // TODO generating the correct list here requires checking for an exclusive resolver. This needs to be checked for later. The generated object does not see exclusive resolver correctly
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
		OrderExpirationDelay: orderParams.OrderExpirationDelay,
		Source:               "", // TODO unsure what this is
	}

	makerTraits, err := CreateMakerTraits(details, extraParams)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating maker traits: %v", err)
	}

	orderInfo := FusionOrderV4{
		Maker:        additionalParams.FromAddress,
		MakerAsset:   orderParams.FromTokenAddress,
		MakingAmount: orderParams.Amount,
		Receiver:     orderParams.Receiver,
		TakerAsset:   takerAsset,
		TakingAmount: preset.AuctionEndAmount,
	}

	postInteractionData, err := CreateSettlementPostInteractionData(details, orderInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating post interaction data: %v", err)
	}

	extension, err := CreateExtension(CreateExtensionParams{
		settlementAddress:   quote.SettlementAddress,
		postInteractionData: postInteractionData,
		orderInfo:           orderInfo,
		details:             details,
		extraParams:         extraParams,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error creating extension: %v", err)
	}

	fusionOrder, err := CreateOrder(CreateOrderDataParams{
		settlementAddress:   quote.SettlementAddress,
		postInteractionData: postInteractionData,
		extension:           extension,
		orderInfo:           orderInfo,
		details:             details,
		extraParams:         extraParams,
		makerTraits:         makerTraits,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error creating fusion order: %v", err)
	}

	//fusionOrderIndented, err := json.MarshalIndent(fusionOrder, "", "  ")
	//if err != nil {
	//	panic("Error marshaling fusion order")
	//}
	//fmt.Printf("Fusion Order indented: %s\n", string(fusionOrderIndented))

	// Add a decode makertraits function to avoid the extra return values

	limitOrder, err := orderbook.CreateLimitOrderMessage(orderbook.CreateOrderParams{
		MakerTraits:  makerTraits,
		Extension:    *fusionOrder.FusionExtension.ConvertToOrderbookExtension(),
		PrivateKey:   additionalParams.PrivateKey,
		Maker:        fusionOrder.OrderInfo.Maker,
		MakerAsset:   fusionOrder.OrderInfo.MakerAsset,
		TakerAsset:   fusionOrder.OrderInfo.TakerAsset,
		TakingAmount: fusionOrder.OrderInfo.TakingAmount,
		MakingAmount: fusionOrder.OrderInfo.MakingAmount,
		Taker:        fusionOrder.OrderInfo.Receiver, // TODO unsure if this is right
	}, additionalParams.NetworkId)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating limit order message: %v", err)
	}

	return &PreparedOrder{
		Order:   *fusionOrder,
		Hash:    limitOrder.OrderHash,
		QuoteId: quote.QuoteId,
	}, limitOrder, nil
}

func BigIntFromString(s string) (*big.Int, error) {
	bigInt, ok := new(big.Int).SetString(s, 10) // base 10 for decimal
	if !ok {
		return nil, fmt.Errorf("failed to convert string (%v) to big.Int", s)
	}
	return bigInt, nil
}

func getPreset(quote GetQuoteOutputFixed, presetType GetQuoteOutputRecommendedPreset) (*PresetClass, error) {
	switch presetType {
	case Custom:
		if quote.Presets.Custom == nil {
			return nil, errors.New("custom preset is not available")
		}
		return quote.Presets.Custom, nil
	case Fast:
		return &quote.Presets.Fast, nil
	case Medium:
		return &quote.Presets.Medium, nil
	case Slow:
		return &quote.Presets.Slow, nil
	}
	panic("Unknown preset type")
}

var CalcAuctionStartTimeFunc func(uint32, uint32) uint32 = CalcAuctionStartTime

func CalcAuctionStartTime(startAuctionIn uint32, additionalWaitPeriod uint32) uint32 {
	currentTime := time.Now().Unix()
	return uint32(currentTime) + additionalWaitPeriod + startAuctionIn
}

func CreateAuctionDetails(preset *PresetClass, additionalWaitPeriod float32) (*AuctionDetails, error) {
	pointsFixed := make([]AuctionPointClassFixed, 0)
	for _, point := range preset.Points {
		pointsFixed = append(pointsFixed, AuctionPointClassFixed{
			Coefficient: uint32(point.Coefficient),
			Delay:       uint16(point.Delay),
		})
	}

	gasPriceEstimateFixed, err := strconv.ParseUint(preset.GasCost.GasPriceEstimate, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("error parsing gas price estimate: %v", err)
	}

	gasCostFixed := GasCostConfigClassFixed{
		GasBumpEstimate:  uint32(preset.GasCost.GasBumpEstimate),
		GasPriceEstimate: uint32(gasPriceEstimateFixed),
	}

	return &AuctionDetails{
		StartTime:       CalcAuctionStartTimeFunc(uint32(preset.StartAuctionIn), uint32(additionalWaitPeriod)),
		Duration:        uint32(preset.AuctionDuration),
		InitialRateBump: uint32(preset.InitialRateBump),
		Points:          pointsFixed,
		GasCost:         gasCostFixed,
	}, nil
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

var chainToWrapper = map[NetworkEnum]common.Address{
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

var (
	feeBase          = big.NewInt(100_000)
	bpsBase          = big.NewInt(10_000)
	bpsToRatioNumber = new(big.Int).Div(feeBase, bpsBase)
)

func bpsToRatioFormat(bps *big.Int) *big.Int {
	if bps == nil || bps.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0)
	}

	return bps.Mul(bps, bpsToRatioNumber)
}
