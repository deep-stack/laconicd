package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	onboardingtypes "git.vdb.to/cerc-io/laconicd/x/onboarding"
)

var _ onboardingtypes.QueryServer = queryServer{}

type queryServer struct {
	k *Keeper
}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k *Keeper) onboardingtypes.QueryServer {
	return queryServer{k}
}

// Participants implements Participants.QueryServer.
func (qs queryServer) Participants(c context.Context, _ *onboardingtypes.QueryParticipantsRequest) (*onboardingtypes.QueryParticipantsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	resp, err := qs.k.ListParticipants(ctx)
	if err != nil {
		return nil, err
	}

	return &onboardingtypes.QueryParticipantsResponse{Participants: resp}, nil
}
