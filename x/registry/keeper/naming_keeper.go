package keeper

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	auctiontypes "git.vdb.to/cerc-io/laconic2d/x/auction"
	registrytypes "git.vdb.to/cerc-io/laconic2d/x/registry"
	"git.vdb.to/cerc-io/laconic2d/x/registry/helpers"
)

// HasNameAuthority - checks if a name/authority exists.
func (k Keeper) HasNameAuthority(ctx sdk.Context, name string) (bool, error) {
	has, err := k.Authorities.Has(ctx, name)
	if err != nil {
		return false, err
	}

	return has, nil
}

// GetNameAuthority - gets a name authority from the store.
func (k Keeper) GetNameAuthority(ctx sdk.Context, name string) (registrytypes.NameAuthority, error) {
	authority, err := k.Authorities.Get(ctx, name)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return registrytypes.NameAuthority{}, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Name authority not found.")
		}
		return registrytypes.NameAuthority{}, err
	}

	return authority, nil
}

// ListNameAuthorityRecords - get all name authority records.
func (k Keeper) ListNameAuthorityRecords(ctx sdk.Context) (map[string]registrytypes.NameAuthority, error) {
	nameAuthorityRecords := make(map[string]registrytypes.NameAuthority)

	err := k.Authorities.Walk(ctx, nil, func(key string, value registrytypes.NameAuthority) (bool, error) {
		nameAuthorityRecords[key] = value
		return false, nil
	})
	if err != nil {
		return map[string]registrytypes.NameAuthority{}, err
	}

	return nameAuthorityRecords, nil
}

// HasNameRecord - checks if a name record exists.
func (k Keeper) HasNameRecord(ctx sdk.Context, crn string) (bool, error) {
	return k.NameRecords.Has(ctx, crn)
}

// GetNameRecord - gets a name record from the store.
func (k Keeper) GetNameRecord(ctx sdk.Context, crn string) (*registrytypes.NameRecord, error) {
	nameRecord, err := k.NameRecords.Get(ctx, crn)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &nameRecord, nil
}

// LookupNameRecord - gets a name record which is not stale and under active authority.
func (k Keeper) LookupNameRecord(ctx sdk.Context, crn string) (*registrytypes.NameRecord, error) {
	_, _, authority, err := k.getAuthority(ctx, crn)
	if err != nil || authority.Status != registrytypes.AuthorityActive {
		// If authority is not active (or any other error), lookup fails.
		return nil, nil
	}

	nameRecord, err := k.GetNameRecord(ctx, crn)

	// Name record may not exist.
	if nameRecord == nil {
		return nil, err
	}

	// Name lookup should fail if the name record is stale.
	// i.e. authority was registered later than the name.
	if authority.Height > nameRecord.Latest.Height {
		return nil, nil
	}

	return nameRecord, nil
}

// ListNameRecords - get all name records.
func (k Keeper) ListNameRecords(ctx sdk.Context) ([]registrytypes.NameEntry, error) {
	var nameEntries []registrytypes.NameEntry

	err := k.NameRecords.Walk(ctx, nil, func(key string, value registrytypes.NameRecord) (stop bool, err error) {
		nameEntries = append(nameEntries, registrytypes.NameEntry{
			Name:  key,
			Entry: &value,
		})

		return false, nil
	})

	return nameEntries, err
}

// SaveNameRecord - sets a name record.
func (k Keeper) SaveNameRecord(ctx sdk.Context, crn string, id string) error {
	var nameRecord registrytypes.NameRecord
	existingNameRecord, err := k.GetNameRecord(ctx, crn)
	if err != nil {
		return err
	}

	if existingNameRecord != nil {
		nameRecord = *existingNameRecord
		nameRecord.History = append(nameRecord.History, nameRecord.Latest)
	}

	nameRecord.Latest = &registrytypes.NameRecordEntry{
		Id:     id,
		Height: uint64(ctx.BlockHeight()),
	}

	// TODO: Check if index gets updated on entry updates
	if err := k.NameRecords.Set(ctx, crn, nameRecord); err != nil {
		return err
	}

	// TODO
	// Update changeSet for name.
	// k.updateBlockChangeSetForName(ctx, crn)

	return nil
}

