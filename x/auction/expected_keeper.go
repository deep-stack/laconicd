package auction

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AuctionUsageKeeper keep track of auction usage in other modules.
// Used to, for example, prevent deletion of a auction that's in use.
type AuctionUsageKeeper interface {
	ModuleName() string
	UsesAuction(ctx sdk.Context, auctionId string) bool

	OnAuctionWinnerSelected(ctx sdk.Context, auctionId string)
}

// AuctionHooksWrapper is a wrapper for modules to inject AuctionUsageKeeper using depinject.
// Reference: https://github.com/cosmos/cosmos-sdk/tree/v0.50.3/core/appmodule#resolving-circular-dependencies
type AuctionHooksWrapper struct{ AuctionUsageKeeper }

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AuctionHooksWrapper) IsOnePerModuleType() {}
