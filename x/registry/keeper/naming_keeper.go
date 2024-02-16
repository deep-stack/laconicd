package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	registrytypes "git.vdb.to/cerc-io/laconic2d/x/registry"
)

// GetNameAuthority - gets a name authority from the store.
func (k Keeper) GetNameAuthority(ctx sdk.Context, name string) registrytypes.NameAuthority {
	panic("unimplemented")
}

// HasNameRecord - checks if a name record exists.
func (k Keeper) HasNameRecord(ctx sdk.Context, crn string) bool {
	panic("unimplemented")
}

// GetNameRecord - gets a name record from the store.
func (k Keeper) GetNameRecord(ctx sdk.Context, crn string) *registrytypes.NameRecord {
	panic("unimplemented")
}

// ListNameRecords - get all name records.
func (k Keeper) ListNameRecords(ctx sdk.Context) []registrytypes.NameEntry {
	panic("unimplemented")
}

// ProcessSetName creates a CRN -> Record ID mapping.
func (k Keeper) ProcessSetName(ctx sdk.Context, msg registrytypes.MsgSetName) error {
	panic("unimplemented")
}

// ProcessReserveAuthority reserves a name authority.
func (k Keeper) ProcessReserveAuthority(ctx sdk.Context, msg registrytypes.MsgReserveAuthority) error {
	panic("unimplemented")
}

func (k Keeper) ProcessSetAuthorityBond(ctx sdk.Context, msg registrytypes.MsgSetAuthorityBond) error {
	panic("unimplemented")
}

// ProcessDeleteName removes a CRN -> Record ID mapping.
func (k Keeper) ProcessDeleteName(ctx sdk.Context, msg registrytypes.MsgDeleteNameAuthority) error {
	panic("unimplemented")
}

func (k Keeper) GetAuthorityExpiryQueue(ctx sdk.Context) []*registrytypes.ExpiryQueueRecord {
	panic("unimplemented")
}

// ResolveCRN resolves a CRN to a record.
func (k Keeper) ResolveCRN(ctx sdk.Context, crn string) *registrytypes.Record {
	panic("unimplemented")
}
