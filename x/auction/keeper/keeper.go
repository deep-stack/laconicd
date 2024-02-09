package keeper

import (
	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	auctiontypes "git.vdb.to/cerc-io/laconic2d/x/auction"
)

type Keeper struct {
	// Codecs
	cdc codec.BinaryCodec

	// External keepers
	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper

	// Track auction usage in other cosmos-sdk modules (more like a usage tracker).
	// usageKeepers []types.AuctionUsageKeeper

	// state management
	Schema collections.Schema
	Params collections.Item[auctiontypes.Params]
	// TODO
	// Auctions ...
}

// NewKeeper creates a new Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	accountKeeper auth.AccountKeeper,
	bankKeeper bank.Keeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:           cdc,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		Params:        collections.NewItem(sb, auctiontypes.ParamsKeyPrefix, "params", codec.CollValue[auctiontypes.Params](cdc)),
		// Auctions:	 ...
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

// func (k *Keeper) SetUsageKeepers(usageKeepers []types.AuctionUsageKeeper) {
// 	k.usageKeepers = usageKeepers
// }
