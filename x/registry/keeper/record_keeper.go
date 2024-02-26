package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	auctionkeeper "git.vdb.to/cerc-io/laconic2d/x/auction/keeper"
	registrytypes "git.vdb.to/cerc-io/laconic2d/x/registry"
)

// TODO: Add methods

// RecordKeeper exposes the bare minimal read-only API for other modules.
type RecordKeeper struct {
	cdc           codec.BinaryCodec // The wire codec for binary encoding/decoding.
	auctionKeeper auctionkeeper.Keeper
	// storeKey      storetypes.StoreKey // Unexposed key to access store from sdk.Context
}

// ProcessRenewRecord renews a record.
func (k Keeper) ProcessRenewRecord(ctx sdk.Context, msg registrytypes.MsgRenewRecord) error {
	panic("unimplemented")
}

// ProcessAssociateBond associates a record with a bond.
func (k Keeper) ProcessAssociateBond(ctx sdk.Context, msg registrytypes.MsgAssociateBond) error {
	panic("unimplemented")
}

// ProcessDissociateBond dissociates a record from its bond.
func (k Keeper) ProcessDissociateBond(ctx sdk.Context, msg registrytypes.MsgDissociateBond) error {
	panic("unimplemented")
}

// ProcessDissociateRecords dissociates all records associated with a given bond.
func (k Keeper) ProcessDissociateRecords(ctx sdk.Context, msg registrytypes.MsgDissociateRecords) error {
	panic("unimplemented")
}

// ProcessReAssociateRecords switches records from and old to new bond.
func (k Keeper) ProcessReAssociateRecords(ctx sdk.Context, msg registrytypes.MsgReAssociateRecords) error {
	panic("unimplemented")
}
