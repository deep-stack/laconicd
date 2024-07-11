package onboarding

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgOnboardParticipant{}

func (msg MsgOnboardParticipant) ValidateBasic() error {
	if len(msg.Participant) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrorInvalidSigner, msg.Participant)
	}

	if len(msg.EthPayload.Address) != 42 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, msg.Participant)
	}

	if len(msg.EthSignature) != 132 {
		return errorsmod.Wrap(sdkerrors.ErrNoSignatures, "Invalid signature.")
	}

	return nil
}
