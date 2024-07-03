package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"git.vdb.to/cerc-io/laconicd/utils"
	"git.vdb.to/cerc-io/laconicd/x/onboarding"
)

type msgServer struct {
	k Keeper
}

var _ onboarding.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper *Keeper) onboarding.MsgServer {
	return &msgServer{k: *keeper}
}

func (ms msgServer) OnboardParticipant(c context.Context, msg *onboarding.MsgOnboardParticipant) (*onboarding.MsgOnboardParticipantResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	ctx = *utils.CtxWithCustomKVGasConfig(&ctx)

	signerAddress, err := sdk.AccAddressFromBech32(msg.Participant)
	if err != nil {
		return nil, err
	}

	_, err = ms.k.OnboardParticipant(ctx, msg, signerAddress)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			onboarding.EventTypeOnboardParticipant,
			sdk.NewAttribute(onboarding.AttributeKeySigner, msg.Participant),
			sdk.NewAttribute(onboarding.AttributeKeyEthAddress, msg.EthPayload.Address),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, onboarding.AttributeValueCategory),
			sdk.NewAttribute(onboarding.AttributeKeySigner, msg.Participant),
		),
	})

	utils.LogTxGasConsumed(ctx, ms.k.Logger(ctx), "OnboardParticipant")

	return &onboarding.MsgOnboardParticipantResponse{}, nil
}
