package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	bondtypes "git.vdb.to/cerc-io/laconic2d/x/bond"
)

type Keeper struct {
	// Codecs
	cdc codec.BinaryCodec

	// External keepers
	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper

	// Track bond usage in other cosmos-sdk modules (more like a usage tracker).
	// usageKeepers []types.BondUsageKeeper

	// paramSubspace paramtypes.Subspace

	// State management
	Schema collections.Schema
	Bonds  collections.Map[string, bondtypes.Bond]
}

// NewKeeper creates new instances of the bond Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	accountKeeper auth.AccountKeeper,
	bankKeeper bank.Keeper,
	// usageKeepers []types.BondUsageKeeper,
	// ps paramtypes.Subspace,
) Keeper {
	// set KeyTable if it has not already been set
	// if !ps.HasKeyTable() {
	// 	ps = ps.WithKeyTable(types.ParamKeyTable())
	// }

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:           cdc,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		Bonds:         collections.NewMap(sb, bondtypes.BondsKey, "bonds", collections.StringKey, codec.CollValue[bondtypes.Bond](cdc)),
		// usageKeepers:  usageKeepers,
		// paramSubspace: ps,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

// BondID simplifies generation of bond IDs.
type BondID struct {
	Address  sdk.Address
	AccNum   uint64
	Sequence uint64
}

// Generate creates the bond ID.
func (bondID BondID) Generate() string {
	hasher := sha256.New()
	str := fmt.Sprintf("%s:%d:%d", bondID.Address.String(), bondID.AccNum, bondID.Sequence)
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

// CreateBond creates a new bond.
func (k Keeper) CreateBond(ctx sdk.Context, ownerAddress sdk.AccAddress, coins sdk.Coins) (*bondtypes.Bond, error) {
	// Check if account has funds.
	for _, coin := range coins {
		balance := k.bankKeeper.HasBalance(ctx, ownerAddress, coin)
		if !balance {
			return nil, errorsmod.Wrap(sdkerrors.ErrInsufficientFunds, "failed to create bond; Insufficient funds")
		}
	}

	// Generate bond ID.
	account := k.accountKeeper.GetAccount(ctx, ownerAddress)
	bondID := BondID{
		Address:  ownerAddress,
		AccNum:   account.GetAccountNumber(),
		Sequence: account.GetSequence(),
	}.Generate()

	maxBondAmount := k.getMaxBondAmount(ctx)

	bond := bondtypes.Bond{Id: bondID, Owner: ownerAddress.String(), Balance: coins}
	if bond.Balance.IsAnyGT(maxBondAmount) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Max bond amount exceeded.")
	}

	// Move funds into the bond account module.
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, ownerAddress, bondtypes.ModuleName, bond.Balance)
	if err != nil {
		return nil, err
	}

	// Save bond in store.
	if err := k.Bonds.Set(ctx, bond.Id, bond); err != nil {
		return nil, err
	}

	return &bond, nil
}

// ListBonds - get all bonds.
func (k Keeper) ListBonds(ctx sdk.Context) ([]*bondtypes.Bond, error) {
	var bonds []*bondtypes.Bond

	iter, err := k.Bonds.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	for ; iter.Valid(); iter.Next() {
		bond, err := iter.Value()
		if err != nil {
			return nil, err
		}

		bonds = append(bonds, &bond)
	}

	return bonds, nil
}

func (k Keeper) getMaxBondAmount(ctx sdk.Context) sdk.Coins {
	params := k.GetParams(ctx)
	maxBondAmount := params.MaxBondAmount
	return sdk.NewCoins(maxBondAmount)
}
