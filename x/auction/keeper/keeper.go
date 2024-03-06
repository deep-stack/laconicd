package keeper

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"
	storetypes "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	wnsUtils "git.vdb.to/cerc-io/laconic2d/utils"
	auctiontypes "git.vdb.to/cerc-io/laconic2d/x/auction"
)

// CompletedAuctionDeleteTimeout => Completed auctions are deleted after this timeout (after reveals end time).
const CompletedAuctionDeleteTimeout = time.Hour * 24

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

type BidsIndexes struct {
	Bidder *indexes.ReversePair[string, string, auctiontypes.Bid]
}

func (b BidsIndexes) IndexesList() []collections.Index[collections.Pair[string, string], auctiontypes.Bid] {
	return []collections.Index[collections.Pair[string, string], auctiontypes.Bid]{b.Bidder}
}

func newBidsIndexes(sb *collections.SchemaBuilder) BidsIndexes {
	return BidsIndexes{
		Bidder: indexes.NewReversePair[auctiontypes.Bid](
			sb, auctiontypes.BidderAuctionIdIndexPrefix, "auction_id_by_bidder",
			collections.PairKeyCodec(collections.StringKey, collections.StringKey),
		),
	}
}

type Keeper struct {
	// Codecs
	cdc codec.BinaryCodec

	// External keepers
	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper

	// Track auction usage in other cosmos-sdk modules (more like a usage tracker).
	usageKeepers []auctiontypes.AuctionUsageKeeper

	// state management
	Schema   collections.Schema
	Params   collections.Item[auctiontypes.Params]
	Auctions *collections.IndexedMap[string, auctiontypes.Auction, AuctionsIndexes]                   // map: auctionId -> Auction, index: owner -> Auctions
	Bids     *collections.IndexedMap[collections.Pair[string, string], auctiontypes.Bid, BidsIndexes] // map: (auctionId, bidder) -> Bid, index: bidder -> auctionId
}

// NewKeeper creates a new Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	accountKeeper auth.AccountKeeper,
	bankKeeper bank.Keeper,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:           cdc,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		Params:        collections.NewItem(sb, auctiontypes.ParamsPrefix, "params", codec.CollValue[auctiontypes.Params](cdc)),
		Auctions:      collections.NewIndexedMap(sb, auctiontypes.AuctionsPrefix, "auctions", collections.StringKey, codec.CollValue[auctiontypes.Auction](cdc), newAuctionIndexes(sb)),
		Bids:          collections.NewIndexedMap(sb, auctiontypes.BidsPrefix, "bids", collections.PairKeyCodec(collections.StringKey, collections.StringKey), codec.CollValue[auctiontypes.Bid](cdc), newBidsIndexes(sb)),
		usageKeepers:  nil,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return &k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return logger(ctx)
}

func logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", auctiontypes.ModuleName)
}

func (k *Keeper) SetUsageKeepers(usageKeepers []auctiontypes.AuctionUsageKeeper) {
	if k.usageKeepers != nil {
		panic("cannot set auction hooks twice")
	}

	k.usageKeepers = usageKeepers
}

// SaveAuction - saves a auction to the store.
func (k Keeper) SaveAuction(ctx sdk.Context, auction *auctiontypes.Auction) error {
	return k.Auctions.Set(ctx, auction.Id, *auction)
}

// DeleteAuction - deletes the auction.
func (k Keeper) DeleteAuction(ctx sdk.Context, auction auctiontypes.Auction) error {
	// Delete all bids first.
	bids, err := k.GetBids(ctx, auction.Id)
	if err != nil {
		return err
	}

	for _, bid := range bids {
		if err := k.DeleteBid(ctx, *bid); err != nil {
			return err
		}
	}

	return k.Auctions.Remove(ctx, auction.Id)
}

func (k Keeper) HasAuction(ctx sdk.Context, id string) (bool, error) {
	has, err := k.Auctions.Has(ctx, id)
	if err != nil {
		return false, err
	}

	return has, nil
}

func (k Keeper) SaveBid(ctx sdk.Context, bid *auctiontypes.Bid) error {
	key := collections.Join(bid.AuctionId, bid.BidderAddress)
	return k.Bids.Set(ctx, key, *bid)
}

