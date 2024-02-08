package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"git.vdb.to/cerc-io/laconic2d/x/bond"
)

type msgServer struct {
	k Keeper
}

var _ bond.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper Keeper) bond.MsgServer {
	return &msgServer{k: keeper}
}

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

// RefillBond implements bond.MsgServer.
func (ms msgServer) RefillBond(c context.Context, msg *bond.MsgRefillBond) (*bond.MsgRefillBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	_, err = ms.k.RefillBond(ctx, msg.Id, signerAddress, msg.Coins)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			bond.EventTypeRefillBond,
			sdk.NewAttribute(bond.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(bond.AttributeKeyBondId, msg.Id),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, bond.AttributeValueCategory),
			sdk.NewAttribute(bond.AttributeKeySigner, msg.Signer),
		),
	})

	return &bond.MsgRefillBondResponse{}, nil
}

// WithdrawBond implements bond.MsgServer.
func (ms msgServer) WithdrawBond(c context.Context, msg *bond.MsgWithdrawBond) (*bond.MsgWithdrawBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	_, err = ms.k.WithdrawBond(ctx, msg.Id, signerAddress, msg.Coins)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			bond.EventTypeWithdrawBond,
			sdk.NewAttribute(bond.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(bond.AttributeKeyBondId, msg.Id),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Coins.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, bond.AttributeValueCategory),
			sdk.NewAttribute(bond.AttributeKeySigner, msg.Signer),
		),
	})

	return &bond.MsgWithdrawBondResponse{}, nil
}

// CancelBond implements bond.MsgServer.
func (ms msgServer) CancelBond(c context.Context, msg *bond.MsgCancelBond) (*bond.MsgCancelBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	_, err = ms.k.CancelBond(ctx, msg.Id, signerAddress)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			bond.EventTypeCancelBond,
			sdk.NewAttribute(bond.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(bond.AttributeKeyBondId, msg.Id),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, bond.AttributeValueCategory),
			sdk.NewAttribute(bond.AttributeKeySigner, msg.Signer),
		),
	})

	return &bond.MsgCancelBondResponse{}, nil
}
