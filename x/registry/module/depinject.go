package module

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"

	"github.com/cosmos/cosmos-sdk/codec"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	modulev1 "git.vdb.to/cerc-io/laconicd/api/cerc/registry/module/v1"
	"git.vdb.to/cerc-io/laconicd/x/auction"
	auctionkeeper "git.vdb.to/cerc-io/laconicd/x/auction/keeper"
	"git.vdb.to/cerc-io/laconicd/x/bond"
	bondkeeper "git.vdb.to/cerc-io/laconicd/x/bond/keeper"
	"git.vdb.to/cerc-io/laconicd/x/registry/keeper"
)

var _ appmodule.AppModule = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

func init() {
	appmodule.Register(
		&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Cdc          codec.Codec
	StoreService store.KVStoreService

	AccountKeeper auth.AccountKeeper
	BankKeeper    bank.Keeper

	BondKeeper    *bondkeeper.Keeper
	AuctionKeeper *auctionkeeper.Keeper
}

type ModuleOutputs struct {
	depinject.Out

	Keeper keeper.Keeper
	Module appmodule.AppModule

	AuctionHooks auction.AuctionHooksWrapper
	BondHooks    bond.BondHooksWrapper
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.AccountKeeper,
		in.BankKeeper,
		in.BondKeeper,
		in.AuctionKeeper,
	)
	m := NewAppModule(in.Cdc, k)

	recordKeeper := keeper.NewRecordKeeper(in.Cdc, &k, in.AuctionKeeper)

	return ModuleOutputs{
		Module: m, Keeper: k,
		AuctionHooks: auction.AuctionHooksWrapper{AuctionUsageKeeper: recordKeeper},
		BondHooks:    bond.BondHooksWrapper{BondUsageKeeper: recordKeeper},
	}
}
