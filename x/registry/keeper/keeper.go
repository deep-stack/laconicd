package keeper

import (
	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	auctionkeeper "git.vdb.to/cerc-io/laconic2d/x/auction/keeper"
	bondkeeper "git.vdb.to/cerc-io/laconic2d/x/bond/keeper"
	registrytypes "git.vdb.to/cerc-io/laconic2d/x/registry"
)

// TODO: Add required methods

type Keeper struct {
	cdc codec.BinaryCodec

	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper
	recordKeeper  RecordKeeper
	bondKeeper    bondkeeper.Keeper
	auctionKeeper auctionkeeper.Keeper

	// state management
	Schema collections.Schema
	Params collections.Item[registrytypes.Params]
}

// NewKeeper creates a new Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	accountKeeper auth.AccountKeeper,
	bankKeeper bank.Keeper,
	recordKeeper RecordKeeper,
	bondKeeper bondkeeper.Keeper,
	auctionKeeper auctionkeeper.Keeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:    cdc,
		Params: collections.NewItem(sb, registrytypes.ParamsPrefix, "params", codec.CollValue[registrytypes.Params](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}