func (k Keeper) DeleteBid(ctx sdk.Context, bid auctiontypes.Bid) error {
	key := collections.Join(bid.AuctionId, bid.BidderAddress)
	return k.Bids.Remove(ctx, key)
}

func (k Keeper) HasBid(ctx sdk.Context, id string, bidder string) (bool, error) {
	key := collections.Join(id, bidder)
	has, err := k.Bids.Has(ctx, key)
	if err != nil {
		return false, err
	}

	return has, nil
}

func (k Keeper) GetBid(ctx sdk.Context, id string, bidder string) (auctiontypes.Bid, error) {
	key := collections.Join(id, bidder)
	bid, err := k.Bids.Get(ctx, key)
	if err != nil {
		return auctiontypes.Bid{}, err
	}

	return bid, nil
}

// GetBids gets the auction bids.
func (k Keeper) GetBids(ctx sdk.Context, id string) ([]*auctiontypes.Bid, error) {
	var bids []*auctiontypes.Bid

	err := k.Bids.Walk(ctx, collections.NewPrefixedPairRange[string, string](id), func(key collections.Pair[string, string], value auctiontypes.Bid) (stop bool, err error) {
		bids = append(bids, &value)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return bids, nil
}

// ListAuctions - get all auctions.
func (k Keeper) ListAuctions(ctx sdk.Context) ([]auctiontypes.Auction, error) {
	iter, err := k.Auctions.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	return iter.Values()
}

// MatchAuctions - get all matching auctions.
func (k Keeper) MatchAuctions(ctx sdk.Context, matchFn func(*auctiontypes.Auction) (bool, error)) ([]*auctiontypes.Auction, error) {
	var auctions []*auctiontypes.Auction

	err := k.Auctions.Walk(ctx, nil, func(key string, value auctiontypes.Auction) (bool, error) {
		auctionMatched, err := matchFn(&value)
		if err != nil {
			return true, err
		}

		if auctionMatched {
			auctions = append(auctions, &value)
		}

		return false, nil
	})
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return indexes.CollectValues(ctx, k.Auctions, iter)
}

// QueryAuctionsByBidder - query auctions by bidder
func (k Keeper) QueryAuctionsByBidder(ctx sdk.Context, bidderAddress string) ([]auctiontypes.Auction, error) {
	auctions := []auctiontypes.Auction{}

	iter, err := k.Bids.Indexes.Bidder.MatchExact(ctx, bidderAddress)
	if err != nil {
		return nil, err
	}

	for ; iter.Valid(); iter.Next() {
		keyPair, err := iter.PrimaryKey()
		if err != nil {
			return nil, err
		}

		auction, err := k.GetAuctionById(ctx, keyPair.K1())
		if err != nil {
			return nil, err
		}

		auctions = append(auctions, auction)
	}

	return auctions, nil
}

// CreateAuction creates a new auction.
func (k Keeper) CreateAuction(ctx sdk.Context, msg auctiontypes.MsgCreateAuction) (*auctiontypes.Auction, error) {
	// Might be called from another module directly, always validate.
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

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
	if err = k.SaveAuction(ctx, &auction); err != nil {
		return nil, err
	}

	return &auction, nil
}

func (k Keeper) CommitBid(ctx sdk.Context, msg auctiontypes.MsgCommitBid) (*auctiontypes.Bid, error) {
	if has, err := k.HasAuction(ctx, msg.AuctionId); !has {
		if err != nil {
			return nil, err
		}
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Auction not found.")
	}

	auction, err := k.GetAuctionById(ctx, msg.AuctionId)
	if err != nil {
		return nil, err
	}

	if auction.Status != auctiontypes.AuctionStatusCommitPhase {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Auction is not in commit phase.")
	}

	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	// Take auction fees from account.
	totalFee := auction.CommitFee.Add(auction.RevealFee)
	sdkErr := k.bankKeeper.SendCoinsFromAccountToModule(ctx, signerAddress, auctiontypes.ModuleName, sdk.NewCoins(totalFee))
	if sdkErr != nil {
		return nil, sdkErr
	}

	// Check if an old bid already exists, if so, return old bids auction fee (update bid scenario).
	bidder := signerAddress.String()
	bidExists, err := k.HasBid(ctx, msg.AuctionId, bidder)
	if err != nil {
		return nil, err
	}

	if bidExists {
		oldBid, err := k.GetBid(ctx, msg.AuctionId, bidder)
		if err != nil {
			return nil, err
		}

		oldTotalFee := oldBid.CommitFee.Add(oldBid.RevealFee)
		sdkErr := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, signerAddress, sdk.NewCoins(oldTotalFee))
		if sdkErr != nil {
			return nil, sdkErr
		}
	}

	// Save new bid.
	bid := auctiontypes.Bid{
		AuctionId:     msg.AuctionId,
		BidderAddress: bidder,
		Status:        auctiontypes.BidStatusCommitted,
		CommitHash:    msg.CommitHash,
		CommitTime:    ctx.BlockTime(),
		CommitFee:     auction.CommitFee,
		RevealFee:     auction.RevealFee,
	}

	if err = k.SaveBid(ctx, &bid); err != nil {
		return nil, err
	}

	return &bid, nil
}

func (k Keeper) RevealBid(ctx sdk.Context, msg auctiontypes.MsgRevealBid) (*auctiontypes.Auction, error) {
	if has, err := k.HasAuction(ctx, msg.AuctionId); !has {
		if err != nil {
			return nil, err
		}
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Auction not found.")
	}

	auction, err := k.GetAuctionById(ctx, msg.AuctionId)
	if err != nil {
		return nil, err
	}

	if auction.Status != auctiontypes.AuctionStatusRevealPhase {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Auction is not in reveal phase.")
	}

	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	bidder := signerAddress.String()
	bidExists, err := k.HasBid(ctx, msg.AuctionId, bidder)
	if err != nil {
		return nil, err
	}

	if !bidExists {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Bid not found.")
	}

	bid, err := k.GetBid(ctx, msg.AuctionId, bidder)
	if err != nil {
		return nil, err
	}

	if bid.Status != auctiontypes.BidStatusCommitted {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Bid not in committed state.")
	}

	revealBytes, err := hex.DecodeString(msg.Reveal)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Invalid reveal string.")
	}

	cid, err := wnsUtils.CIDFromJSONBytes(revealBytes)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Invalid reveal JSON.")
	}

	if bid.CommitHash != cid {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Commit hash mismatch.")
	}

	var reveal map[string]interface{}
	err = json.Unmarshal(revealBytes, &reveal)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Reveal JSON unmarshal error.")
	}

	chainId, err := wnsUtils.GetAttributeAsString(reveal, "chainId")
	if err != nil || chainId != ctx.ChainID() {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Invalid reveal chainID.")
	}

	auctionId, err := wnsUtils.GetAttributeAsString(reveal, "auctionId")
	if err != nil || auctionId != msg.AuctionId {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Invalid reveal auction Id.")
	}

	bidderAddress, err := wnsUtils.GetAttributeAsString(reveal, "bidderAddress")
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Invalid reveal bid address.")
	}

	if bidderAddress != signerAddress.String() {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Reveal bid address mismatch.")
	}

	bidAmountStr, err := wnsUtils.GetAttributeAsString(reveal, "bidAmount")
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Invalid reveal bid amount.")
	}

	bidAmount, err := sdk.ParseCoinNormalized(bidAmountStr)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Invalid reveal bid amount.")
	}

	if bidAmount.IsLT(auction.MinimumBid) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Bid is lower than minimum bid.")
	}

	// Lock bid amount.
	sdkErr := k.bankKeeper.SendCoinsFromAccountToModule(ctx, signerAddress, auctiontypes.ModuleName, sdk.NewCoins(bidAmount))
	if sdkErr != nil {
		return nil, sdkErr
	}

	// Update bid.
	bid.BidAmount = bidAmount
	bid.RevealTime = ctx.BlockTime()
	bid.Status = auctiontypes.BidStatusRevealed
	if err = k.SaveBid(ctx, &bid); err != nil {
		return nil, err
	}

	return &auction, nil
}

