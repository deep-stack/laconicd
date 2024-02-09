package keeper

// TODO: Add required read methods

// var _ auctiontypes.QueryServer = queryServer{}

type queryServer struct {
	k Keeper
}

// NewQueryServerImpl returns an implementation of the module QueryServer.
// func NewQueryServerImpl(k Keeper) auctiontypes.QueryServer {
// 	return queryServer{k}
// }
