package bond

import (
	"github.com/cosmos/cosmos-sdk/codec"
	types "github.com/cosmos/cosmos-sdk/codec/types"
)

// RegisterLegacyAminoCodec registers the necessary x/bond interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// cdc.RegisterConcrete(&MsgCreateBond{}, "bond/MsgCreateBond", nil)
	// cdc.RegisterConcrete(&MsgRefillBond{}, "bond/MsgRefillBond", nil)
	// cdc.RegisterConcrete(&MsgWithdrawBond{}, "bond/MsgWithdrawBond", nil)
	// cdc.RegisterConcrete(&MsgCancelBond{}, "bond/MsgCancelBond", nil)
}

// RegisterInterfaces registers the interfaces types with the interface registry.
func RegisterInterfaces(registry types.InterfaceRegistry) {
	// registry.RegisterImplementations((*sdk.Msg)(nil),
	// 	&MsgCreateBond{},
	// 	&MsgRefillBond{},
	// 	&MsgCancelBond{},
	// 	&MsgWithdrawBond{},
	// )
	// msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
