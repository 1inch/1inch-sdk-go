package fusionorder

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// TakingFeeInfo contains fee information for order taking
type TakingFeeInfo struct {
	TakingFeeBps      *big.Int // 100 == 1%
	TakingFeeReceiver common.Address
}

// IntegratorFeeRatio represents an integrator fee with ratio format (used in cross-chain)
type IntegratorFeeRatio struct {
	Ratio    *big.Int
	Receiver common.Address
}

// FeesWithBankFee contains integrator fee and bank fee for cross-chain orders
type FeesWithBankFee struct {
	IntFee  IntegratorFeeRatio
	BankFee *big.Int
}

var (
	feeBase          = big.NewInt(100_000)
	bpsBase          = big.NewInt(10_000)
	bpsToRatioNumber = new(big.Int).Div(feeBase, bpsBase)
)

// BpsToRatioFormat converts basis points to ratio format
func BpsToRatioFormat(bps *big.Int) *big.Int {
	if bps == nil || bps.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0)
	}
	return new(big.Int).Mul(bps, bpsToRatioNumber)
}
