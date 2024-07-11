package keeper

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	types "git.vdb.to/cerc-io/laconicd/x/registry"
)

// RegisterInvariants registers all registry invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "module-account", ModuleAccountInvariant(&k))
	ir.RegisterRoute(types.ModuleName, "record-bond", RecordBondInvariant(&k))
}

// AllInvariants runs all invariants of the registry module.
func AllInvariants(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		res, stop := ModuleAccountInvariant(k)(ctx)
		if stop {
			return res, stop
		}

		return RecordBondInvariant(k)(ctx)
	}
}

// ModuleAccountInvariant checks that the 'registry' module account balance is non-negative.
func ModuleAccountInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		moduleAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
		if k.bankKeeper.GetAllBalances(ctx, moduleAddress).IsAnyNegative() {
			return sdk.FormatInvariant(
				types.ModuleName,
				"module-account",
				fmt.Sprintf("Module account '%s' has negative balance.", types.ModuleName),
			), true
		}

		return "", false
	}
}

// RecordBondInvariant checks that for every record:
// if bondId is not null, associated bond exists
func RecordBondInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		err := k.Records.Walk(ctx, nil, func(key string, record types.Record) (bool, error) {
			if record.BondId != "" {
				bondExists, err := k.bondKeeper.HasBond(ctx, record.BondId)
				if err != nil {
					return true, err
				}

				if !bondExists {
					return true, errors.New(record.Id)
				}
			}

			return false, nil
		})
		if err != nil {
			return sdk.FormatInvariant(
				types.ModuleName,
				"record-bond",
				fmt.Sprintf("Bond not found for record id: '%s'.", err.Error()),
			), true
		}

		return "", false
	}
}
