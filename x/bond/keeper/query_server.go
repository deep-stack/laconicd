package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	bondtypes "git.vdb.to/cerc-io/laconic2d/x/bond"
)

var _ bondtypes.QueryServer = queryServer{}

type queryServer struct {
	k *Keeper
}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k *Keeper) bondtypes.QueryServer {
	return queryServer{k}
}

// Params implements bond.QueryServer.
func (qs queryServer) Params(c context.Context, _ *bondtypes.QueryParamsRequest) (*bondtypes.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	params, err := qs.k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	return &bondtypes.QueryParamsResponse{Params: params}, nil
}

// Bonds implements bond.QueryServer.
func (qs queryServer) Bonds(c context.Context, _ *bondtypes.QueryBondsRequest) (*bondtypes.QueryBondsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	resp, err := qs.k.ListBonds(ctx)
	if err != nil {
		return nil, err
	}

	return &bondtypes.QueryBondsResponse{Bonds: resp}, nil
}

// GetBondById implements bond.QueryServer.
func (qs queryServer) GetBondById(c context.Context, req *bondtypes.QueryGetBondByIdRequest) (*bondtypes.QueryGetBondByIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	bondId := req.GetId()
	if len(bondId) == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "bond id required")
	}

	bond, err := qs.k.GetBondById(ctx, bondId)
	if err != nil {
		return nil, err
	}

	return &bondtypes.QueryGetBondByIdResponse{Bond: &bond}, nil
}

// GetBondsByOwner implements bond.QueryServer.
func (qs queryServer) GetBondsByOwner(c context.Context, req *bondtypes.QueryGetBondsByOwnerRequest) (*bondtypes.QueryGetBondsByOwnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	owner := req.GetOwner()
	if len(owner) == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "owner required")
	}

	bonds, err := qs.k.GetBondsByOwner(ctx, owner)
	if err != nil {
		return nil, err
	}

	return &bondtypes.QueryGetBondsByOwnerResponse{Bonds: bonds}, nil
}

// GetBondModuleBalance implements bond.QueryServer.
func (qs queryServer) GetBondModuleBalance(c context.Context, _ *bondtypes.QueryGetBondModuleBalanceRequest) (*bondtypes.QueryGetBondModuleBalanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	balances := qs.k.GetBondModuleBalances(ctx)

	return &bondtypes.QueryGetBondModuleBalanceResponse{Balance: balances}, nil
}