// SetName creates a CRN -> Record ID mapping.
func (k Keeper) SetName(ctx sdk.Context, msg registrytypes.MsgSetName) error {
	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return err
	}
	err = k.checkCRNAccess(ctx, signerAddress, msg.Crn)
	if err != nil {
		return err
	}

	nameRecord, err := k.LookupNameRecord(ctx, msg.Crn)
	if err != nil {
		return err
	}
	if nameRecord != nil && nameRecord.Latest.Id == msg.Cid {
		return nil
	}

	return k.SaveNameRecord(ctx, msg.Crn, msg.Cid)
}

// SaveNameAuthority creates the NameAuthority record.
func (k Keeper) SaveNameAuthority(ctx sdk.Context, name string, authority *registrytypes.NameAuthority) error {
	return k.Authorities.Set(ctx, name, *authority)

	// TODO
	// updateBlockChangeSetForNameAuthority(ctx, codec, store, name)
}

// ReserveAuthority reserves a name authority.
func (k Keeper) ReserveAuthority(ctx sdk.Context, msg registrytypes.MsgReserveAuthority) error {
	crn := fmt.Sprintf("crn://%s", msg.GetName())
	parsedCrn, err := url.Parse(crn)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Invalid name")
	}

	name := parsedCrn.Host
	if fmt.Sprintf("crn://%s", name) != crn {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Invalid name")
	}

	if strings.Contains(name, ".") {
		return k.ReserveSubAuthority(ctx, name, msg)
	}

	err = k.createAuthority(ctx, name, msg.GetSigner(), true)
	if err != nil {
		return err
	}

	return nil
}

// ReserveSubAuthority reserves a sub-authority.
func (k Keeper) ReserveSubAuthority(ctx sdk.Context, name string, msg registrytypes.MsgReserveAuthority) error {
	// Get parent authority name.
	names := strings.Split(name, ".")
	parent := strings.Join(names[1:], ".")

	// Check if parent authority exists.
	if has, err := k.HasNameAuthority(ctx, parent); !has {
		if err != nil {
			return err
		}

		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Parent authority not found.")
	}
	parentAuthority, err := k.GetNameAuthority(ctx, parent)
	if err != nil {
		return err
	}

	// Sub-authority creator needs to be the owner of the parent authority.
	if parentAuthority.OwnerAddress != msg.Signer {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Access denied.")
	}

	// Sub-authority owner defaults to parent authority owner.
	subAuthorityOwner := msg.Signer
	if len(msg.Owner) != 0 {
		// Override sub-authority owner if provided in message.
		subAuthorityOwner = msg.Owner
	}

	sdkErr := k.createAuthority(ctx, name, subAuthorityOwner, false)
	if sdkErr != nil {
		return sdkErr
	}

	return nil
}

