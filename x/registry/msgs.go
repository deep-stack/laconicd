package registry

import (
	"net/url"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewMsgSetRecord is the constructor function for MsgSetRecord.
func NewMsgSetRecord(payload Payload, bondId string, signer sdk.AccAddress) *MsgSetRecord {
	return &MsgSetRecord{
		Payload: payload,
		BondId:  bondId,
		Signer:  signer.String(),
	}
}

func (msg MsgSetRecord) ValidateBasic() error {
	if len(msg.Signer) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}

	owners := msg.Payload.Record.Owners
	for _, owner := range owners {
		if owner == "" {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Record owner not set.")
		}
	}

	if len(msg.BondId) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Bond Id is required.")
	}

	return nil
}

func (msg MsgRenewRecord) ValidateBasic() error {
	if len(msg.RecordId) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "record id is required.")
	}

	if len(msg.Signer) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer.")
	}

	return nil
}

// NewMsgReserveAuthority is the constructor function for MsgReserveName.
func NewMsgReserveAuthority(name string, signer sdk.AccAddress, owner sdk.AccAddress) MsgReserveAuthority {
	return MsgReserveAuthority{
		Name:   name,
		Owner:  owner.String(),
		Signer: signer.String(),
	}
}

// ValidateBasic Implements Msg.
func (msg MsgReserveAuthority) ValidateBasic() error {
	if len(msg.Name) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "name is required.")
	}

	if len(msg.Signer) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer")
	}

	return nil
}

// NewMsgSetAuthorityBond is the constructor function for MsgSetAuthorityBond.
func NewMsgSetAuthorityBond(name string, bondID string, signer sdk.AccAddress) MsgSetAuthorityBond {
	return MsgSetAuthorityBond{
		Name:   name,
		Signer: signer.String(),
		BondId: bondID,
	}
}

func (msg MsgSetAuthorityBond) ValidateBasic() error {
	if len(msg.Name) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "name is required.")
	}

	if len(msg.Signer) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer.")
	}

	if len(msg.BondId) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "bond id is required.")
	}

	return nil
}

// NewMsgSetName is the constructor function for MsgSetName.
func NewMsgSetName(lrn string, cid string, signer sdk.AccAddress) *MsgSetName {
	return &MsgSetName{
		Lrn:    lrn,
		Cid:    cid,
		Signer: signer.String(),
	}
}

// ValidateBasic Implements Msg.
func (msg MsgSetName) ValidateBasic() error {
	if msg.Lrn == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "LRN is required.")
	}

	if msg.Cid == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "CID is required.")
	}

	if len(msg.Signer) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer")
	}

	return nil
}

func (msg MsgDeleteName) ValidateBasic() error {
	if len(msg.Lrn) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "lrn is required.")
	}

	if len(msg.Signer) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer.")
	}

	_, err := url.Parse(msg.Lrn)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid lrn.")
	}

	return nil
}

func (msg MsgAssociateBond) ValidateBasic() error {
	if len(msg.RecordId) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "record id is required.")
	}
	if len(msg.BondId) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "bond id is required.")
	}
	if len(msg.Signer) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer.")
	}

	return nil
}

func (msg MsgDissociateBond) ValidateBasic() error {
	if len(msg.RecordId) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "record id is required.")
	}
	if len(msg.Signer) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer.")
	}

	return nil
}

func (msg MsgDissociateRecords) ValidateBasic() error {
	if len(msg.BondId) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "bond id is required.")
	}
	if len(msg.Signer) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer.")
	}

	return nil
}

func (msg MsgReassociateRecords) ValidateBasic() error {
	if len(msg.OldBondId) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "old-bond-id is required.")
	}
	if len(msg.NewBondId) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "new-bond-id is required.")
	}
	if len(msg.Signer) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer.")
	}

	return nil
}
