package fusionplus

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/1inch/1inch-sdk-go/common"
	"github.com/1inch/1inch-sdk-go/common/fusionorder"
	"github.com/1inch/1inch-sdk-go/constants"
	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
	geth_common "github.com/ethereum/go-ethereum/common"
)

func CreateFusionPlusOrderData(quoteParams QuoterControllerGetQuoteParamsFixed, quote *GetQuoteOutputFixed, orderParams OrderParams, wallet common.Wallet, chainId int) (*PreparedOrder, error) {

	// TODO preset is already gotten earlier for the secret count
	preset, err := GetPreset(quote.Presets, orderParams.Preset)
	if err != nil {
		return nil, fmt.Errorf("failed to get preset: %w", err)
	}

	auctionPointsPlus := make([]AuctionPointClass, 0)
	for _, point := range preset.Points {
		auctionPointsPlus = append(auctionPointsPlus, AuctionPointClass(point))
	}

	gasCostsPlus := GasCostConfigClass{
		GasBumpEstimate:  preset.GasCost.GasBumpEstimate,
		GasPriceEstimate: preset.GasCost.GasPriceEstimate,
	}
	presetPlus := &PresetClassFixed{
		AllowMultipleFills: preset.AllowMultipleFills,
		//ExclusiveResolver: preset.ExclusiveResolver, // TODO This is not working for fusion at the moment
		AllowPartialFills:  preset.AllowPartialFills,
		AuctionDuration:    preset.AuctionDuration,
		AuctionEndAmount:   preset.AuctionEndAmount,
		AuctionStartAmount: preset.AuctionStartAmount,
		GasCost:            gasCostsPlus,
		InitialRateBump:    preset.InitialRateBump,
		Points:             auctionPointsPlus,
		StartAuctionIn:     preset.StartAuctionIn,
	}

	auctionDetailsPlus, err := CreateAuctionDetailsPlus(presetPlus, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create auction details: %w", err)
	}

	takerAsset := quoteParams.DstTokenAddress
	if takerAsset == constants.NativeToken {
		takerAssetWrapped, ok := constants.ChainToWrapper[constants.NetworkEnum(chainId)]
		if !ok {
			return nil, fmt.Errorf("unsupported network for wrapped token: %d", chainId)
		}
		takerAsset = takerAssetWrapped.Hex()
	}

	var takingFeeReceiver geth_common.Address
	if orderParams.TakingFeeReceiver == "" {
		takingFeeReceiver = geth_common.HexToAddress(constants.ZeroAddress)
	} else {
		takingFeeReceiver = geth_common.HexToAddress(orderParams.TakingFeeReceiver)
	}

	fees := Fees{
		IntFee: IntegratorFee{
			Ratio:    fusionorder.BpsToRatioFormat(quoteParams.Fee),
			Receiver: takingFeeReceiver,
		},
		BankFee: big.NewInt(0),
	}

	whitelistAddresses := make([]fusionorder.AuctionWhitelistItem, 0)
	for _, address := range quote.Whitelist {
		whitelistAddresses = append(whitelistAddresses, fusionorder.AuctionWhitelistItem{
			Address:   geth_common.HexToAddress(address),
			AllowFrom: big.NewInt(0), // TODO generating the correct list here requires checking for an exclusive resolver. This needs to be checked for later. The generated object does not see exclusive resolver correctly
		})
	}

	var nonce *big.Int
	if fusionorder.IsNonceRequired(preset.AllowPartialFills, preset.AllowMultipleFills) {
		if orderParams.Nonce != nil {
			nonce = orderParams.Nonce
		} else {
			nonce, err = random_number_generation.BigIntMaxFunc(constants.Uint40Max)
			if err != nil {
				return nil, fmt.Errorf("failed to generate nonce: %w", err)
			}
		}
	} else {
		nonce = orderParams.Nonce
	}

	details := Details{
		Auction:   auctionDetailsPlus,
		Fees:      fees,
		Whitelist: whitelistAddresses,
	}

	extraParams := ExtraParams{
		Nonce:                nonce,
		Permit:               orderParams.Permit,
		AllowPartialFills:    preset.AllowPartialFills,
		AllowMultipleFills:   preset.AllowMultipleFills,
		OrderExpirationDelay: 0,
		Source:               "",
	}

	makerTraits, err := CreateMakerTraits(details, extraParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create maker traits: %w", err)
	}

	orderInfo := CrossChainOrderDto{
		Maker:        quoteParams.WalletAddress,
		MakerAsset:   quoteParams.SrcTokenAddress,
		MakingAmount: quoteParams.Amount,
		Receiver:     orderParams.Receiver,
		TakerAsset:   takerAsset,
		TakingAmount: preset.AuctionEndAmount,
	}

	escrowParams := EscrowExtensionParams{
		HashLock:         orderParams.HashLock,
		DstChainId:       quoteParams.DstChain,
		SrcSafetyDeposit: quote.SrcSafetyDeposit,
		DstSafetyDeposit: quote.DstSafetyDeposit,
		TimeLocks: TimeLocks{
			DstCancellation:       quote.TimeLocks.DstCancellation,
			DstPublicWithdrawal:   quote.TimeLocks.DstPublicWithdrawal,
			DstWithdrawal:         quote.TimeLocks.DstWithdrawal,
			SrcCancellation:       quote.TimeLocks.SrcCancellation,
			SrcPublicCancellation: quote.TimeLocks.SrcPublicCancellation,
			SrcPublicWithdrawal:   quote.TimeLocks.SrcPublicWithdrawal,
			SrcWithdrawal:         quote.TimeLocks.SrcWithdrawal,
		}, // TODO timelocks have many safety checks
	}

	postInteractionData, err := CreateSettlementPostInteractionData(details, orderInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to create post interaction data: %w", err)
	}

	postInteractionDataWithFees, err := CreateSettlementPostInteractionDataWithFees(details, orderInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to create post interaction data with fees: %w", err)
	}

	extension, err := NewEscrowExtension(EscrowExtensionParams{
		ExtensionParamsPlus: ExtensionParamsPlus{
			SettlementContract:  quote.SrcEscrowFactory,
			PostInteractionData: postInteractionDataWithFees,
			AuctionDetails:      auctionDetailsPlus,
			Asset:               quoteParams.SrcTokenAddress,
			Permit:              orderParams.Permit, // TODO unsure about this permit value
		},
		HashLock:         orderParams.HashLock,
		DstChainId:       quoteParams.DstChain,
		DstToken:         geth_common.HexToAddress(takerAsset),
		SrcSafetyDeposit: quote.SrcSafetyDeposit,
		DstSafetyDeposit: quote.DstSafetyDeposit,
		TimeLocks:        escrowParams.TimeLocks,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create extension: %w", err)
	}

	fusionPlusOrder, err := CreateOrder(CreateOrderDataParams{
		srcEscrowFactory:    quote.SrcEscrowFactory,
		orderInfo:           orderInfo,
		escrowParams:        escrowParams,
		details:             details,
		extraParams:         extraParams,
		extension:           extension,
		makerTraits:         makerTraits,
		postInteractionData: postInteractionData,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create fusion order: %w", err)
	}

	extensionOrderbook, err := extension.ConvertToOrderbookExtension()
	if err != nil {
		return nil, fmt.Errorf("failed to convert extension to orderbook extension: %w", err)
	}

	extensionEncoded, err := extensionOrderbook.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode extension: %w", err)
	}

	salt, err := orderbook.GenerateSalt(extensionEncoded, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	limitOrder, err := orderbook.CreateLimitOrderMessage(orderbook.CreateOrderParams{
		Wallet:           wallet,
		MakerTraits:      makerTraits,
		Extension:        *extensionOrderbook,
		ExtensionEncoded: extensionEncoded,
		Salt:             salt,
		Maker:            orderInfo.Maker,
		MakerAsset:       orderInfo.MakerAsset,
		TakerAsset:       orderInfo.TakerAsset,
		TakingAmount:     orderInfo.TakingAmount,
		MakingAmount:     orderInfo.MakingAmount,
		Taker:            orderInfo.Receiver,
	}, chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to create limit order message: %w", err)
	}

	return &PreparedOrder{
		Order:      *fusionPlusOrder,
		Hash:       limitOrder.OrderHash,
		QuoteId:    quote.QuoteId,
		LimitOrder: limitOrder,
	}, nil
}

func GetPreset(presets QuotePresets, presetType GetQuoteOutputRecommendedPreset) (*Preset, error) {
	switch presetType {
	case Custom:
		if presets.Custom == nil {
			return nil, errors.New("custom preset unavailable")
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

func CreateAuctionDetails(preset *Preset, additionalWaitPeriod float32) (*fusionorder.AuctionDetails, error) {
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

var timeNow func() int64 = GetCurrentTime

func GetCurrentTime() int64 {
	return time.Now().Unix()
}

func CreateSettlementPostInteractionData(details Details, orderInfo CrossChainOrderDto) (*SettlementPostInteractionData, error) {
	resolverStartTime := details.ResolvingStartTime
	if details.ResolvingStartTime == nil || details.ResolvingStartTime.Cmp(big.NewInt(0)) == 0 {
		resolverStartTime = big.NewInt(timeNow())
	}
	return NewSettlementPostInteractionData(SettlementSuffixData{
		Whitelist:          details.Whitelist,
		IntegratorFee:      &details.Fees.IntFee,
		BankFee:            details.Fees.BankFee,
		ResolvingStartTime: resolverStartTime,
		CustomReceiver:     geth_common.HexToAddress(orderInfo.Receiver),
	})
}

type CreateOrderDataParams struct {
	srcEscrowFactory    string
	orderInfo           CrossChainOrderDto
	escrowParams        EscrowExtensionParams
	details             Details
	extraParams         ExtraParams
	extension           *EscrowExtension
	makerTraits         *orderbook.MakerTraits
	postInteractionData *SettlementPostInteractionData
}

func CreateOrder(params CreateOrderDataParams) (*Order, error) {

	salt, err := params.extension.GenerateSalt()
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	extensionHash, err := params.extension.Keccak256()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate extension hash: %w", err)
	}

	return &Order{
		EscExtension: params.extension,
		Inner: orderbook.OrderData{
			MakerAsset:   params.orderInfo.MakerAsset,
			TakerAsset:   params.orderInfo.TakerAsset,
			MakingAmount: params.orderInfo.MakingAmount,
			TakingAmount: params.orderInfo.TakingAmount,
			Salt:         fmt.Sprintf("%x", salt),
			Maker:        params.orderInfo.Maker,
			Receiver:     params.orderInfo.Receiver,
			MakerTraits:  params.makerTraits.Encode(),
			Extension:    fmt.Sprintf("%x", extensionHash),
		},
		OrderInfo:           params.orderInfo,
		AuctionDetails:      params.details.Auction,
		PostInteractionData: params.postInteractionData,
		Extra: fusionorder.ExtraData{
			UnwrapWETH:           params.extraParams.unwrapWeth,
			Nonce:                params.extraParams.Nonce,
			Permit:               params.extraParams.Permit,
			AllowPartialFills:    params.extraParams.AllowPartialFills,
			AllowMultipleFills:   params.extraParams.AllowMultipleFills,
			OrderExpirationDelay: params.extraParams.OrderExpirationDelay,
			EnablePermit2:        params.extraParams.EnablePermit2,
			Source:               params.extraParams.Source,
		},
	}, nil
}

func CreateAuctionDetailsPlus(preset *PresetClassFixed, additionalWaitPeriod float32) (*fusionorder.AuctionDetails, error) {
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

// CreateMakerTraits creates MakerTraits from Details and ExtraParams
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

// CreateSettlementPostInteractionDataWithFees creates settlement post interaction data with fee information
// for use in extension encoding.
func CreateSettlementPostInteractionDataWithFees(details Details, orderInfo CrossChainOrderDto) (*SettlementPostInteractionData, error) {
	resolverStartTime := details.ResolvingStartTime
	if details.ResolvingStartTime == nil || details.ResolvingStartTime.Cmp(big.NewInt(0)) == 0 {
		resolverStartTime = big.NewInt(timeNow())
	}
	return NewSettlementPostInteractionDataWithFees(SettlementSuffixData{
		Whitelist:          details.Whitelist,
		IntegratorFee:      &details.Fees.IntFee,
		BankFee:            details.Fees.BankFee,
		ResolvingStartTime: resolverStartTime,
		CustomReceiver:     geth_common.HexToAddress(orderInfo.Receiver),
	})
}
