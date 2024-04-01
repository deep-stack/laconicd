package module

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"golang.org/x/exp/maps"

	"github.com/cosmos/cosmos-sdk/codec"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	modulev1 "git.vdb.to/cerc-io/laconicd/api/cerc/auction/module/v1"
	"git.vdb.to/cerc-io/laconicd/x/auction"
	"git.vdb.to/cerc-io/laconicd/x/auction/keeper"
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
		appmodule.Invoke(InvokeSetAuctionHooks),
	)
}

type ModuleInputs struct {
	depinject.In

	Cdc          codec.Codec
	StoreService store.KVStoreService

	AccountKeeper auth.AccountKeeper
	BankKeeper    bank.Keeper
}

type ModuleOutputs struct {
	depinject.Out

	// Use * as required by InvokeSetAuctionHooks
	// https://github.com/cosmos/cosmos-sdk/tree/v0.50.3/core/appmodule#invoker-invocation-details
	// https://github.com/cosmos/cosmos-sdk/tree/v0.50.3/core/appmodule#regular-golang-types
	Keeper *keeper.Keeper

	Module appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	k := keeper.NewKeeper(in.Cdc, in.StoreService, in.AccountKeeper, in.BankKeeper)
	m := NewAppModule(in.Cdc, k)

	return ModuleOutputs{Module: m, Keeper: k}
}

func InvokeSetAuctionHooks(
	config *modulev1.Module,
	keeper *keeper.Keeper,
	auctionHooks map[string]auction.AuctionHooksWrapper,
) error {
	// all arguments to invokers are optional
	if keeper == nil || config == nil {
		return nil
	}

	var usageKeepers []auction.AuctionUsageKeeper

	for _, modName := range maps.Keys(auctionHooks) {
		hook := auctionHooks[modName]
		usageKeepers = append(usageKeepers, hook)
	}

	keeper.SetUsageKeepers(usageKeepers)

	return nil
}
