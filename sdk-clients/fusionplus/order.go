package fusionplus

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/1inch/1inch-sdk-go/common"
	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
	"github.com/1inch/1inch-sdk-go/sdk-clients/fusion"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
	geth_common "github.com/ethereum/go-ethereum/common"
)

var uint40Max = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 40), big.NewInt(1))

func CreateFusionPlusOrderData(quoteParams QuoterControllerGetQuoteParamsFixed, quote *GetQuoteOutputFixed, orderParams OrderParams, wallet common.Wallet, chainId int) (*PreparedOrder, error) {

	// TODO preset is already gotten earlier for the secret count
	preset, err := GetPreset(quote.Presets, orderParams.Preset)
	if err != nil {
		return nil, fmt.Errorf("error getting preset: %v", err)
	}

	auctionPointsFusion := make([]fusion.AuctionPointClass, 0)
	for _, point := range preset.Points {
		auctionPointsFusion = append(auctionPointsFusion, fusion.AuctionPointClass{
			Coefficient: point.Coefficient,
			Delay:       point.Delay,
		})
	}

	gasCostsFusion := fusion.GasCostConfigClass{
		GasBumpEstimate:  preset.GasCost.GasBumpEstimate,
		GasPriceEstimate: preset.GasCost.GasPriceEstimate,
	}
	presetFusion := &fusion.PresetClassFixed{
		AllowMultipleFills: preset.AllowMultipleFills,
		//ExclusiveResolver: preset.ExclusiveResolver, // TODO This is not working for fusion at the moment
		AllowPartialFills:  preset.AllowPartialFills,
		AuctionDuration:    preset.AuctionDuration,
		AuctionEndAmount:   preset.AuctionEndAmount,
		AuctionStartAmount: preset.AuctionStartAmount,
		GasCost:            gasCostsFusion,
		InitialRateBump:    preset.InitialRateBump,
		Points:             auctionPointsFusion,
		StartAuctionIn:     preset.StartAuctionIn,
	}

	auctionDetails, err := CreateAuctionDetails(preset, 0) // No extra delay for now
	if err != nil {
		return nil, fmt.Errorf("error creating auction details: %v", err)
	}

	auctionDetailsFusion, err := fusion.CreateAuctionDetails(presetFusion, 0)
	if err != nil {
		return nil, fmt.Errorf("error creating auction details: %v", err)
	}

	takerAsset := quoteParams.DstTokenAddress
	if takerAsset == NativeToken {
		takerAssetWrapped, ok := chainToWrapper[NetworkEnum(chainId)]
		if !ok {
			return nil, fmt.Errorf("unable to get address for taker asset's wrapped token. unrecognized network: %v", chainId)
		}
		takerAsset = takerAssetWrapped.Hex()
	}

	var takingFreeReceiver geth_common.Address
	if orderParams.TakingFeeReceiver == "" {
		takingFreeReceiver = geth_common.HexToAddress("0x0000000000000000000000000000000000000000")
	} else {
		takingFreeReceiver = geth_common.HexToAddress(orderParams.TakingFeeReceiver)
	}

	fees := Fees{
		IntFee: IntegratorFee{
			Ratio:    bpsToRatioFormat(quoteParams.Fee),
			Receiver: takingFreeReceiver,
		},
		BankFee: big.NewInt(0),
	}
	feesFusion := fusion.Fees{
		IntFee: fusion.IntegratorFee{
			Ratio:    bpsToRatioFormat(quoteParams.Fee),
			Receiver: takingFreeReceiver,
		},
		BankFee: big.NewInt(0),
	}

	whitelistAddresses := make([]AuctionWhitelistItem, 0)
	for _, address := range quote.Whitelist {
		whitelistAddresses = append(whitelistAddresses, AuctionWhitelistItem{
			Address:   geth_common.HexToAddress(address),
			AllowFrom: big.NewInt(0), // TODO generating the correct list here requires checking for an exclusive resolver. This needs to be checked for later. The generated object does not see exclusive resolver correctly
		})
	}
	whitelistAddressesFusion := make([]fusion.AuctionWhitelistItem, 0)
	for _, address := range quote.Whitelist {
		whitelistAddressesFusion = append(whitelistAddressesFusion, fusion.AuctionWhitelistItem{
			Address:   geth_common.HexToAddress(address),
			AllowFrom: big.NewInt(0), // TODO generating the correct list here requires checking for an exclusive resolver. This needs to be checked for later. The generated object does not see exclusive resolver correctly
		})
	}

	var nonce *big.Int
	if isNonceRequired(preset.AllowPartialFills, preset.AllowMultipleFills) {
		if orderParams.Nonce != nil {
			nonce = orderParams.Nonce
		} else {
			nonce, err = random_number_generation.BigIntMaxFunc(uint40Max)
			if err != nil {
				return nil, fmt.Errorf("error generating nonce: %v\n", err)
			}
		}
	} else {
		nonce = orderParams.Nonce
	}

	details := Details{
		Auction:   auctionDetails,
		Fees:      fees,
		Whitelist: whitelistAddresses,
	}
	detailsFusion := fusion.Details{
		Auction:   auctionDetailsFusion,
		Fees:      feesFusion,
		Whitelist: whitelistAddressesFusion,
	}

	extraParams := ExtraParams{
		Nonce:                nonce,
		Permit:               orderParams.Permit,
		AllowPartialFills:    preset.AllowPartialFills,
		AllowMultipleFills:   preset.AllowMultipleFills,
		OrderExpirationDelay: 0,
		Source:               "",
	}
	extraParamsFusion := fusion.ExtraParams{
		Nonce:                nonce,
		Permit:               orderParams.Permit,
		AllowPartialFills:    preset.AllowPartialFills,
		AllowMultipleFills:   preset.AllowMultipleFills,
		OrderExpirationDelay: 0,
		Source:               "",
	}

	makerTraitsFusion, err := fusion.CreateMakerTraits(detailsFusion, extraParamsFusion)
	if err != nil {
		return nil, fmt.Errorf("error creating maker traits: %v", err)
	}

	orderInfo := CrossChainOrderDto{
		Maker:        quoteParams.WalletAddress,
		MakerAsset:   quoteParams.SrcTokenAddress,
		MakingAmount: quoteParams.Amount,
		Receiver:     orderParams.Receiver,
		TakerAsset:   takerAsset,
		TakingAmount: preset.AuctionEndAmount,
	}
	orderInfoFusion := fusion.FusionOrderV4{
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
		return nil, fmt.Errorf("error creating post interaction data: %v", err)
	}
	postInteractionDataFusion, err := fusion.CreateSettlementPostInteractionData(detailsFusion, orderInfoFusion)
	if err != nil {
		return nil, fmt.Errorf("error creating post interaction data: %v", err)
	}

	extension, err := NewEscrowExtension(EscrowExtensionParams{
		ExtensionParams: fusion.ExtensionParams{
			SettlementContract:  quote.SrcEscrowFactory,
			PostInteractionData: postInteractionDataFusion,
			AuctionDetails:      auctionDetailsFusion,
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
		return nil, fmt.Errorf("error creating extension: %v", err)
	}

	fusionPlusOrder, err := CreateOrder(CreateOrderDataParams{
		srcEscrowFactory:    quote.SrcEscrowFactory,
		orderInfo:           orderInfo,
		escrowParams:        escrowParams,
		details:             details,
		extraParams:         extraParams,
		extension:           extension,
		makerTraits:         makerTraitsFusion,
		postInteractionData: postInteractionData,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating fusion order: %v", err)
	}

	extensionOrderbook, err := extension.ConvertToOrderbookExtension()
	if err != nil {
		return nil, fmt.Errorf("error converting extension to orderbook extension: %v", err)
	}

	limitOrder, err := orderbook.CreateLimitOrderMessage(orderbook.CreateOrderParams{
		Wallet:       wallet,
		MakerTraits:  makerTraitsFusion,
		Extension:    *extensionOrderbook,
		Maker:        orderInfo.Maker,
		MakerAsset:   orderInfo.MakerAsset,
		TakerAsset:   orderInfo.TakerAsset,
		TakingAmount: orderInfo.TakingAmount,
		MakingAmount: orderInfo.MakingAmount,
		Taker:        orderInfo.Receiver,
	}, chainId)
	if err != nil {
		return nil, fmt.Errorf("error creating limit order message: %v", err)
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
			return nil, errors.New("custom preset is not available")
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

func CreateAuctionDetails(preset *Preset, additionalWaitPeriod float32) (*AuctionDetails, error) {
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
		return nil, fmt.Errorf("error generating salt: %v", err)
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
			Extension:    fmt.Sprintf("%x", params.extension.Keccak256()),
		},
		OrderInfo:           params.orderInfo,
		AuctionDetails:      params.details.Auction,
		PostInteractionData: params.postInteractionData,
		Extra: ExtraData{
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
