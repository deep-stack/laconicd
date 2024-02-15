package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"

	auctionkeeper "git.vdb.to/cerc-io/laconic2d/x/auction/keeper"
)

// TODO: Add methods

// RecordKeeper exposes the bare minimal read-only API for other modules.
type RecordKeeper struct {
	cdc           codec.BinaryCodec // The wire codec for binary encoding/decoding.
	auctionKeeper auctionkeeper.Keeper
	// storeKey      storetypes.StoreKey // Unexposed key to access store from sdk.Context
}
