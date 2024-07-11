package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	auctiontypes "git.vdb.to/cerc-io/laconicd/x/auction"
)

var _ auctiontypes.QueryServer = queryServer{}

type queryServer struct {
	k *Keeper
}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k *Keeper) auctiontypes.QueryServer {
	return queryServer{k}
}

// Params implements the params query command
func (qs queryServer) Params(c context.Context, req *auctiontypes.QueryParamsRequest) (*auctiontypes.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	params, err := qs.k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	return &auctiontypes.QueryParamsResponse{Params: params}, nil
}

// Auctions queries all auctions
func (qs queryServer) Auctions(c context.Context, req *auctiontypes.QueryAuctionsRequest) (*auctiontypes.QueryAuctionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	auctions, err := qs.k.ListAuctions(ctx)
	if err != nil {
		return nil, err
	}

	return &auctiontypes.QueryAuctionsResponse{Auctions: &auctiontypes.Auctions{Auctions: auctions}}, nil
}

// GetAuction queries an auction by id
func (qs queryServer) GetAuction(c context.Context, req *auctiontypes.QueryGetAuctionRequest) (*auctiontypes.QueryGetAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if req.Id == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "auction id is required")
	}

	auction, err := qs.k.GetAuctionById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &auctiontypes.QueryGetAuctionResponse{Auction: &auction}, nil
}

// GetBid queries an auction bid by auction-id and bidder
func (qs queryServer) GetBid(c context.Context, req *auctiontypes.QueryGetBidRequest) (*auctiontypes.QueryGetBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if req.AuctionId == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "auction id is required")
	}

	if req.Bidder == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "bidder address is required")
	}

	bid, err := qs.k.GetBid(ctx, req.AuctionId, req.Bidder)
	if err != nil {
		return nil, err
	}

	return &auctiontypes.QueryGetBidResponse{Bid: &bid}, nil
}

// GetBids queries all auction bids
func (qs queryServer) GetBids(c context.Context, req *auctiontypes.QueryGetBidsRequest) (*auctiontypes.QueryGetBidsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if req.AuctionId == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "auction id is required")
	}

	bids, err := qs.k.GetBids(ctx, req.AuctionId)
	if err != nil {
		return nil, err
	}

	return &auctiontypes.QueryGetBidsResponse{Bids: bids}, nil
}

// AuctionsByBidder queries auctions by bidder
func (qs queryServer) AuctionsByBidder(
	c context.Context,
	req *auctiontypes.QueryAuctionsByBidderRequest,
) (*auctiontypes.QueryAuctionsByBidderResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if req.BidderAddress == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "bidder address is required")
	}

	auctions, err := qs.k.QueryAuctionsByBidder(ctx, req.BidderAddress)
	if err != nil {
		return nil, err
	}

	return &auctiontypes.QueryAuctionsByBidderResponse{
		Auctions: &auctiontypes.Auctions{
			Auctions: auctions,
		},
	}, nil
}

// AuctionsByOwner queries auctions by owner
func (qs queryServer) AuctionsByOwner(
	c context.Context,
	req *auctiontypes.QueryAuctionsByOwnerRequest,
) (*auctiontypes.QueryAuctionsByOwnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if req.OwnerAddress == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "owner address is required")
	}

	auctions, err := qs.k.GetAuctionsByOwner(ctx, req.OwnerAddress)
	if err != nil {
		return nil, err
	}

	return &auctiontypes.QueryAuctionsByOwnerResponse{Auctions: &auctiontypes.Auctions{Auctions: auctions}}, nil
}

// GetAuctionModuleBalance queries the auction module account balance
func (qs queryServer) GetAuctionModuleBalance(
	c context.Context,
	req *auctiontypes.QueryGetAuctionModuleBalanceRequest,
) (*auctiontypes.QueryGetAuctionModuleBalanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	balances := qs.k.GetAuctionModuleBalances(ctx)

	return &auctiontypes.QueryGetAuctionModuleBalanceResponse{Balance: balances}, nil
}
