package fusionorder

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// WhitelistItem represents an entry in the resolver whitelist (encoded format)
type WhitelistItem struct {
	// AddressHalf is the last 10 bytes of address, no 0x prefix
	AddressHalf string
	// Delay from previous resolver in seconds
	// For first resolver delay from `resolvingStartTime`
	Delay *big.Int
}

// NewWhitelistItem creates a new WhitelistItem
func NewWhitelistItem(addressHalf string, delay *big.Int) *WhitelistItem {
	return &WhitelistItem{
		AddressHalf: addressHalf,
		Delay:       delay,
	}
}

// AuctionWhitelistItem represents a whitelist entry with full address and timestamp
type AuctionWhitelistItem struct {
	Address common.Address
	// AllowFrom is the timestamp in seconds at which this address can start resolving
	AllowFrom *big.Int
}

// GenerateWhitelist converts a list of address strings into WhitelistItems with delays
func GenerateWhitelist(addresses []string, resolvingStartTime *big.Int) ([]WhitelistItem, error) {
	if len(addresses) == 0 {
		return nil, errors.New("whitelist cannot be empty")
	}

	sumDelay := big.NewInt(0)
	whitelist := make([]WhitelistItem, len(addresses))

	for i, addr := range addresses {
		allowFrom := new(big.Int).Set(resolvingStartTime)

		zero := big.NewInt(0)
		delay := new(big.Int).Sub(allowFrom, resolvingStartTime)
		delay.Sub(delay, sumDelay)
		
		if delay.Cmp(zero) == 0 {
			delay = zero
		}
		
		whitelist[i] = WhitelistItem{
			AddressHalf: strings.ToLower(addr)[len(addr)-20:],
			Delay:       delay,
		}

		sumDelay.Add(sumDelay, whitelist[i].Delay)

		if whitelist[i].Delay.Cmp(Uint16Max) >= 0 {
			return nil, fmt.Errorf("delay too big - %d must be less than %d", whitelist[i].Delay, Uint16Max)
		}
	}

	return whitelist, nil
}

// GenerateWhitelistFromItems converts AuctionWhitelistItems into WhitelistItems with delays.
// Items are sorted by AllowFrom before processing.
func GenerateWhitelistFromItems(items []AuctionWhitelistItem, resolvingStartTime *big.Int) ([]WhitelistItem, error) {
	if len(items) == 0 {
		return nil, errors.New("whitelist cannot be empty")
	}

	// Sort by AllowFrom timestamp
	sortedItems := make([]AuctionWhitelistItem, len(items))
	copy(sortedItems, items)
	sortWhitelistByAllowFrom(sortedItems)

	sumDelay := big.NewInt(0)
	whitelist := make([]WhitelistItem, len(sortedItems))

	for i, item := range sortedItems {
		allowFrom := item.AllowFrom
		if allowFrom.Cmp(resolvingStartTime) < 0 {
			allowFrom = resolvingStartTime
		}

		zero := big.NewInt(0)
		delay := new(big.Int).Sub(allowFrom, resolvingStartTime)
		delay.Sub(delay, sumDelay)
		
		if delay.Cmp(zero) == 0 {
			delay = zero
		}
		
		whitelist[i] = WhitelistItem{
			AddressHalf: strings.ToLower(item.Address.Hex())[len(item.Address.Hex())-20:],
			Delay:       delay,
		}

		sumDelay.Add(sumDelay, whitelist[i].Delay)

		if whitelist[i].Delay.Cmp(Uint16Max) >= 0 {
			return nil, fmt.Errorf("delay too big - %d must be less than %d", whitelist[i].Delay, Uint16Max)
		}
	}

	return whitelist, nil
}

// sortWhitelistByAllowFrom sorts whitelist items by their AllowFrom timestamp
func sortWhitelistByAllowFrom(items []AuctionWhitelistItem) {
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].AllowFrom.Cmp(items[j].AllowFrom) > 0 {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}
