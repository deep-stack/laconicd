package module

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"cosmossdk.io/client/v2/autocli"
	// registryv1 "git.vdb.to/cerc-io/laconic2d/api/cerc/registry/v1"
)

var _ autocli.HasAutoCLIConfig = AppModule{}

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: nil,
		Tx:    nil,
	}
}
