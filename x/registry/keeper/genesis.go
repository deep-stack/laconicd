package keeper

import (
	"context"

	"git.vdb.to/cerc-io/laconic2d/x/registry"
)

// InitGenesis initializes the module state from a genesis state.
func (k *Keeper) InitGenesis(ctx context.Context, data *registry.GenesisState) error {
	if err := k.Params.Set(ctx, data.Params); err != nil {
		return err
	}

	// for _, record := range data.Records {
	// 	keeper.PutRecord(ctx, record)

	// 	// Add to record expiry queue if expiry time is in the future.
	// 	expiryTime, err := time.Parse(time.RFC3339, record.ExpiryTime)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	if expiryTime.After(ctx.BlockTime()) {
	// 		keeper.InsertRecordExpiryQueue(ctx, record)
	// 	}

	// 	// Note: Bond genesis runs first, so bonds will already be present.
	// 	if record.BondId != "" {
	// 		keeper.AddBondToRecordIndexEntry(ctx, record.BondId, record.Id)
	// 	}
	// }

	// for _, authority := range data.Authorities {
	// 	// Only import authorities that are marked active.
	// 	if authority.Entry.Status == types.AuthorityActive {
	// 		keeper.SetNameAuthority(ctx, authority.Name, authority.Entry)

	// 		// Add authority name to expiry queue.
	// 		keeper.InsertAuthorityExpiryQueue(ctx, authority.Name, authority.Entry.ExpiryTime)

	// 		// Note: Bond genesis runs first, so bonds will already be present.
	// 		if authority.Entry.BondId != "" {
	// 			keeper.AddBondToAuthorityIndexEntry(ctx, authority.Entry.BondId, authority.Name)
	// 		}
	// 	}
	// }

	// for _, nameEntry := range data.Names {
	// 	keeper.SetNameRecord(ctx, nameEntry.Name, nameEntry.Entry.Latest.Id)
	// }

	return nil
}

// ExportGenesis exports the module state to a genesis state.
func (k *Keeper) ExportGenesis(ctx context.Context) (*registry.GenesisState, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	// records := keeper.ListRecords(ctx)

	// authorities := keeper.ListNameAuthorityRecords(ctx)
	// authorityEntries := []types.AuthorityEntry{}
	// // #nosec G705
	// for name, record := range authorities {
	// 	authorityEntries = append(authorityEntries, types.AuthorityEntry{
	// 		Name:  name,
	// 		Entry: &record, //nolint: all
	// 	}) // #nosec G601
	// }

	// names := keeper.ListNameRecords(ctx)

	return &registry.GenesisState{
		Params: params,
	}, nil
}