// GetParams gets the auction module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (*auctiontypes.Params, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &params, nil
}

// GetAuctionModuleBalances gets the auction module account(s) balances.
func (k Keeper) GetAuctionModuleBalances(ctx sdk.Context) sdk.Coins {
	moduleAddress := k.accountKeeper.GetModuleAddress(auctiontypes.ModuleName)
	balances := k.bankKeeper.GetAllBalances(ctx, moduleAddress)

	return balances
}

func (k Keeper) EndBlockerProcessAuctions(ctx sdk.Context) error {
	// Transition auction state (commit, reveal, expired, completed).
	if err := k.processAuctionPhases(ctx); err != nil {
		return err
	}

	// Delete stale auctions.
	return k.deleteCompletedAuctions(ctx)
}

func (k Keeper) processAuctionPhases(ctx sdk.Context) error {
	auctions, err := k.MatchAuctions(ctx, func(_ *auctiontypes.Auction) (bool, error) {
		return true, nil
	})
	if err != nil {
		return err
	}

	for _, auction := range auctions {
		// Commit -> Reveal state.
		if auction.Status == auctiontypes.AuctionStatusCommitPhase && ctx.BlockTime().After(auction.CommitsEndTime) {
			auction.Status = auctiontypes.AuctionStatusRevealPhase
			if err = k.SaveAuction(ctx, auction); err != nil {
				return err
			}

			k.Logger(ctx).Info(fmt.Sprintf("Moved auction %s to reveal phase.", auction.Id))
		}

		// Reveal -> Expired state.
		if auction.Status == auctiontypes.AuctionStatusRevealPhase && ctx.BlockTime().After(auction.RevealsEndTime) {
			auction.Status = auctiontypes.AuctionStatusExpired
			if err = k.SaveAuction(ctx, auction); err != nil {
				return err
			}

			k.Logger(ctx).Info(fmt.Sprintf("Moved auction %s to expired state.", auction.Id))
		}

		// If auction has expired, pick a winner from revealed bids.
		if auction.Status == auctiontypes.AuctionStatusExpired {
			if err = k.pickAuctionWinner(ctx, auction); err != nil {
				return err
			}
		}
	}

	return nil
}

