package module

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"golang.org/x/exp/maps"

	"github.com/cosmos/cosmos-sdk/codec"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	modulev1 "git.vdb.to/cerc-io/laconicd/api/cerc/bond/module/v1"
	"git.vdb.to/cerc-io/laconicd/x/bond"
	"git.vdb.to/cerc-io/laconicd/x/bond/keeper"
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
		appmodule.Invoke(InvokeSetBondHooks),
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

	Keeper *keeper.Keeper
	Module appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	k := keeper.NewKeeper(in.Cdc, in.StoreService, in.AccountKeeper, in.BankKeeper)
	m := NewAppModule(in.Cdc, k)

	return ModuleOutputs{Module: m, Keeper: k}
}

func InvokeSetBondHooks(
	config *modulev1.Module,
	keeper *keeper.Keeper,
	bondHooks map[string]bond.BondHooksWrapper,
) error {
	// all arguments to invokers are optional
	if keeper == nil || config == nil {
		return nil
	}

	var usageKeepers []bond.BondUsageKeeper

	for _, modName := range maps.Keys(bondHooks) {
		hook := bondHooks[modName]
		usageKeepers = append(usageKeepers, hook)
	}

	keeper.SetUsageKeepers(usageKeepers)

	return nil
}
