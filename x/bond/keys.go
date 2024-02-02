package bond

import "cosmossdk.io/collections"

const (
	ModuleName = "bond"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// TODO: Add required keys
)

// Store prefixes
var (
	BondsKey = collections.NewPrefix(0)
)