func (k Keeper) createAuthority(ctx sdk.Context, name string, owner string, isRoot bool) error {
	moduleParams, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	has, err := k.HasNameAuthority(ctx, name)
	if err != nil {
		return err
	}
	if has {
		authority, err := k.GetNameAuthority(ctx, name)
		if err != nil {
			return err
		}

		if authority.Status != registrytypes.AuthorityExpired {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Name already reserved.")
		}
	}

	ownerAddress, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Invalid owner address.")
	}
	ownerAccount := k.accountKeeper.GetAccount(ctx, ownerAddress)
	if ownerAccount == nil {
		return errorsmod.Wrap(sdkerrors.ErrUnknownAddress, "Owner account not found.")
	}

	authority := registrytypes.NameAuthority{
		OwnerPublicKey: getAuthorityPubKey(ownerAccount.GetPubKey()),
		OwnerAddress:   owner,
		Height:         uint64(ctx.BlockHeight()),
		Status:         registrytypes.AuthorityActive,
		AuctionId:      "",
		BondId:         "",
		ExpiryTime:     ctx.BlockTime().Add(moduleParams.AuthorityGracePeriod),
	}

	if isRoot && moduleParams.AuthorityAuctionEnabled {
		// If auctions are enabled, clear out owner fields. They will be set after a winner is picked.
		authority.OwnerAddress = ""
		authority.OwnerPublicKey = ""
		// Reset bond ID if required.
		if authority.BondId != "" || len(authority.BondId) != 0 {
			// TODO
			// k.RemoveBondToAuthorityIndexEntry(ctx, authority.BondId, name)
			authority.BondId = ""
		}

		params := auctiontypes.Params{
			CommitsDuration: moduleParams.AuthorityAuctionCommitsDuration,
			RevealsDuration: moduleParams.AuthorityAuctionRevealsDuration,
			CommitFee:       moduleParams.AuthorityAuctionCommitFee,
			RevealFee:       moduleParams.AuthorityAuctionRevealFee,
			MinimumBid:      moduleParams.AuthorityAuctionMinimumBid,
		}

		// Create an auction.
		msg := auctiontypes.NewMsgCreateAuction(params, ownerAddress)

		auction, sdkErr := k.auctionKeeper.CreateAuction(ctx, msg)
		if sdkErr != nil {
			return sdkErr
		}

		// TODO
		// Create auction ID -> authority name index.
		// k.AddAuctionToAuthorityMapping(ctx, auction.Id, name)

		authority.Status = registrytypes.AuthorityUnderAuction
		authority.AuctionId = auction.Id
		authority.ExpiryTime = auction.RevealsEndTime.Add(moduleParams.AuthorityGracePeriod)
	}

	// Save name authority in store.
	if err = k.SaveNameAuthority(ctx, name, &authority); err != nil {
		return err
	}

	// TODO
	// k.InsertAuthorityExpiryQueue(ctx, name, authority.ExpiryTime)

	return nil
}

func (k Keeper) SetAuthorityBond(ctx sdk.Context, msg registrytypes.MsgSetAuthorityBond) error {
	name := msg.GetName()
	signer := msg.GetSigner()

	if has, err := k.HasNameAuthority(ctx, name); !has {
		if err != nil {
			return err
		}

		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Name authority not found.")
	}

	authority, err := k.GetNameAuthority(ctx, name)
	if err != nil {
		return err
	}
	if authority.OwnerAddress != signer {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Access denied.")
	}

	if has, err := k.bondKeeper.HasBond(ctx, msg.BondId); !has {
		if err != nil {
			return err
		}
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Bond not found.")
	}

	bond, err := k.bondKeeper.GetBondById(ctx, msg.BondId)
	if err != nil {
		return err
	}
	if bond.Owner != signer {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Bond owner mismatch.")
	}

	// No-op if bond hasn't changed.
	if authority.BondId == msg.BondId {
		return nil
	}

	// TODO
	// Remove old bond ID mapping, if any.
	// if authority.BondId != "" {
	// 	k.RemoveBondToAuthorityIndexEntry(ctx, authority.BondId, name)
	// }

	// Update bond id and save name authority in store.
	authority.BondId = bond.Id
	if err = k.SaveNameAuthority(ctx, name, &authority); err != nil {
		return err
	}

	// TODO
	// Add new bond ID mapping.
	// k.AddBondToAuthorityIndexEntry(ctx, authority.BondId, name)

	return nil
}

