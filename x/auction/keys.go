package auction

import "cosmossdk.io/collections"

const (
	ModuleName = "auction"

	// AuctionBurnModuleAccountName is the name of the auction burn module account.
	AuctionBurnModuleAccountName = "auction_burn"
)

// Store prefixes
var (
	// ParamsKey is the prefix for params key
	ParamsKeyPrefix = collections.NewPrefix(0)

	AuctionsKeyPrefix       = collections.NewPrefix(1)
	AuctionOwnerIndexPrefix = collections.NewPrefix(2)
)
