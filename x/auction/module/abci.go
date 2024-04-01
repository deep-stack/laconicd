package module

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"git.vdb.to/cerc-io/laconicd/x/auction/keeper"
)

// EndBlocker is called every block
func EndBlocker(ctx context.Context, k *keeper.Keeper) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	return k.EndBlockerProcessAuctions(sdkCtx)
}