// DeleteName removes a CRN -> Record ID mapping.
func (k Keeper) DeleteName(ctx sdk.Context, msg registrytypes.MsgDeleteNameAuthority) error {
	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return err
	}
	err = k.checkCRNAccess(ctx, signerAddress, msg.Crn)
	if err != nil {
		return err
	}

	crnExists, err := k.HasNameRecord(ctx, msg.Crn)
	if err != nil {
		return err
	}
	if !crnExists {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Name not found.")
	}

	// Set CID to empty string.
	return k.SaveNameRecord(ctx, msg.Crn, "")
}

func (k Keeper) GetAuthorityExpiryQueue(ctx sdk.Context) []*registrytypes.ExpiryQueueRecord {
	panic("unimplemented")
}

// ResolveCRN resolves a CRN to a record.
func (k Keeper) ResolveCRN(ctx sdk.Context, crn string) (*registrytypes.Record, error) {
	_, _, authority, err := k.getAuthority(ctx, crn)
	if err != nil || authority.Status != registrytypes.AuthorityActive {
		// If authority is not active (or any other error), resolution fails.
		return nil, err
	}

	// Name should not resolve if it's stale.
	// i.e. authority was registered later than the name.
	record, nameRecord, err := k.resolveCRNRecord(ctx, crn)
	if err != nil {
		return nil, err
	}
	if authority.Height > nameRecord.Latest.Height {
		return nil, nil
	}

	return record, nil
}

func (k Keeper) resolveCRNRecord(ctx sdk.Context, crn string) (*registrytypes.Record, *registrytypes.NameRecord, error) {
	nameRecord, err := k.GetNameRecord(ctx, crn)
	if nameRecord == nil {
		return nil, nil, err
	}

	latestRecordId := nameRecord.Latest.Id
	if latestRecordId == "" {
		return nil, nameRecord, nil
	}

	if has, err := k.HasRecord(ctx, latestRecordId); !has {
		if err != nil {
			return nil, nil, err
		}

		return nil, nameRecord, nil
	}

	record, err := k.GetRecordById(ctx, latestRecordId)
	if err != nil {
		return nil, nil, err
	}

	return &record, nameRecord, nil
}

func (k Keeper) getAuthority(ctx sdk.Context, crn string) (string, *url.URL, *registrytypes.NameAuthority, error) {
	parsedCRN, err := url.Parse(crn)
	if err != nil {
		return "", nil, nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Invalid CRN.")
	}

	name := parsedCRN.Host
	if has, err := k.HasNameAuthority(ctx, name); !has {
		if err != nil {
			return "", nil, nil, err
		}

		return "", nil, nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Name authority not found.")
	}

	authority, err := k.GetNameAuthority(ctx, name)
	if err != nil {
		return "", nil, nil, err
	}

	return name, parsedCRN, &authority, nil
}

func (k Keeper) checkCRNAccess(ctx sdk.Context, signer sdk.AccAddress, crn string) error {
	name, parsedCRN, authority, err := k.getAuthority(ctx, crn)
	if err != nil {
		return err
	}

	formattedCRN := fmt.Sprintf("crn://%s%s", name, parsedCRN.RequestURI())
	if formattedCRN != crn {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Invalid CRN.")
	}

	if authority.OwnerAddress != signer.String() {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Access denied.")
	}

	if authority.Status != registrytypes.AuthorityActive {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Authority is not active.")
	}

	if authority.BondId == "" || len(authority.BondId) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Authority bond not found.")
	}

	if authority.OwnerPublicKey == "" {
		// Try to set owner public key if account has it available now.
		ownerAccount := k.accountKeeper.GetAccount(ctx, signer)
		pubKey := ownerAccount.GetPubKey()
		if pubKey != nil {
			// Update public key in authority record.
			authority.OwnerPublicKey = getAuthorityPubKey(pubKey)
			if err = k.SaveNameAuthority(ctx, name, authority); err != nil {
				return err
			}
		}
	}

	return nil
}

func getAuthorityPubKey(pubKey cryptotypes.PubKey) string {
	if pubKey == nil {
		return ""
	}

	return helpers.BytesToBase64(pubKey.Bytes())
}
