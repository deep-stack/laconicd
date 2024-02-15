package keeper

// import (
// 	registrytypes "git.vdb.to/cerc-io/laconic2d/x/registry"
// )

// var _ registrytypes.QueryServer = queryServer{}

type queryServer struct {
	k Keeper
}

// // NewQueryServerImpl returns an implementation of the module QueryServer.
// func NewQueryServerImpl(k Keeper) registrytypes.QueryServer {
// 	return queryServer{k}
// }
