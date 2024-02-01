package keeper

import (
	"git.vdb.to/cerc-io/laconic2d/x/bond"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module state from a genesis state.
func (k *Keeper) InitGenesis(ctx sdk.Context, data *bond.GenesisState) error {
	k.SetParams(ctx, data.Params)

	for _, bond := range data.Bonds {
		k.SaveBond(ctx, bond)
	}

	return nil
}

// ExportGenesis exports the module state to a genesis state.
func (k *Keeper) ExportGenesis(ctx sdk.Context) (*bond.GenesisState, error) {
	params := k.GetParams(ctx)
	bonds := k.ListBonds(ctx)

	return &bond.GenesisState{Params: params, Bonds: bonds}, nil
}
