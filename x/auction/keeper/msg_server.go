package keeper

import (
	"context"

	auctiontypes "git.vdb.to/cerc-io/laconic2d/x/auction"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ auctiontypes.MsgServer = msgServer{}

type msgServer struct {
	k *Keeper
}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper *Keeper) auctiontypes.MsgServer {
	return &msgServer{k: keeper}
}

func (ms msgServer) CreateAuction(c context.Context, msg *auctiontypes.MsgCreateAuction) (*auctiontypes.MsgCreateAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	resp, err := ms.k.CreateAuction(ctx, *msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			auctiontypes.EventTypeCreateAuction,
			sdk.NewAttribute(auctiontypes.AttributeKeyCommitsDuration, msg.CommitsDuration.String()),
			sdk.NewAttribute(auctiontypes.AttributeKeyCommitFee, msg.CommitFee.String()),
			sdk.NewAttribute(auctiontypes.AttributeKeyRevealFee, msg.RevealFee.String()),
			sdk.NewAttribute(auctiontypes.AttributeKeyMinimumBid, msg.MinimumBid.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, auctiontypes.AttributeValueCategory),
			sdk.NewAttribute(auctiontypes.AttributeKeySigner, signerAddress.String()),
		),
	})

	return &auctiontypes.MsgCreateAuctionResponse{Auction: resp}, nil
}

// CommitBid is the command for committing a bid
// nolint: all
func (ms msgServer) CommitBid(c context.Context, msg *auctiontypes.MsgCommitBid) (*auctiontypes.MsgCommitBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	resp, err := ms.k.CommitBid(ctx, *msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			auctiontypes.EventTypeCommitBid,
			sdk.NewAttribute(auctiontypes.AttributeKeyAuctionId, msg.AuctionId),
			sdk.NewAttribute(auctiontypes.AttributeKeyCommitHash, msg.CommitHash),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, auctiontypes.AttributeValueCategory),
			sdk.NewAttribute(auctiontypes.AttributeKeySigner, signerAddress.String()),
		),
	})

	return &auctiontypes.MsgCommitBidResponse{Bid: resp}, nil
}

// RevealBid is the command for revealing a bid
// nolint: all
func (ms msgServer) RevealBid(c context.Context, msg *auctiontypes.MsgRevealBid) (*auctiontypes.MsgRevealBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	resp, err := ms.k.RevealBid(ctx, *msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			auctiontypes.EventTypeRevealBid,
			sdk.NewAttribute(auctiontypes.AttributeKeyAuctionId, msg.AuctionId),
			sdk.NewAttribute(auctiontypes.AttributeKeyReveal, msg.Reveal),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, auctiontypes.AttributeValueCategory),
			sdk.NewAttribute(auctiontypes.AttributeKeySigner, signerAddress.String()),
		),
	})

	return &auctiontypes.MsgRevealBidResponse{Auction: resp}, nil
}
