package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"git.vdb.to/cerc-io/laconicd/utils"
	"git.vdb.to/cerc-io/laconicd/x/bond"
)

var _ bond.MsgServer = msgServer{}

type msgServer struct {
	k *Keeper
}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper *Keeper) bond.MsgServer {
	return &msgServer{k: keeper}
}

func (ms msgServer) CreateBond(c context.Context, msg *bond.MsgCreateBond) (*bond.MsgCreateBondResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	ctx = *utils.CtxWithCustomKVGasConfig(&ctx)

	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	resp, err := ms.k.CreateBond(ctx, signerAddress, msg.Coins)
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

	utils.LogTxGasConsumed(ctx, ms.k.Logger(ctx), "CreateBond")

	return &bond.MsgCreateBondResponse{Id: resp.Id}, nil
}

// RefillBond implements bond.MsgServer.
func (ms msgServer) RefillBond(c context.Context, msg *bond.MsgRefillBond) (*bond.MsgRefillBondResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	ctx = *utils.CtxWithCustomKVGasConfig(&ctx)

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

	utils.LogTxGasConsumed(ctx, ms.k.Logger(ctx), "RefillBond")

	return &bond.MsgRefillBondResponse{}, nil
}

// WithdrawBond implements bond.MsgServer.
func (ms msgServer) WithdrawBond(c context.Context, msg *bond.MsgWithdrawBond) (*bond.MsgWithdrawBondResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	ctx = *utils.CtxWithCustomKVGasConfig(&ctx)

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

	utils.LogTxGasConsumed(ctx, ms.k.Logger(ctx), "WithdrawBond")

	return &bond.MsgWithdrawBondResponse{}, nil
}

// CancelBond implements bond.MsgServer.
func (ms msgServer) CancelBond(c context.Context, msg *bond.MsgCancelBond) (*bond.MsgCancelBondResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	ctx = *utils.CtxWithCustomKVGasConfig(&ctx)

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

	utils.LogTxGasConsumed(ctx, ms.k.Logger(ctx), "CancelBond")

	return &bond.MsgCancelBondResponse{}, nil
}