// Delete completed stale auctions.
func (k Keeper) deleteCompletedAuctions(ctx sdk.Context) error {
	auctions, err := k.MatchAuctions(ctx, func(auction *auctiontypes.Auction) (bool, error) {
		deleteTime := auction.RevealsEndTime.Add(CompletedAuctionDeleteTimeout)
		return auction.Status == auctiontypes.AuctionStatusCompleted && ctx.BlockTime().After(deleteTime), nil
	})
	if err != nil {
		return err
	}

	for _, auction := range auctions {
		k.Logger(ctx).Info(fmt.Sprintf("Deleting completed auction %s after timeout.", auction.Id))
		if err := k.DeleteAuction(ctx, *auction); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) pickAuctionWinner(ctx sdk.Context, auction *auctiontypes.Auction) error {
	k.Logger(ctx).Info(fmt.Sprintf("Picking auction %s winner.", auction.Id))

	var highestBid *auctiontypes.Bid
	var secondHighestBid *auctiontypes.Bid

	bids, err := k.GetBids(ctx, auction.Id)
	if err != nil {
		return err
	}

	for _, bid := range bids {
		k.Logger(ctx).Info(fmt.Sprintf("Processing bid %s %s", bid.BidderAddress, bid.BidAmount.String()))

		// Only consider revealed bids.
		if bid.Status != auctiontypes.BidStatusRevealed {
			k.Logger(ctx).Info(fmt.Sprintf("Ignoring unrevealed bid %s %s", bid.BidderAddress, bid.BidAmount.String()))
			continue
		}

		// Init highest bid.
		if highestBid == nil {
			highestBid = bid
			k.Logger(ctx).Info(fmt.Sprintf("Initializing 1st bid %s %s", bid.BidderAddress, bid.BidAmount.String()))
			continue
		}

		//nolint: all
		if highestBid.BidAmount.IsLT(bid.BidAmount) {
			k.Logger(ctx).Info(fmt.Sprintf("New highest bid %s %s", bid.BidderAddress, bid.BidAmount.String()))

			secondHighestBid = highestBid
			highestBid = bid

			k.Logger(ctx).Info(fmt.Sprintf("Updated 1st bid %s %s", highestBid.BidderAddress, highestBid.BidAmount.String()))
			k.Logger(ctx).Info(fmt.Sprintf("Updated 2nd bid %s %s", secondHighestBid.BidderAddress, secondHighestBid.BidAmount.String()))

		} else if secondHighestBid == nil || secondHighestBid.BidAmount.IsLT(bid.BidAmount) {
			k.Logger(ctx).Info(fmt.Sprintf("New 2nd highest bid %s %s", bid.BidderAddress, bid.BidAmount.String()))

			secondHighestBid = bid
			k.Logger(ctx).Info(fmt.Sprintf("Updated 2nd bid %s %s", secondHighestBid.BidderAddress, secondHighestBid.BidAmount.String()))
		} else {
			k.Logger(ctx).Info(fmt.Sprintf("Ignoring bid as it doesn't affect 1st/2nd price %s %s", bid.BidderAddress, bid.BidAmount.String()))
		}
	}

	// Highest bid is the winner, but pays second highest bid price.
	auction.Status = auctiontypes.AuctionStatusCompleted

	if highestBid != nil {
		auction.WinnerAddress = highestBid.BidderAddress
		auction.WinningBid = highestBid.BidAmount

		// Winner pays 2nd price, if a 2nd price exists.
		auction.WinningPrice = highestBid.BidAmount
		if secondHighestBid != nil {
			auction.WinningPrice = secondHighestBid.BidAmount
		}
		k.Logger(ctx).Info(fmt.Sprintf("Auction %s winner %s.", auction.Id, auction.WinnerAddress))
		k.Logger(ctx).Info(fmt.Sprintf("Auction %s winner bid %s.", auction.Id, auction.WinningBid.String()))
		k.Logger(ctx).Info(fmt.Sprintf("Auction %s winner price %s.", auction.Id, auction.WinningPrice.String()))
	} else {
		k.Logger(ctx).Info(fmt.Sprintf("Auction %s has no valid revealed bids (no winner).", auction.Id))
	}

	if err := k.SaveAuction(ctx, auction); err != nil {
		return err
	}

	for _, bid := range bids {
		bidderAddress, err := sdk.AccAddressFromBech32(bid.BidderAddress)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Invalid bidderAddress address. %v", err))
			panic("Invalid bidder address.")
		}

		if bid.Status == auctiontypes.BidStatusRevealed {
			// Send reveal fee back to bidders that've revealed the bid.
			sdkErr := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, bidderAddress, sdk.NewCoins(bid.RevealFee))
			if sdkErr != nil {
				k.Logger(ctx).Error(fmt.Sprintf("Auction error returning reveal fee: %v", sdkErr))
				panic(sdkErr)
			}
		}

		// Send back locked bid amount to all bidders.
		sdkErr := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, bidderAddress, sdk.NewCoins(bid.BidAmount))
		if sdkErr != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Auction error returning bid amount: %v", sdkErr))
			panic(sdkErr)
		}
	}

	// Process winner account (if nobody bids, there won't be a winner).
	if auction.WinnerAddress != "" {
		winnerAddress, err := sdk.AccAddressFromBech32(auction.WinnerAddress)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Invalid winner address. %v", err))
			panic("Invalid winner address.")
		}

		// Take 2nd price from winner.
		sdkErr := k.bankKeeper.SendCoinsFromAccountToModule(ctx, winnerAddress, auctiontypes.ModuleName, sdk.NewCoins(auction.WinningPrice))
		if sdkErr != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Auction error taking funds from winner: %v", sdkErr))
			panic(sdkErr)
		}

		// Burn anything over the min. bid amount.
		amountToBurn := auction.WinningPrice.Sub(auction.MinimumBid)
		if amountToBurn.IsNegative() {
			k.Logger(ctx).Error("Auction coins to burn cannot be negative.")
			panic("Auction coins to burn cannot be negative.")
		}

		// Use auction burn module account instead of actually burning coins to better keep track of supply.
		sdkErr = k.bankKeeper.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, auctiontypes.AuctionBurnModuleAccountName, sdk.NewCoins(amountToBurn))
		if sdkErr != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Auction error burning coins: %v", sdkErr))
			panic(sdkErr)
		}
	}

	// Notify other modules (hook).
	k.Logger(ctx).Info(fmt.Sprintf("Auction %s notifying %d modules.", auction.Id, len(k.usageKeepers)))
	for _, keeper := range k.usageKeepers {
		k.Logger(ctx).Info(fmt.Sprintf("Auction %s notifying module %s.", auction.Id, keeper.ModuleName()))
		keeper.OnAuctionWinnerSelected(ctx, auction.Id)
	}

	return nil
}
