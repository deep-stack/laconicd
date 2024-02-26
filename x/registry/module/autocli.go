package module

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"cosmossdk.io/client/v2/autocli"

	registryv1 "git.vdb.to/cerc-io/laconic2d/api/cerc/registry/v1"
)

var _ autocli.HasAutoCLIConfig = AppModule{}

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: registryv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod:      "Params",
					Use:            "params",
					Short:          "Get the current registry parameters",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
				{
					RpcMethod:      "Records",
					Use:            "list",
					Short:          "List records",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
				{
					RpcMethod: "GetRecord",
					Use:       "get [record-id]",
					Short:     "Get record info by record id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "GetRecordsByBondId",
					Use:       "get-records-by-bond-id [bond-id]",
					Short:     "Get records by bond id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "Whois",
					Use:       "whois [name]",
					Short:     "Get name authority info",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "name"},
					},
				},
				{
					RpcMethod:      "NameRecords",
					Use:            "names",
					Short:          "List name records",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
				{
					RpcMethod: "LookupCrn",
					Use:       "lookup [crn]",
					Short:     "Get naming info for CRN",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "crn"},
					},
				},
				{
					RpcMethod: "ResolveCrn",
					Use:       "resolve [crn]",
					Short:     "Resolve CRN to record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "crn"},
					},
				},
				{
					RpcMethod:      "GetRegistryModuleBalance",
					Use:            "balance",
					Short:          "Get registry module account balances",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: registryv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "RenewRecord",
					Use:       "renew-record [record-id]",
					Short:     "Renew (expired) record",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "record_id"},
					},
				},
				{
					RpcMethod: "ReserveName",
					Use:       "reserve-name [name] [owner]",
					Short:     "Reserve name",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "name"},
						{ProtoField: "owner"},
					},
				},
				{
					RpcMethod: "SetAuthorityBond",
					Use:       "authority-bond [name] [bond-id]",
					Short:     "Associate authority with bond",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "name"},
						{ProtoField: "bond_id"},
					},
				},
				{
					RpcMethod: "SetName",
					Use:       "set-name [crn] [cid]",
					Short:     "Set CRN to CID mapping",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "crn"},
						{ProtoField: "cid"},
					},
				},
				{
					RpcMethod: "DeleteName",
					Use:       "delete-name [crn]",
					Short:     "Delete CRN",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "crn"},
					},
				},
				{
					RpcMethod: "AssociateBond",
					Use:       "associate-bond [record-id] [bond-id]",
					Short:     "Associate record with a bond",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "record_id"},
						{ProtoField: "bond_id"},
					},
				},
				{
					RpcMethod: "DissociateBond",
					Use:       "dissociate-bond [record-id]",
					Short:     "Dissociate record from (existing) bond",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "record_id"},
					},
				},
				{
					RpcMethod: "DissociateRecords",
					Use:       "dissociate-records [bond-id]",
					Short:     "Dissociate all records from a bond",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "bond_id"},
					},
				},
				{
					RpcMethod: "ReassociateRecords",
					Use:       "reassociate-records [old-bond-id] [new-bond-id]",
					Short:     "Re-associate all records from an old to a new bond",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "old_bond_id"},
						{ProtoField: "new_bond_id"},
					},
				},
			},
			EnhanceCustomCommand: true, // Allow additional manual commands
		},
	}
}
