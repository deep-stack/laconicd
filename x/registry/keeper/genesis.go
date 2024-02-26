package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"git.vdb.to/cerc-io/laconic2d/x/registry"
)

// InitGenesis initializes the module state from a genesis state.
func (k *Keeper) InitGenesis(ctx sdk.Context, data *registry.GenesisState) error {
	if err := k.Params.Set(ctx, data.Params); err != nil {
		return err
	}

	for _, record := range data.Records {
		if err := k.SaveRecord(ctx, record); err != nil {
			return err
		}

		// Add to record expiry queue if expiry time is in the future.
		expiryTime, err := time.Parse(time.RFC3339, record.ExpiryTime)
		if err != nil {
			panic(err)
		}

		if expiryTime.After(ctx.BlockTime()) {
			// TODO
			// k.InsertRecordExpiryQueue(ctx, record)
		}
	}

	for _, authority := range data.Authorities {
		// Only import authorities that are marked active.
		if authority.Entry.Status == registry.AuthorityActive {
			if err := k.SaveNameAuthority(ctx, authority.Name, authority.Entry); err != nil {
				return err
			}

			// TODO
			// Add authority name to expiry queue.
			// k.InsertAuthorityExpiryQueue(ctx, authority.Name, authority.Entry.ExpiryTime)

			// TODO
			// Note: Bond genesis runs first, so bonds will already be present.
			// if authority.Entry.BondId != "" {
			// 	k.AddBondToAuthorityIndexEntry(ctx, authority.Entry.BondId, authority.Name)
			// }
		}
	}

	for _, nameEntry := range data.Names {
		if err := k.SaveNameRecord(ctx, nameEntry.Name, nameEntry.Entry.Latest.Id); err != nil {
			return err
		}
	}

	return nil
}

// ExportGenesis exports the module state to a genesis state.
func (k *Keeper) ExportGenesis(ctx sdk.Context) (*registry.GenesisState, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	records, err := k.ListRecords(ctx)
	if err != nil {
		return nil, err
	}

	authorities, err := k.ListNameAuthorityRecords(ctx)
	if err != nil {
		return nil, err
	}

	authorityEntries := []registry.AuthorityEntry{}
	// #nosec G705
	for name, record := range authorities {
		authorityEntries = append(authorityEntries, registry.AuthorityEntry{
			Name:  name,
			Entry: &record, //nolint: all
		}) // #nosec G601
	}

	names, err := k.ListNameRecords(ctx)
	if err != nil {
		return nil, err
	}

	return &registry.GenesisState{
		Params:      params,
		Records:     records,
		Authorities: authorityEntries,
		Names:       names,
	}, nil
}
