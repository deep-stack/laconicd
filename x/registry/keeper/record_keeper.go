package keeper

import (
	"time"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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

// RenewRecord renews a record.
func (k Keeper) RenewRecord(ctx sdk.Context, msg registrytypes.MsgRenewRecord) error {
	if has, err := k.HasRecord(ctx, msg.RecordId); !has {
		if err != nil {
			return err
		}
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Record not found.")
	}

	// Check if renewal is required (i.e. expired record marked as deleted).
	record, err := k.GetRecordById(ctx, msg.RecordId)
	if err != nil {
		return err
	}

	expiryTime, err := time.Parse(time.RFC3339, record.ExpiryTime)
	if err != nil {
		return err
	}

	if !record.Deleted || expiryTime.After(ctx.BlockTime()) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Renewal not required.")
	}

	readableRecord := record.ToReadableRecord()
	return k.processRecord(ctx, &readableRecord, true)
}

// AssociateBond associates a record with a bond.
func (k Keeper) AssociateBond(ctx sdk.Context, msg registrytypes.MsgAssociateBond) error {
	if has, err := k.HasRecord(ctx, msg.RecordId); !has {
		if err != nil {
			return err
		}
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Record not found.")
	}

	if has, err := k.bondKeeper.HasBond(ctx, msg.BondId); !has {
		if err != nil {
			return err
		}
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Bond not found.")
	}

	// Check if already associated with a bond.
	record, err := k.GetRecordById(ctx, msg.RecordId)
	if err != nil {
		return err
	}

	if len(record.BondId) != 0 {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Bond already exists.")
	}

	// Only the bond owner can associate a record with the bond.
	bond, err := k.bondKeeper.GetBondById(ctx, msg.BondId)
	if err != nil {
		return err
	}
	if msg.Signer != bond.Owner {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Bond owner mismatch.")
	}

	record.BondId = msg.BondId
	if err = k.SaveRecord(ctx, record); err != nil {
		return err
	}

	// Required so that renewal is triggered (with new bond ID) for expired records.
	if record.Deleted {
		return k.insertRecordExpiryQueue(ctx, record)
	}

	return nil
}

// DissociateBond dissociates a record from its bond.
func (k Keeper) DissociateBond(ctx sdk.Context, msg registrytypes.MsgDissociateBond) error {
	if has, err := k.HasRecord(ctx, msg.RecordId); !has {
		if err != nil {
			return err
		}
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Record not found.")
	}

	record, err := k.GetRecordById(ctx, msg.RecordId)
	if err != nil {
		return err
	}

	// Check if record associated with a bond.
	bondId := record.BondId
	if len(bondId) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Bond not found.")
	}

	// Only the bond owner can dissociate a record with the bond.
	bond, err := k.bondKeeper.GetBondById(ctx, bondId)
	if err != nil {
		return err
	}
	if msg.Signer != bond.Owner {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Bond owner mismatch.")
	}

	// Clear bond Id.
	record.BondId = ""
	return k.SaveRecord(ctx, record)
}

// DissociateRecords dissociates all records associated with a given bond.
func (k Keeper) DissociateRecords(ctx sdk.Context, msg registrytypes.MsgDissociateRecords) error {
	if has, err := k.bondKeeper.HasBond(ctx, msg.BondId); !has {
		if err != nil {
			return err
		}
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Bond not found.")
	}

	// Only the bond owner can dissociate all records from the bond.
	bond, err := k.bondKeeper.GetBondById(ctx, msg.BondId)
	if err != nil {
		return err
	}
	if msg.Signer != bond.Owner {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Bond owner mismatch.")
	}

	// Dissociate all records from the bond.
	records, err := k.GetRecordsByBondId(ctx, msg.BondId)
	if err != nil {
		return err
	}

	for _, record := range records {
		// Clear bond Id.
		record.BondId = ""
		if err = k.SaveRecord(ctx, record); err != nil {
			return err
		}
	}

	return nil
}

// ReassociateRecords switches records from and old to new bond.
func (k Keeper) ReassociateRecords(ctx sdk.Context, msg registrytypes.MsgReassociateRecords) error {
	if has, err := k.bondKeeper.HasBond(ctx, msg.OldBondId); !has {
		if err != nil {
			return err
		}
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Old bond not found.")
	}

	if has, err := k.bondKeeper.HasBond(ctx, msg.NewBondId); !has {
		if err != nil {
			return err
		}
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "New bond not found.")
	}

	// Only the bond owner can re-associate all records.
	oldBond, err := k.bondKeeper.GetBondById(ctx, msg.OldBondId)
	if err != nil {
		return err
	}
	if msg.Signer != oldBond.Owner {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Old bond owner mismatch.")
	}

	newBond, err := k.bondKeeper.GetBondById(ctx, msg.NewBondId)
	if err != nil {
		return err
	}
	if msg.Signer != newBond.Owner {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "New bond owner mismatch.")
	}

	// Re-associate all records.
	records, err := k.GetRecordsByBondId(ctx, msg.OldBondId)
	if err != nil {
		return err
	}

	for _, record := range records {
		// Switch bond ID.
		record.BondId = msg.NewBondId
		if err = k.SaveRecord(ctx, record); err != nil {
			return err
		}

		// Required so that renewal is triggered (with new bond ID) for expired records.
		if record.Deleted {
			if err = k.insertRecordExpiryQueue(ctx, record); err != nil {
				return err
			}
		}
	}

	return nil
}
