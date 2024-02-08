package bond

import "cosmossdk.io/collections"

const (
	ModuleName = "bond"

	// StoreKey is the string store representation
	StoreKey = ModuleName
)

// Store prefixes
var (
	// ParamsKey is the prefix for params key
	ParamsKeyPrefix = collections.NewPrefix(0)

	BondsKeyPrefix       = collections.NewPrefix(1)
	BondOwnerIndexPrefix = collections.NewPrefix(2)
)
