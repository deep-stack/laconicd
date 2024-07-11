package bond

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateBond{}

// NewMsgCreateBond is the constructor function for MsgCreateBond.
func NewMsgCreateBond(coins sdk.Coins, signer sdk.AccAddress) MsgCreateBond {
	return MsgCreateBond{
		Coins:  coins,
		Signer: signer.String(),
	}
}

func (msg MsgCreateBond) ValidateBasic() error {
	if len(msg.Signer) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}
	if len(msg.Coins) == 0 || !msg.Coins.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "Invalid amount.")
	}
	return nil
}

func (msg MsgRefillBond) ValidateBasic() error {
	if len(msg.Id) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, msg.Id)
	}
	if len(msg.Signer) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}
	if len(msg.Coins) == 0 || !msg.Coins.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "Invalid amount.")
	}
	return nil
}

func (msg MsgWithdrawBond) ValidateBasic() error {
	if len(msg.Id) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, msg.Id)
	}
	if len(msg.Signer) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}
	if len(msg.Coins) == 0 || !msg.Coins.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "Invalid amount.")
	}
	return nil
}

func (msg MsgCancelBond) ValidateBasic() error {
	if len(msg.Id) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, msg.Id)
	}
	if len(msg.Signer) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}
	return nil
}
