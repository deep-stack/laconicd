package module

import (
	"context"

	"git.vdb.to/cerc-io/laconic2d/x/auction/keeper"
)

// EndBlocker is called every block
func EndBlocker(ctx context.Context, k keeper.Keeper) error {
	// TODO: Implement
	// k.EndBlockerProcessAuctions(ctx)

	return nil
}
