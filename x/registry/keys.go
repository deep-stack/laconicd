package registry

import "cosmossdk.io/collections"

const (
	// ModuleName is the name of the registry module
	ModuleName = "registry"

	// RecordRentModuleAccountName is the name of the module account that keeps track of record rents paid.
	RecordRentModuleAccountName = "record_rent"

	// AuthorityRentModuleAccountName is the name of the module account that keeps track of authority rents paid.
	AuthorityRentModuleAccountName = "authority_rent"
)

// Store prefixes
var (
	// ParamsKey is the prefix for params key
	ParamsPrefix = collections.NewPrefix(0)
)
