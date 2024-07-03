package keeper

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"git.vdb.to/cerc-io/laconicd/utils"
	"git.vdb.to/cerc-io/laconicd/x/onboarding"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	// authority is the address capable of executing a MsgUpdateParams and other authority-gated message.
	// typically, this should be the x/gov module account.
	authority string

	// state management
	Schema       collections.Schema
	Params       collections.Item[onboarding.Params]
	Participants collections.Map[string, onboarding.Participant]
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, authority string) *Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,
		Params:       collections.NewItem(sb, onboarding.ParamsPrefix, "params", codec.CollValue[onboarding.Params](cdc)),
		Participants: collections.NewMap(sb, onboarding.ParticipantsPrefix, "participants", collections.StringKey, codec.CollValue[onboarding.Participant](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return &k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", onboarding.ModuleName)
}

func (k Keeper) OnboardParticipant(ctx sdk.Context, msg *onboarding.MsgOnboardParticipant, signerAddress sdk.AccAddress) (*onboarding.MsgOnboardParticipantResponse, error) {
	message, err := json.Marshal(msg.EthPayload)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Invalid format for payload")
	}

	ethereumAddress, err := utils.DecodeEthereumAddress(message, msg.EthSignature)
	if ethereumAddress != msg.EthPayload.Address {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Recovered ethereum address does not match the address set in payload")
	}

	participant := &onboarding.Participant{
		CosmosAddress:   signerAddress.String(),
		EthereumAddress: ethereumAddress,
	}

	if err := k.StoreParticipant(ctx, participant); err != nil {
		return nil, err
	}

	return nil, err
}

func (k Keeper) StoreParticipant(ctx sdk.Context, participant *onboarding.Participant) error {
	key := participant.CosmosAddress
	k.Participants.Set(ctx, key, *participant)

	return nil
}
