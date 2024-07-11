package integration_test

import (
	"context"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	cmtprototypes "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	auctionTypes "git.vdb.to/cerc-io/laconicd/x/auction"
	auctionkeeper "git.vdb.to/cerc-io/laconicd/x/auction/keeper"
	auctionmodule "git.vdb.to/cerc-io/laconicd/x/auction/module"
	bondTypes "git.vdb.to/cerc-io/laconicd/x/bond"
	bondkeeper "git.vdb.to/cerc-io/laconicd/x/bond/keeper"
	bondmodule "git.vdb.to/cerc-io/laconicd/x/bond/module"
	registryTypes "git.vdb.to/cerc-io/laconicd/x/registry"
	registrykeeper "git.vdb.to/cerc-io/laconicd/x/registry/keeper"
	registrymodule "git.vdb.to/cerc-io/laconicd/x/registry/module"
)

type TestFixture struct {
	App *integration.App

	SdkCtx sdk.Context
	cdc    codec.Codec
	keys   map[string]*storetypes.KVStoreKey

	AccountKeeper authkeeper.AccountKeeper
	BankKeeper    bankkeeper.Keeper

	AuctionKeeper  *auctionkeeper.Keeper
	BondKeeper     *bondkeeper.Keeper
	RegistryKeeper registrykeeper.Keeper
}

func (tf *TestFixture) Setup() error {
	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, auctionTypes.StoreKey, bondTypes.StoreKey, registryTypes.StoreKey,
	)
	cdc := moduletestutil.MakeTestEncodingConfig(
		auth.AppModuleBasic{},
		auctionmodule.AppModule{},
		bondmodule.AppModule{},
		registrymodule.AppModule{},
	).Codec

	logger := log.NewNopLogger() // Use log.NewTestLogger(kts.T()) for help with debugging
	cms := integration.CreateMultiStore(keys, logger)

	newCtx := sdk.NewContext(cms, cmtprototypes.Header{}, true, logger)

	authority := authtypes.NewModuleAddress("gov")

	maccPerms := map[string][]string{
		minttypes.ModuleName:                         {authtypes.Minter},
		auctionTypes.ModuleName:                      {},
		auctionTypes.AuctionBurnModuleAccountName:    {},
		bondTypes.ModuleName:                         {},
		registryTypes.ModuleName:                     {},
		registryTypes.RecordRentModuleAccountName:    {},
		registryTypes.AuthorityRentModuleAccountName: {},
	}

	accountKeeper := authkeeper.NewAccountKeeper(
		cdc,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		addresscodec.NewBech32Codec(sdk.Bech32MainPrefix),
		sdk.Bech32MainPrefix,
		authority.String(),
	)

	blockedAddresses := map[string]bool{
		accountKeeper.GetAuthority(): false,
	}
	bankKeeper := bankkeeper.NewBaseKeeper(
		cdc,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		accountKeeper,
		blockedAddresses,
		authority.String(),
		log.NewNopLogger(),
	)

	auctionKeeper := auctionkeeper.NewKeeper(cdc, runtime.NewKVStoreService(keys[auctionTypes.StoreKey]), accountKeeper, bankKeeper)

	bondKeeper := bondkeeper.NewKeeper(cdc, runtime.NewKVStoreService(keys[bondTypes.StoreKey]), accountKeeper, bankKeeper)

	registryKeeper := registrykeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(keys[registryTypes.StoreKey]),
		accountKeeper,
		bankKeeper,
		bondKeeper,
		auctionKeeper,
	)

	authModule := auth.NewAppModule(cdc, accountKeeper, authsims.RandomGenesisAccounts, nil)
	bankModule := bank.NewAppModule(cdc, bankKeeper, accountKeeper, nil)
	auctionModule := auctionmodule.NewAppModule(cdc, auctionKeeper)
	bondModule := bondmodule.NewAppModule(cdc, bondKeeper)
	registryModule := registrymodule.NewAppModule(cdc, registryKeeper)

	integrationApp := integration.NewIntegrationApp(newCtx, logger, keys, cdc, map[string]appmodule.AppModule{
		authtypes.ModuleName:     authModule,
		banktypes.ModuleName:     bankModule,
		auctionTypes.ModuleName:  auctionModule,
		bondTypes.ModuleName:     bondModule,
		registryTypes.ModuleName: registryModule,
	})

	sdkCtx := sdk.UnwrapSDKContext(integrationApp.Context())

	// Register MsgServer and QueryServer
	auctionTypes.RegisterMsgServer(integrationApp.MsgServiceRouter(), auctionkeeper.NewMsgServerImpl(auctionKeeper))
	auctionTypes.RegisterQueryServer(integrationApp.QueryHelper(), auctionkeeper.NewQueryServerImpl(auctionKeeper))

	bondTypes.RegisterMsgServer(integrationApp.MsgServiceRouter(), bondkeeper.NewMsgServerImpl(bondKeeper))
	bondTypes.RegisterQueryServer(integrationApp.QueryHelper(), bondkeeper.NewQueryServerImpl(bondKeeper))

	registryTypes.RegisterMsgServer(integrationApp.MsgServiceRouter(), registrykeeper.NewMsgServerImpl(registryKeeper))
	registryTypes.RegisterQueryServer(integrationApp.QueryHelper(), registrykeeper.NewQueryServerImpl(registryKeeper))

	// set default params
	if err := auctionKeeper.Params.Set(sdkCtx, auctionTypes.DefaultParams()); err != nil {
		return err
	}
	if err := bondKeeper.Params.Set(sdkCtx, bondTypes.DefaultParams()); err != nil {
		return err
	}
	if err := registryKeeper.Params.Set(sdkCtx, registryTypes.DefaultParams()); err != nil {
		return err
	}

	tf.App = integrationApp
	tf.SdkCtx, tf.cdc, tf.keys = sdkCtx, cdc, keys
	tf.AccountKeeper, tf.BankKeeper = accountKeeper, bankKeeper
	tf.AuctionKeeper, tf.BondKeeper, tf.RegistryKeeper = auctionKeeper, bondKeeper, registryKeeper

	return nil
}

type BondDenomProvider struct{}

func (bdp BondDenomProvider) BondDenom(ctx context.Context) (string, error) {
	return sdk.DefaultBondDenom, nil
}
