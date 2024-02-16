package registry

import (
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
