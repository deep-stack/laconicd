package bond

import "cosmossdk.io/collections"

const (
	ModuleName = "bond"
)

// Store prefixes
var (
	ParamsPrefix = collections.NewPrefix(0)

	BondsPrefix          = collections.NewPrefix(1)
	BondOwnerIndexPrefix = collections.NewPrefix(2)
)
