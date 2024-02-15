package keeper

// import "git.vdb.to/cerc-io/laconic2d/x/registry"

// var _ registry.MsgServer = msgServer{}

type msgServer struct {
	k Keeper
}

// // NewMsgServerImpl returns an implementation of the module MsgServer interface.
// func NewMsgServerImpl(keeper Keeper) registry.MsgServer {
// 	return &msgServer{k: keeper}
// }
