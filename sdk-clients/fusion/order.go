package fusion

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/1inch/1inch-sdk-go/common"
	"github.com/1inch/1inch-sdk-go/internal/times"
	geth_common "github.com/ethereum/go-ethereum/common"

	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
	"github.com/1inch/1inch-sdk-go/sdk-clients/fusionorder"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)


func CreateFusionOrderData(quote GetQuoteOutputFixed, orderParams OrderParams, wallet common.Wallet, chainId uint64) (*PreparedOrder, *orderbook.Order, error) {

	preset, err := getPreset(quote.Presets, orderParams.Preset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get preset: %w", err)
	}

	auctionDetails, err := CreateAuctionDetails(preset, orderParams.DelayAuctionStartTimeBy)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create auction details: %w", err)
	}

	takerAsset := orderParams.ToTokenAddress
	if takerAsset == fusionorder.NativeToken {
		takerAssetWrapped, ok := fusionorder.ChainToWrapper[fusionorder.NetworkEnum(chainId)]
		if !ok {
			return nil, nil, fmt.Errorf("unsupported network for wrapped token: %d", chainId)
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
	if fusionorder.IsNonceRequired(orderParams.AllowPartialFills, orderParams.AllowMultipleFills) {
		if orderParams.Nonce != nil {
			nonce = orderParams.Nonce
		} else {
			nonce, err = random_number_generation.BigIntMaxFunc(fusionorder.Uint40Max)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to generate nonce: %w", err)
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
		return nil, nil, fmt.Errorf("failed to create maker traits: %w", err)
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
		return nil, nil, fmt.Errorf("failed to generate whitelist: %w", err)
	}
	postInteractionData, err := CreateSettlementPostInteractionData(details, whitelist, orderInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create post interaction data: %w", err)
	}

	marketAmountBig := big.NewInt(0)
	_, ok := marketAmountBig.SetString(quote.MarketAmount, 10)
	if !ok {
		return nil, nil, fmt.Errorf("failed to parse market amount: %s", quote.MarketAmount)
	}

	surplusFee, err := fusionorder.FromPercent(float64(quote.SurplusFee), fusionorder.GetDefaultBase())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse surplus fee: %w", err)
	}

	surplus, err := NewSurplusParams(marketAmountBig, surplusFee)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create surplus params: %w", err)
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
		return nil, nil, fmt.Errorf("failed to create extension: %w", err)
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
		return nil, nil, fmt.Errorf("failed to create fusion order: %w", err)
	}

	orderbookExtension := fusionOrder.FusionExtension.ConvertToOrderbookExtension()
	orderbookExtensionEncoded, err := orderbookExtension.Encode()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode orderbook extension: %w", err)
	}

	salt, err := orderbook.GenerateSalt(orderbookExtensionEncoded, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate salt: %w", err)
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
		return nil, nil, fmt.Errorf("failed to create limit order message: %w", err)
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
			return nil, errors.New("custom preset requires custom preset data")
		}
		return presets.Custom, nil
	case Fast:
		return &presets.Fast, nil
	case Medium:
		return &presets.Medium, nil
	case Slow:
		return &presets.Slow, nil
	}
	return nil, fmt.Errorf("unsupported preset type: %v", presetType)
}

// CalcAuctionStartTimeFunc allows overriding the auction start time calculation for testing
var CalcAuctionStartTimeFunc func(uint32, uint32) uint32 = fusionorder.CalcAuctionStartTime

// CalcAuctionStartTime is a convenience alias for fusionorder.CalcAuctionStartTime
var CalcAuctionStartTime = fusionorder.CalcAuctionStartTime

func CreateAuctionDetails(preset *PresetClassFixed, additionalWaitPeriod float32) (*AuctionDetails, error) {
	points := make([]fusionorder.AuctionPointInput, len(preset.Points))
	for i, point := range preset.Points {
		points[i] = fusionorder.AuctionPointInput{
			Coefficient: point.Coefficient,
			Delay:       point.Delay,
		}
	}
	return fusionorder.CreateAuctionDetailsFromParams(fusionorder.CreateAuctionDetailsParams{
		StartAuctionIn:       preset.StartAuctionIn,
		AdditionalWaitPeriod: additionalWaitPeriod,
		AuctionDuration:      preset.AuctionDuration,
		InitialRateBump:      preset.InitialRateBump,
		Points:               points,
		GasCost: fusionorder.GasCostInput{
			GasBumpEstimate:  preset.GasCost.GasBumpEstimate,
			GasPriceEstimate: preset.GasCost.GasPriceEstimate,
		},
	})
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
	return fusionorder.CreateMakerTraits(fusionorder.MakerTraitsParams{
		AuctionStartTime:     details.Auction.StartTime,
		AuctionDuration:      details.Auction.Duration,
		OrderExpirationDelay: extraParams.OrderExpirationDelay,
		Nonce:                extraParams.Nonce,
		AllowPartialFills:    extraParams.AllowPartialFills,
		AllowMultipleFills:   extraParams.AllowMultipleFills,
		UnwrapWeth:           extraParams.unwrapWeth,
		EnablePermit2:        extraParams.EnablePermit2,
	})
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
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	extensionHash, err := params.Extension.Keccak256()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate extension hash: %w", err)
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
			Extension:    fmt.Sprintf("%x", extensionHash),
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

// bpsToRatioFormat is an alias for fusionorder.BpsToRatioFormat
var bpsToRatioFormat = fusionorder.BpsToRatioFormat
