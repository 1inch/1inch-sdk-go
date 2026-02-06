package fusionorder

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// CanExecuteAt checks if an executor can execute at the given time based on whitelist rules.
// This is a shared helper used by SettlementPostInteractionData types in fusion and fusionplus.
func CanExecuteAt(whitelist []WhitelistItem, resolvingStartTime *big.Int, executor common.Address, executionTime *big.Int) bool {
	// Whitelist AddressHalf is stored in lowercase, so we need to lowercase for comparison
	addressHalf := strings.ToLower(executor.Hex()[len(executor.Hex())-20:])

	allowedFrom := new(big.Int).Set(resolvingStartTime)

	for _, wl := range whitelist {
		allowedFrom.Add(allowedFrom, wl.Delay)

		if addressHalf == wl.AddressHalf {
			return executionTime.Cmp(allowedFrom) >= 0
		} else if executionTime.Cmp(allowedFrom) < 0 {
			return false
		}
	}

	return false
}

// IsExclusiveResolver checks if a wallet is an exclusive resolver based on whitelist rules.
// This is a shared helper used by SettlementPostInteractionData types in fusion and fusionplus.
func IsExclusiveResolver(whitelist []WhitelistItem, wallet common.Address) bool {
	// Whitelist AddressHalf is stored in lowercase, so we need to lowercase for comparison
	addressHalf := strings.ToLower(wallet.Hex()[len(wallet.Hex())-20:])

	if len(whitelist) == 0 {
		return false
	}

	if len(whitelist) == 1 {
		return addressHalf == whitelist[0].AddressHalf
	}

	if whitelist[0].Delay.Cmp(whitelist[1].Delay) == 0 {
		return false
	}

	return addressHalf == whitelist[0].AddressHalf
}
