package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/bond interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSetName{}, "registry/SetName", nil)
	cdc.RegisterConcrete(&MsgReserveAuthority{}, "registry/ReserveAuthority", nil)
	cdc.RegisterConcrete(&MsgDeleteNameAuthority{}, "registry/DeleteAuthority", nil)
	cdc.RegisterConcrete(&MsgSetAuthorityBond{}, "registry/SetAuthorityBond", nil)

	cdc.RegisterConcrete(&MsgSetRecord{}, "registry/SetRecord", nil)
	cdc.RegisterConcrete(&MsgRenewRecord{}, "registry/RenewRecord", nil)
	cdc.RegisterConcrete(&MsgAssociateBond{}, "registry/AssociateBond", nil)
	cdc.RegisterConcrete(&MsgDissociateBond{}, "registry/DissociateBond", nil)
	cdc.RegisterConcrete(&MsgDissociateRecords{}, "registry/DissociateRecords", nil)
	cdc.RegisterConcrete(&MsgReAssociateRecords{}, "registry/ReassociateRecords", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetName{},
		&MsgReserveAuthority{},
		&MsgDeleteNameAuthority{},
		&MsgSetAuthorityBond{},

		&MsgSetRecord{},
		&MsgRenewRecord{},
		&MsgAssociateBond{},
		&MsgDissociateBond{},
		&MsgDissociateRecords{},
		&MsgReAssociateRecords{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
