package module

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"

	"github.com/cosmos/cosmos-sdk/codec"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	modulev1 "git.vdb.to/cerc-io/laconic2d/api/cerc/registry/module/v1"
	auctionkeeper "git.vdb.to/cerc-io/laconic2d/x/auction/keeper"
	bondkeeper "git.vdb.to/cerc-io/laconic2d/x/bond/keeper"
	"git.vdb.to/cerc-io/laconic2d/x/registry/keeper"
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

	BondKeeper    bondkeeper.Keeper
	AuctionKeeper auctionkeeper.Keeper
}

type ModuleOutputs struct {
	depinject.Out

	Keeper keeper.Keeper
	Module appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.AccountKeeper,
		in.BankKeeper,
		keeper.RecordKeeper{},
		in.BondKeeper,
		in.AuctionKeeper,
	)
	m := NewAppModule(in.Cdc, k)

	return ModuleOutputs{Module: m, Keeper: k}
}
