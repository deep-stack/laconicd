package utils

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CtxWithCustomKVGasConfig(ctx *sdk.Context) *sdk.Context {
	updatedCtx := ctx.WithKVGasConfig(storetypes.GasConfig{
		HasCost:          0,
		DeleteCost:       0,
		ReadCostFlat:     0,
		ReadCostPerByte:  0,
		WriteCostFlat:    0,
		WriteCostPerByte: 0,
		IterNextCostFlat: 0,
	})

	return &updatedCtx
}

func LogTxGasConsumed(ctx sdk.Context, logger log.Logger, tx string) {
	gasConsumed := ctx.GasMeter().GasConsumed()
	logger.Info("tx executed", "method", tx, "gas_consumed", fmt.Sprintf("%d", gasConsumed))
}
