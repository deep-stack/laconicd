package registry

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetName{},
		&MsgReserveAuthority{},
		&MsgDeleteName{},
		&MsgSetAuthorityBond{},

		&MsgSetRecord{},
		&MsgRenewRecord{},
		&MsgAssociateBond{},
		&MsgDissociateBond{},
		&MsgDissociateRecords{},
		&MsgReassociateRecords{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
