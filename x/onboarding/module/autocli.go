package module

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	onboardingv1 "git.vdb.to/cerc-io/laconicd/api/cerc/onboarding/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: onboardingv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod:      "Participants",
					Use:            "list",
					Short:          "List participants",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
			},
		},
		// TODO: Use JSON file for input
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: onboardingv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "OnboardParticipant",
					Use:       "enroll",
					Short:     "Enroll a testnet validator",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "eth_payload"},
						{ProtoField: "eth_signature"},
					},
				},
			},
		},
	}
}
