package keeper

import (
	"context"

	"git.vdb.to/cerc-io/laconic2d/x/bond"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: Add remaining query methods

type queryServer struct {
	k Keeper
}

var _ bond.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k Keeper) bond.QueryServer {
	return queryServer{k}
}

func (qs queryServer) Bonds(c context.Context, _ *bond.QueryGetBondsRequest) (*bond.QueryGetBondsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	resp, err := qs.k.ListBonds(ctx)
	if err != nil {
		return nil, err
	}

	return &bond.QueryGetBondsResponse{Bonds: resp}, nil
}
