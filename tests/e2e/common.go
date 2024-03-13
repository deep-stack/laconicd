package e2e

import (
	"fmt"
	"os"

	"cosmossdk.io/log"
	pruningtypes "cosmossdk.io/store/pruning/types"

	dbm "github.com/cosmos/cosmos-db"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"

	laconicApp "git.vdb.to/cerc-io/laconic2d/app"
	auctionmodule "git.vdb.to/cerc-io/laconic2d/x/auction/module"
	bondmodule "git.vdb.to/cerc-io/laconic2d/x/bond/module"
	registrymodule "git.vdb.to/cerc-io/laconic2d/x/registry/module"

	_ "git.vdb.to/cerc-io/laconic2d/app/params" // import for side-effects (see init)
	"git.vdb.to/cerc-io/laconic2d/testutil/network"
)

// NewTestNetworkFixture returns a new LaconicApp AppConstructor for network simulation tests
func NewTestNetworkFixture() network.TestFixture {
	dir, err := os.MkdirTemp("", "laconic")
	if err != nil {
		panic(fmt.Sprintf("failed creating temporary directory: %v", err))
	}
	defer os.RemoveAll(dir)

	app, err := laconicApp.NewLaconicApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.NewAppOptionsWithFlagHome(dir))
	if err != nil {
		panic(fmt.Sprintf("failed to create laconic app: %v", err))
	}

	appCtr := func(val network.ValidatorI) servertypes.Application {
		app, err := laconicApp.NewLaconicApp(
			val.GetCtx().Logger, dbm.NewMemDB(), nil, true,
			simtestutil.NewAppOptionsWithFlagHome(val.GetCtx().Config.RootDir),
			bam.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
			bam.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
			bam.SetChainID(val.GetCtx().Viper.GetString(flags.FlagChainID)),
		)
		if err != nil {
			panic(fmt.Sprintf("failed creating temporary directory: %v", err))
		}

		return app
	}

	return network.TestFixture{
		AppConstructor: appCtr,
		GenesisState:   app.DefaultGenesis(),
		EncodingConfig: testutil.MakeTestEncodingConfig(
			auth.AppModuleBasic{},
			bank.AppModuleBasic{},
			staking.AppModuleBasic{},
			auctionmodule.AppModule{},
			bondmodule.AppModule{},
			registrymodule.AppModule{},
		),
	}
}
