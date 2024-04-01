package keeper

import (
	"git.vdb.to/cerc-io/laconicd/x/bond"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module state from a genesis state.
func (k *Keeper) InitGenesis(ctx sdk.Context, data *bond.GenesisState) error {
	if err := k.Params.Set(ctx, data.Params); err != nil {
		return err
	}

	// Save bonds in store.
	for _, bond := range data.Bonds {
		if err := k.Bonds.Set(ctx, bond.Id, *bond); err != nil {
			return err
		}
	}

	return nil
}

// ExportGenesis exports the module state to a genesis state.
func (k *Keeper) ExportGenesis(ctx sdk.Context) (*bond.GenesisState, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	bonds, err := k.ListBonds(ctx)
	if err != nil {
		return nil, err
	}

	return &bond.GenesisState{Params: params, Bonds: bonds}, nil
}
