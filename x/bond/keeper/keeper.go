package keeper

import (
	"cosmossdk.io/core/address"
	"github.com/cosmos/cosmos-sdk/codec"
)

// TODO: Add genesis.go?

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	// TODO: state management

	// TODO: Later: add bond usage keepers
}

// NewKeeper creates a new Keeper instance
// TODO: Implement
func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec) Keeper {
	return Keeper{}
}
