package module

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"git.vdb.to/cerc-io/laconicd/x/registry/keeper"
)

// EndBlocker is called every block
func EndBlocker(ctx context.Context, k keeper.Keeper) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := k.ProcessRecordExpiryQueue(sdkCtx); err != nil {
		return err
	}

	if err := k.ProcessAuthorityExpiryQueue(sdkCtx); err != nil {
		return err
	}

	return nil
}
