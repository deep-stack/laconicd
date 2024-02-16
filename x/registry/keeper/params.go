package keeper

import (
	registrytypes "git.vdb.to/cerc-io/laconic2d/x/registry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams - Get all parameters as types.Params.
func (k Keeper) GetParams(ctx sdk.Context) (*registrytypes.Params, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &params, nil
}
