package keeper

import (
	"errors"

	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"
	storetypes "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	auctiontypes "git.vdb.to/cerc-io/laconic2d/x/auction"
)

type AuctionsIndexes struct {
	Owner *indexes.Multi[string, string, auctiontypes.Auction]
}

func (a AuctionsIndexes) IndexesList() []collections.Index[string, auctiontypes.Auction] {
	return []collections.Index[string, auctiontypes.Auction]{a.Owner}
}

func newAuctionIndexes(sb *collections.SchemaBuilder) AuctionsIndexes {
	return AuctionsIndexes{
		Owner: indexes.NewMulti(
			sb, auctiontypes.AuctionOwnerIndexPrefix, "auctions_by_owner",
			collections.StringKey, collections.StringKey,
			func(_ string, v auctiontypes.Auction) (string, error) {
				return v.OwnerAddress, nil
			},
		),
	}
}

// TODO: Add required methods

type Keeper struct {
	// Codecs
	cdc codec.BinaryCodec

	// External keepers
	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper

	// Track auction usage in other cosmos-sdk modules (more like a usage tracker).
	// usageKeepers []types.AuctionUsageKeeper

	// state management
	Schema   collections.Schema
	Params   collections.Item[auctiontypes.Params]
	Auctions *collections.IndexedMap[string, auctiontypes.Auction, AuctionsIndexes]
}

// NewKeeper creates a new Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	accountKeeper auth.AccountKeeper,
	bankKeeper bank.Keeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:           cdc,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		Params:        collections.NewItem(sb, auctiontypes.ParamsKeyPrefix, "params", codec.CollValue[auctiontypes.Params](cdc)),
		Auctions:      collections.NewIndexedMap(sb, auctiontypes.AuctionsKeyPrefix, "auctions", collections.StringKey, codec.CollValue[auctiontypes.Auction](cdc), newAuctionIndexes(sb)),
		// usageKeepers:  usageKeepers,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

// func (k *Keeper) SetUsageKeepers(usageKeepers []types.AuctionUsageKeeper) {
// 	k.usageKeepers = usageKeepers
// }

// SaveAuction - saves a auction to the store.
func (k Keeper) SaveAuction(ctx sdk.Context, auction *auctiontypes.Auction) error {
	return k.Auctions.Set(ctx, auction.Id, *auction)

	// // Notify interested parties.
	// for _, keeper := range k.usageKeepers {
	// 	keeper.OnAuction(ctx, auction.Id)
	// }
	// return nil
}

// ListAuctions - get all auctions.
func (k Keeper) ListAuctions(ctx sdk.Context) ([]auctiontypes.Auction, error) {
	var auctions []auctiontypes.Auction

	iter, err := k.Auctions.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	for ; iter.Valid(); iter.Next() {
		auction, err := iter.Value()
		if err != nil {
			return nil, err
		}

		auctions = append(auctions, auction)
	}

	return auctions, nil
}

// GetAuction - gets a record from the store.
func (k Keeper) GetAuctionById(ctx sdk.Context, id string) (auctiontypes.Auction, error) {
	auction, err := k.Auctions.Get(ctx, id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return auctiontypes.Auction{}, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Auction not found.")
		}
		return auctiontypes.Auction{}, err
	}

	return auction, nil
}

func (k Keeper) GetAuctionsByOwner(ctx sdk.Context, owner string) ([]auctiontypes.Auction, error) {
	iter, err := k.Auctions.Indexes.Owner.MatchExact(ctx, owner)
	if err != nil {
		return []auctiontypes.Auction{}, err
	}

	return indexes.CollectValues(ctx, k.Auctions, iter)
}

// CreateAuction creates a new auction.
func (k Keeper) CreateAuction(ctx sdk.Context, msg auctiontypes.MsgCreateAuction) (*auctiontypes.Auction, error) {
	// TODO: Setup checks
	// Might be called from another module directly, always validate.
	// err := msg.ValidateBasic()
	// if err != nil {
	// 	return nil, err
	// }

	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	// Generate auction Id.
	account := k.accountKeeper.GetAccount(ctx, signerAddress)
	if account == nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "Account not found.")
	}

	auctionId := auctiontypes.AuctionId{
		Address:  signerAddress,
		AccNum:   account.GetAccountNumber(),
		Sequence: account.GetSequence(),
	}.Generate()

	// Compute timestamps.
	now := ctx.BlockTime()
	commitsEndTime := now.Add(msg.CommitsDuration)
	revealsEndTime := now.Add(msg.CommitsDuration + msg.RevealsDuration)

	auction := auctiontypes.Auction{
		Id:             auctionId,
		Status:         auctiontypes.AuctionStatusCommitPhase,
		OwnerAddress:   signerAddress.String(),
		CreateTime:     now,
		CommitsEndTime: commitsEndTime,
		RevealsEndTime: revealsEndTime,
		CommitFee:      msg.CommitFee,
		RevealFee:      msg.RevealFee,
		MinimumBid:     msg.MinimumBid,
	}

	// Save auction in store.
	k.SaveAuction(ctx, &auction)

	return &auction, nil
}

func (k Keeper) CommitBid(ctx sdk.Context, msg auctiontypes.MsgCommitBid) (*auctiontypes.Bid, error) {
	panic("unimplemented")
}

func (k Keeper) RevealBid(ctx sdk.Context, msg auctiontypes.MsgRevealBid) (*auctiontypes.Auction, error) {
	panic("unimplemented")
}

// GetParams gets the auction module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (*auctiontypes.Params, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &params, nil
}
