package module

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"git.vdb.to/cerc-io/laconic2d/x/bond"
	"git.vdb.to/cerc-io/laconic2d/x/bond/keeper"
)

// TODO: Port remaining AppModule methods

var (
	_ module.AppModuleBasic = AppModule{}
	_ appmodule.AppModule   = AppModule{}
	_ module.HasGenesis     = AppModule{}
	_ module.HasServices    = AppModule{}
)

// ConsensusVersion defines the current module consensus version
const ConsensusVersion = 1

type AppModule struct {
	cdc    codec.Codec
	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
	return AppModule{
		cdc:    cdc,
		keeper: keeper,
	}
}

func NewAppModuleBasic(m AppModule) module.AppModuleBasic {
	return module.CoreAppModuleBasicAdaptor(m.Name(), m)
}

// module.AppModuleBasic

// Name returns the bond module's name.
func (AppModule) Name() string { return bond.ModuleName }

// RegisterLegacyAminoCodec registers the bond module's types on the LegacyAmino codec.
// New modules do not need to support Amino.
func (AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	bond.RegisterLegacyAminoCodec(cdc)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the bond module.
func (AppModule) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	// if err := bond.RegisterQueryHandlerClient(context.Background(), mux, bond.NewQueryClient(clientCtx)); err != nil {
	// 	panic(err)
	// }
}

// RegisterInterfaces registers interfaces and implementations of the bond module.
func (AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	bond.RegisterInterfaces(registry)
}

// module.HasGenesis

// DefaultGenesis returns default genesis state as raw bytes for the module.
func (AppModule) DefaultGenesis(jsonCodec codec.JSONCodec) json.RawMessage {
	return jsonCodec.MustMarshalJSON(bond.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the module.
func (AppModule) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, message json.RawMessage) error {
	var data bond.GenesisState
	if err := cdc.UnmarshalJSON(message, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", bond.ModuleName, err)
	}

	return data.Validate()
}

// InitGenesis performs genesis initialization for the checkers module.
// It returns no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	var genesisState bond.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	if err := am.keeper.InitGenesis(ctx, &genesisState); err != nil {
		panic(fmt.Sprintf("failed to initialize %s genesis state: %v", bond.ModuleName, err))
	}
}

// ExportGenesis returns the exported genesis state as raw bytes for the circuit
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs, err := am.keeper.ExportGenesis(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to export %s genesis state: %v", bond.ModuleName, err))
	}

	return cdc.MustMarshalJSON(gs)
}

// module.HasServices

func (am AppModule) RegisterServices(cfg module.Configurator) {
	// Register servers
	// bond.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	// bond.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))
}
