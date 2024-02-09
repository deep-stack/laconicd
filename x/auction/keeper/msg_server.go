package keeper

// TODO: Add required read methods

// var _ auctiontypes.MsgServer = msgServer{}

type msgServer struct {
	k Keeper
}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
// func NewMsgServerImpl(keeper Keeper) auctiontypes.MsgServer {
//     return &msgServer{k: keeper}
// }
