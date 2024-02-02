package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"git.vdb.to/cerc-io/laconic2d/x/bond"
)

type msgServer struct {
	k Keeper
}

// TODO: Generate types
var _ bond.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper Keeper) bond.MsgServer {
	return &msgServer{k: keeper}
}

// TODO: Add remaining write methods

func (ms msgServer) CreateBond(c context.Context, msg *bond.MsgCreateBond) (*bond.MsgCreateBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	_, err = ms.k.CreateBond(ctx, signerAddress, msg.Coins)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			bond.EventTypeCreateBond,
			sdk.NewAttribute(bond.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, bond.AttributeValueCategory),
			sdk.NewAttribute(bond.AttributeKeySigner, msg.Signer),
		),
	})

	return &bond.MsgCreateBondResponse{}, nil
}
