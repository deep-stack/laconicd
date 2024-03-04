package bond

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgCreateBond{}
)

// NewMsgCreateBond is the constructor function for MsgCreateBond.
func NewMsgCreateBond(coins sdk.Coins, signer sdk.AccAddress) MsgCreateBond {
	return MsgCreateBond{
		Coins:  coins,
		Signer: signer.String(),
	}
}
