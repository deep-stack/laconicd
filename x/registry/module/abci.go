package module

import (
	"context"

	"git.vdb.to/cerc-io/laconic2d/x/registry/keeper"
)

// EndBlocker is called every block
func EndBlocker(ctx context.Context, k keeper.Keeper) error {
	// TODO: Implement

	// k.ProcessRecordExpiryQueue(ctx)
	// k.ProcessAuthorityExpiryQueue(ctx)

	return nil
}
