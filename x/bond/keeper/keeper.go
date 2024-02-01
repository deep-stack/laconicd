package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"git.vdb.to/cerc-io/laconic2d/x/bond"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: Add genesis.go?

type Keeper struct {
	// Store keys
	storeKey storetypes.StoreKey

	// Codecs
	cdc codec.BinaryCodec

	// External keepers
	// accountKeeper auth.AccountKeeper
	// bankKeeper    bank.Keeper

	// Track bond usage in other cosmos-sdk modules (more like a usage tracker).
	// usageKeepers []types.BondUsageKeeper

	// paramSubspace paramtypes.Subspace
}

// NewKeeper creates new instances of the bond Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	// accountKeeper auth.AccountKeeper,
	// bankKeeper bank.Keeper,
	// usageKeepers []types.BondUsageKeeper,
	storeKey storetypes.StoreKey,
	// ps paramtypes.Subspace,
) Keeper {
	// set KeyTable if it has not already been set
	// if !ps.HasKeyTable() {
	// 	ps = ps.WithKeyTable(types.ParamKeyTable())
	// }

	return Keeper{
		// accountKeeper: accountKeeper,
		// bankKeeper:    bankKeeper,
		storeKey: storeKey,
		cdc:      cdc,
		// usageKeepers:  usageKeepers,
		// paramSubspace: ps,
	}
}

// TODO: Add keeper methods

// SaveBond - saves a bond to the store.
func (k Keeper) SaveBond(ctx sdk.Context, bond *bond.Bond) {
	// TODO: Implement
}

// ListBonds - get all bonds.
func (k Keeper) ListBonds(ctx sdk.Context) []*bond.Bond {
	// TODO: Implement
	var bonds []*bond.Bond
	return bonds
}
