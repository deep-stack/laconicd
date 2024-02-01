package keeper

import (
	"git.vdb.to/cerc-io/laconic2d/x/bond"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetMaxBondAmount max bond amount
func (k Keeper) GetMaxBondAmount(ctx sdk.Context) (res sdk.Coin) {
	// TODO: Implement
	return sdk.Coin{}
}

// GetParams - Get all parameter as types.Params.
func (k Keeper) GetParams(ctx sdk.Context) (params bond.Params) {
	getMaxBondAmount := k.GetMaxBondAmount(ctx)
	return bond.Params{MaxBondAmount: getMaxBondAmount}
}

// SetParams - set the params.
func (k Keeper) SetParams(ctx sdk.Context, params bond.Params) {
	// TODO: Implement
}
