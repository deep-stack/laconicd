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
	onboardingTypes "git.vdb.to/cerc-io/laconicd/x/onboarding"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	// authority is the address capable of executing a MsgUpdateParams and other authority-gated message.
	// typically, this should be the x/gov module account.
	authority string

	// state management
	Schema       collections.Schema
	Params       collections.Item[onboardingTypes.Params]
	Participants collections.Map[string, onboardingTypes.Participant]
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
		Params:       collections.NewItem(sb, onboardingTypes.ParamsPrefix, "params", codec.CollValue[onboardingTypes.Params](cdc)),
		Participants: collections.NewMap(
			sb, onboardingTypes.ParticipantsPrefix, "participants", collections.StringKey, codec.CollValue[onboardingTypes.Participant](cdc),
		),
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
	return ctx.Logger().With("module", onboardingTypes.ModuleName)
}

func (k Keeper) OnboardParticipant(
	ctx sdk.Context,
	msg *onboardingTypes.MsgOnboardParticipant,
	signerAddress sdk.AccAddress,
) (*onboardingTypes.MsgOnboardParticipantResponse, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	if !params.OnboardingEnabled {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Validator onboarding is disabled")
	}

	message, err := json.Marshal(msg.EthPayload)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Invalid format for payload")
	}

	ethereumAddress, err := utils.DecodeEthereumAddress(message, msg.EthSignature)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Failed to decode Ethereum address")
	}

	if ethereumAddress != msg.EthPayload.Address {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Recovered ethereum address does not match the address set in payload")
	}

	participant := &onboardingTypes.Participant{
		CosmosAddress:   signerAddress.String(),
		EthereumAddress: ethereumAddress,
	}

	if err := k.StoreParticipant(ctx, participant); err != nil {
		return nil, err
	}

	return &onboardingTypes.MsgOnboardParticipantResponse{}, nil
}

func (k Keeper) StoreParticipant(ctx sdk.Context, participant *onboardingTypes.Participant) error {
	key := participant.CosmosAddress
	return k.Participants.Set(ctx, key, *participant)
}

// ListParticipants - get all participants.
func (k Keeper) ListParticipants(ctx sdk.Context) ([]*onboardingTypes.Participant, error) {
	var participants []*onboardingTypes.Participant

	iter, err := k.Participants.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	for ; iter.Valid(); iter.Next() {
		participant, err := iter.Value()
		if err != nil {
			return nil, err
		}

		participants = append(participants, &participant)
	}

	return participants, nil
}
