package auction

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgCreateAuction{}
	_ sdk.Msg = &MsgCommitBid{}
	_ sdk.Msg = &MsgRevealBid{}
)

// NewMsgCreateAuction is the constructor function for MsgCreateAuction.
func NewMsgCreateAuction(params Params, signer sdk.AccAddress) MsgCreateAuction {
	return MsgCreateAuction{
		CommitsDuration: params.CommitsDuration,
		RevealsDuration: params.RevealsDuration,
		CommitFee:       params.CommitFee,
		RevealFee:       params.RevealFee,
		MinimumBid:      params.MinimumBid,
		Signer:          signer.String(),
	}
}

// NewMsgCommitBid is the constructor function for MsgCommitBid.
func NewMsgCommitBid(auctionId string, commitHash string, signer sdk.AccAddress) MsgCommitBid {
	return MsgCommitBid{
		AuctionId:  auctionId,
		CommitHash: commitHash,
		Signer:     signer.String(),
	}
}

// ValidateBasic Implements Msg.
func (msg MsgCommitBid) ValidateBasic() error {
	if msg.Signer == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer address")
	}

	if msg.AuctionId == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid auction id")
	}

	if msg.CommitHash == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid commit hash")
	}

	return nil
}

// NewMsgRevealBid is the constructor function for MsgRevealBid.
func NewMsgRevealBid(auctionId string, reveal string, signer sdk.AccAddress) MsgRevealBid {
	return MsgRevealBid{
		AuctionId: auctionId,
		Reveal:    reveal,
		Signer:    signer.String(),
	}
}

// ValidateBasic Implements Msg.
func (msg MsgRevealBid) ValidateBasic() error {
	if msg.Signer == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer address")
	}

	if msg.AuctionId == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid auction id")
	}

	if msg.Reveal == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid reveal data")
	}

	return nil
}
