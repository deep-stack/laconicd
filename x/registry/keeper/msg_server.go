package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	registrytypes "git.vdb.to/cerc-io/laconic2d/x/registry"
)

var _ registrytypes.MsgServer = msgServer{}

type msgServer struct {
	k Keeper
}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper Keeper) registrytypes.MsgServer {
	return &msgServer{k: keeper}
}

func (ms msgServer) SetRecord(c context.Context, msg *registrytypes.MsgSetRecord) (*registrytypes.MsgSetRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	record, err := ms.k.SetRecord(ctx, *msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			registrytypes.EventTypeSetRecord,
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.GetSigner()),
			sdk.NewAttribute(registrytypes.AttributeKeyBondId, msg.GetBondId()),
			sdk.NewAttribute(registrytypes.AttributeKeyPayload, msg.Payload.Record.Id),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, registrytypes.AttributeValueCategory),
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
		),
	})

	return &registrytypes.MsgSetRecordResponse{Id: record.Id}, nil
}

// nolint: all
func (ms msgServer) SetName(c context.Context, msg *registrytypes.MsgSetName) (*registrytypes.MsgSetNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	err = ms.k.SetName(ctx, *msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			registrytypes.EventTypeSetRecord,
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(registrytypes.AttributeKeyCRN, msg.Crn),
			sdk.NewAttribute(registrytypes.AttributeKeyCID, msg.Cid),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, registrytypes.AttributeValueCategory),
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
		),
	})
	return &registrytypes.MsgSetNameResponse{}, nil
}

func (ms msgServer) ReserveName(c context.Context, msg *registrytypes.MsgReserveAuthority) (*registrytypes.MsgReserveAuthorityResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	_, err = sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	err = ms.k.ReserveAuthority(ctx, *msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			registrytypes.EventTypeReserveNameAuthority,
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(registrytypes.AttributeKeyName, msg.Name),
			sdk.NewAttribute(registrytypes.AttributeKeyOwner, msg.Owner),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, registrytypes.AttributeValueCategory),
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
		),
	})
	return &registrytypes.MsgReserveAuthorityResponse{}, nil
}

// nolint: all
func (ms msgServer) SetAuthorityBond(c context.Context, msg *registrytypes.MsgSetAuthorityBond) (*registrytypes.MsgSetAuthorityBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	err = ms.k.SetAuthorityBond(ctx, *msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			registrytypes.EventTypeAuthorityBond,
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(registrytypes.AttributeKeyName, msg.Name),
			sdk.NewAttribute(registrytypes.AttributeKeyBondId, msg.BondId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, registrytypes.AttributeValueCategory),
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
		),
	})
	return &registrytypes.MsgSetAuthorityBondResponse{}, nil
}

func (ms msgServer) DeleteName(c context.Context, msg *registrytypes.MsgDeleteNameAuthority) (*registrytypes.MsgDeleteNameAuthorityResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	err = ms.k.DeleteName(ctx, *msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			registrytypes.EventTypeDeleteName,
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(registrytypes.AttributeKeyCRN, msg.Crn),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, registrytypes.AttributeValueCategory),
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
		),
	})
	return &registrytypes.MsgDeleteNameAuthorityResponse{}, nil
}

func (ms msgServer) RenewRecord(c context.Context, msg *registrytypes.MsgRenewRecord) (*registrytypes.MsgRenewRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	err = ms.k.RenewRecord(ctx, *msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			registrytypes.EventTypeRenewRecord,
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(registrytypes.AttributeKeyRecordId, msg.RecordId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, registrytypes.AttributeValueCategory),
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
		),
	})
	return &registrytypes.MsgRenewRecordResponse{}, nil
}

// nolint: all
func (ms msgServer) AssociateBond(c context.Context, msg *registrytypes.MsgAssociateBond) (*registrytypes.MsgAssociateBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	err = ms.k.ProcessAssociateBond(ctx, *msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			registrytypes.EventTypeAssociateBond,
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(registrytypes.AttributeKeyRecordId, msg.RecordId),
			sdk.NewAttribute(registrytypes.AttributeKeyBondId, msg.BondId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, registrytypes.AttributeValueCategory),
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
		),
	})
	return &registrytypes.MsgAssociateBondResponse{}, nil
}

func (ms msgServer) DissociateBond(c context.Context, msg *registrytypes.MsgDissociateBond) (*registrytypes.MsgDissociateBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	err = ms.k.ProcessDissociateBond(ctx, *msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			registrytypes.EventTypeDissociateBond,
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(registrytypes.AttributeKeyRecordId, msg.RecordId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, registrytypes.AttributeValueCategory),
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
		),
	})
	return &registrytypes.MsgDissociateBondResponse{}, nil
}

func (ms msgServer) DissociateRecords(c context.Context, msg *registrytypes.MsgDissociateRecords) (*registrytypes.MsgDissociateRecordsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	err = ms.k.ProcessDissociateRecords(ctx, *msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			registrytypes.EventTypeDissociateRecords,
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(registrytypes.AttributeKeyBondId, msg.BondId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, registrytypes.AttributeValueCategory),
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
		),
	})
	return &registrytypes.MsgDissociateRecordsResponse{}, nil
}

func (ms msgServer) ReAssociateRecords(c context.Context, msg *registrytypes.MsgReAssociateRecords) (*registrytypes.MsgReAssociateRecordsResponse, error) { //nolint: all
	ctx := sdk.UnwrapSDKContext(c)
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	err = ms.k.ProcessReAssociateRecords(ctx, *msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			registrytypes.EventTypeReAssociateRecords,
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(registrytypes.AttributeKeyOldBondId, msg.OldBondId),
			sdk.NewAttribute(registrytypes.AttributeKeyNewBondId, msg.NewBondId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, registrytypes.AttributeValueCategory),
			sdk.NewAttribute(registrytypes.AttributeKeySigner, msg.Signer),
		),
	})
	return &registrytypes.MsgReAssociateRecordsResponse{}, nil
}
