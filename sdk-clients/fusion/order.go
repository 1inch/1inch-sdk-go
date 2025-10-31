package fusion

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/1inch/1inch-sdk-go/common"
	"github.com/1inch/1inch-sdk-go/internal/times"
	geth_common "github.com/ethereum/go-ethereum/common"

	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

var uint40Max = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 40), big.NewInt(1))

func CreateFusionOrderData(quote GetQuoteOutputFixed, orderParams OrderParams, wallet common.Wallet, chainId uint64) (*PreparedOrder, *orderbook.Order, error) {

	preset, err := getPreset(quote.Presets, orderParams.Preset)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting preset: %v", err)
	}

	auctionDetails, err := CreateAuctionDetails(preset, orderParams.DelayAuctionStartTimeBy)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating auction details: %v", err)
	}

	takerAsset := orderParams.ToTokenAddress
	if takerAsset == NativeToken {
		takerAssetWrapped, ok := chainToWrapper[NetworkEnum(chainId)]
		if !ok {
			return nil, nil, fmt.Errorf("unable to get address for taker asset's wrapped token. unrecognized network: %v", chainId)
		}
		takerAsset = takerAssetWrapped.Hex()
	}

	whitelistAddresses := make([]AuctionWhitelistItem, 0)
	whitelistAddressesStrings := make([]string, 0)
	for _, address := range quote.Whitelist {
		whitelistAddresses = append(whitelistAddresses, AuctionWhitelistItem{
			Address:   geth_common.HexToAddress(address),
			AllowFrom: big.NewInt(0), // TODO generating the correct list here requires checking for an exclusive resolver. This needs to be checked for later. The generated object does not see exclusive resolver correctly
		})
		whitelistAddressesStrings = append(whitelistAddressesStrings, address)
	}

	var nonce *big.Int
	if isNonceRequired(orderParams.AllowPartialFills, orderParams.AllowMultipleFills) {
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

	details := Details{
		Auction: auctionDetails,
		FeesIntAndRes: &FeesIntegratorAndResolver{
			Resolver:   ResolverFee{},
			Integrator: IntegratorFee{},
		},
		Whitelist:          whitelistAddresses,
		ResolvingStartTime: big.NewInt(int64(auctionDetails.StartTime)),
	}
	extraParams := ExtraParams{
		Nonce:                nonce,
		Permit:               "",
		AllowPartialFills:    orderParams.AllowPartialFills,
		AllowMultipleFills:   orderParams.AllowMultipleFills,
		OrderExpirationDelay: orderParams.OrderExpirationDelay,
		Source:               "",
	}

	makerTraits, err := CreateMakerTraits(details, extraParams)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating maker traits: %v", err)
	}

	orderInfo := FusionOrderV4{
		Maker:        orderParams.WalletAddress,
		MakerAsset:   orderParams.FromTokenAddress,
		MakingAmount: orderParams.Amount,
		Receiver:     orderParams.Receiver,
		TakerAsset:   takerAsset,
		TakingAmount: preset.AuctionEndAmount,
	}

	whitelist, err := GenerateWhitelist(whitelistAddressesStrings, details.ResolvingStartTime)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating whitelist: %v", err)
	}
	postInteractionData, err := CreateSettlementPostInteractionData(details, whitelist, orderInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating post interaction data: %v", err)
	}

	marketAmountBig := big.NewInt(0)
	_, ok := marketAmountBig.SetString(quote.MarketAmount, 10)
	if !ok {
		return nil, nil, fmt.Errorf("error parsing marketAmount: %v", quote.MarketAmount)
	}

	surplus, err := NewSurplusParams(marketAmountBig, FromPercent(float64(quote.SurplusFee), GetDefaultBase()))
	if err != nil {
		return nil, nil, fmt.Errorf("error creating surplus: %v", err)
	}

	extension, err := NewExtension(ExtensionParams{
		SettlementContract:  quote.SettlementAddress,
		AuctionDetails:      auctionDetails,
		PostInteractionData: postInteractionData,
		Asset:               orderInfo.MakerAsset,
		Permit:              extraParams.Permit,
		ResolvingStartTime:  details.ResolvingStartTime,
		Surplus:             surplus,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error creating extension: %v", err)
	}

	fusionOrder, err := CreateOrder(CreateOrderDataParams{
		SettlementAddress:   quote.SettlementAddress,
		PostInteractionData: postInteractionData,
		Extension:           extension,
		orderInfo:           orderInfo,
		Details:             details,
		ExtraParams:         extraParams,
		MakerTraits:         makerTraits,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error creating fusion order: %v", err)
	}

	orderbookExtension := fusionOrder.FusionExtension.ConvertToOrderbookExtension()
	orderbookExtensionEncoded, err := orderbookExtension.Encode()
	if err != nil {
		return nil, nil, fmt.Errorf("error encoding orderbookExtension: %v", err)
	}

	salt, err := orderbook.GenerateSalt(orderbookExtensionEncoded, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating salt for orderbook: %v", err)
	}

	limitOrder, err := orderbook.CreateLimitOrderMessage(orderbook.CreateOrderParams{
		Wallet:           wallet,
		MakerTraits:      makerTraits,
		Extension:        *orderbookExtension,
		ExtensionEncoded: orderbookExtensionEncoded,
		Salt:             salt,
		Maker:            fusionOrder.OrderInfo.Maker,
		MakerAsset:       fusionOrder.OrderInfo.MakerAsset,
		TakerAsset:       fusionOrder.OrderInfo.TakerAsset,
		TakingAmount:     fusionOrder.OrderInfo.TakingAmount,
		MakingAmount:     fusionOrder.OrderInfo.MakingAmount,
		Taker:            fusionOrder.OrderInfo.Receiver,
	}, int(chainId))
	if err != nil {
		return nil, nil, fmt.Errorf("error creating limit order message: %v", err)
	}

	return &PreparedOrder{
		Order:   *fusionOrder,
		Hash:    limitOrder.OrderHash,
		QuoteId: quote.QuoteId,
	}, limitOrder, nil
}

func getPreset(presets QuotePresetsClassFixed, presetType GetQuoteOutputRecommendedPreset) (*PresetClassFixed, error) {
	switch presetType {
	case Custom:
		if presets.Custom == nil {
			return nil, errors.New("custom preset selected, but no custom preset data provided")
		}
		return presets.Custom, nil
	case Fast:
		return &presets.Fast, nil
	case Medium:
		return &presets.Medium, nil
	case Slow:
		return &presets.Slow, nil
	}
	return nil, fmt.Errorf("unknown preset type: %v", presetType)
}

var CalcAuctionStartTimeFunc func(uint32, uint32) uint32 = CalcAuctionStartTime

func CalcAuctionStartTime(startAuctionIn uint32, additionalWaitPeriod uint32) uint32 {
	currentTime := time.Now().Unix()
	return uint32(currentTime) + additionalWaitPeriod + startAuctionIn
}

func CreateAuctionDetails(preset *PresetClassFixed, additionalWaitPeriod float32) (*AuctionDetails, error) {
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

func CreateSettlementPostInteractionData(details Details, whitelist []WhitelistItem, orderInfo FusionOrderV4) (*SettlementPostInteractionData, error) {
	resolverStartTime := details.ResolvingStartTime
	if details.ResolvingStartTime == nil || details.ResolvingStartTime.Cmp(big.NewInt(0)) == 0 {
		resolverStartTime = big.NewInt(times.Now())
	}

	return &SettlementPostInteractionData{
		Whitelist:          whitelist,
		AuctionFees:        details.FeesIntAndRes,
		ResolvingStartTime: resolverStartTime,
		CustomReceiver:     geth_common.HexToAddress(orderInfo.Receiver),
	}, nil
}

func CreateMakerTraits(details Details, extraParams ExtraParams) (*orderbook.MakerTraits, error) {
	deadline := details.Auction.StartTime + details.Auction.Duration + extraParams.OrderExpirationDelay
	var nonce int64
	if extraParams.Nonce == nil {
		nonce = 0
	} else {
		nonce = extraParams.Nonce.Int64()
	}
	makerTraitParms := orderbook.MakerTraitsParams{
		Expiry:             int64(deadline),
		AllowPartialFills:  extraParams.AllowPartialFills,
		AllowMultipleFills: extraParams.AllowMultipleFills,
		HasPostInteraction: true,
		UnwrapWeth:         extraParams.unwrapWeth,
		UsePermit2:         extraParams.EnablePermit2,
		HasExtension:       true,
		Nonce:              nonce,
	}
	makerTraits, err := orderbook.NewMakerTraits(makerTraitParms)
	if err != nil {
		return nil, fmt.Errorf("error creating maker traits: %v", err)
	}
	if makerTraits.IsBitInvalidatorMode() {
		if extraParams.Nonce == nil || extraParams.Nonce.Cmp(big.NewInt(0)) == 0 {
			return nil, errors.New("nonce required when partial fill or multiple fill disallowed")
		}
	}
	return makerTraits, nil
}

type CreateOrderDataParams struct {
	SettlementAddress   string
	PostInteractionData *SettlementPostInteractionData
	Extension           *Extension
	orderInfo           FusionOrderV4
	Details             Details
	ExtraParams         ExtraParams
	MakerTraits         *orderbook.MakerTraits
}

func getReceiver(fees *FeesIntegratorAndResolver, settlementAddress string, receiver string) string {
	if fees != nil {
		return settlementAddress
	}
	return receiver
}

func CreateOrder(params CreateOrderDataParams) (*Order, error) {
	receiver := getReceiver(params.Details.FeesIntAndRes, params.SettlementAddress, params.orderInfo.Receiver)

	salt, err := params.Extension.GenerateSalt()
	if err != nil {
		return nil, fmt.Errorf("error generating salt: %v", err)
	}

	return &Order{
		FusionExtension: params.Extension,
		Inner: orderbook.OrderData{
			MakerAsset:   params.orderInfo.MakerAsset,
			TakerAsset:   params.orderInfo.TakerAsset,
			MakingAmount: params.orderInfo.MakingAmount,
			TakingAmount: params.orderInfo.TakingAmount,
			Salt:         fmt.Sprintf("%x", salt),
			Maker:        params.orderInfo.Maker,
			Receiver:     receiver,
			MakerTraits:  params.MakerTraits.Encode(),
			Extension:    fmt.Sprintf("%x", params.Extension.Keccak256()),
		},
		SettlementExtension: geth_common.HexToAddress(params.SettlementAddress),
		OrderInfo:           params.orderInfo,
		AuctionDetails:      params.Details.Auction,
		PostInteractionData: params.PostInteractionData,
		Extra: ExtraData{
			UnwrapWETH:           params.ExtraParams.unwrapWeth,
			Nonce:                params.ExtraParams.Nonce,
			Permit:               params.ExtraParams.Permit,
			AllowPartialFills:    params.ExtraParams.AllowPartialFills,
			AllowMultipleFills:   params.ExtraParams.AllowMultipleFills,
			OrderExpirationDelay: params.ExtraParams.OrderExpirationDelay,
			EnablePermit2:        params.ExtraParams.EnablePermit2,
			Source:               params.ExtraParams.Source,
		},
	}, nil
}

func isNonceRequired(allowPartialFills, allowMultipleFills bool) bool {
	return !allowPartialFills || !allowMultipleFills
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
