package keeper

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"
	storetypes "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/gibson042/canonicaljson-go"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/node/basicnode"

	auctionkeeper "git.vdb.to/cerc-io/laconic2d/x/auction/keeper"
	bondkeeper "git.vdb.to/cerc-io/laconic2d/x/bond/keeper"
	registrytypes "git.vdb.to/cerc-io/laconic2d/x/registry"
	"git.vdb.to/cerc-io/laconic2d/x/registry/helpers"
)

// TODO: Add required methods

type RecordsIndexes struct {
	BondId *indexes.Multi[string, string, registrytypes.Record]
}

func (b RecordsIndexes) IndexesList() []collections.Index[string, registrytypes.Record] {
	return []collections.Index[string, registrytypes.Record]{b.BondId}
}

func newRecordIndexes(sb *collections.SchemaBuilder) RecordsIndexes {
	return RecordsIndexes{
		BondId: indexes.NewMulti(
			sb, registrytypes.RecordsByBondIdIndexPrefix, "records_by_bond_id",
			collections.StringKey, collections.StringKey,
			func(_ string, v registrytypes.Record) (string, error) {
				return v.BondId, nil
			},
		),
	}
}

// TODO
type AuthoritiesIndexes struct {
}

func (b AuthoritiesIndexes) IndexesList() []collections.Index[string, registrytypes.NameAuthority] {
	return []collections.Index[string, registrytypes.NameAuthority]{}
}

func newAuthorityIndexes(sb *collections.SchemaBuilder) AuthoritiesIndexes {
	return AuthoritiesIndexes{}
}

type NameRecordsIndexes struct {
	Cid *indexes.Multi[string, string, registrytypes.NameRecord]
}

func (b NameRecordsIndexes) IndexesList() []collections.Index[string, registrytypes.NameRecord] {
	return []collections.Index[string, registrytypes.NameRecord]{b.Cid}
}

func newNameRecordIndexes(sb *collections.SchemaBuilder) NameRecordsIndexes {
	return NameRecordsIndexes{
		Cid: indexes.NewMulti(
			sb, registrytypes.NameRecordsByCidIndexPrefix, "name_records_by_cid",
			collections.StringKey, collections.StringKey,
			func(_ string, v registrytypes.NameRecord) (string, error) {
				return v.Latest.Id, nil
			},
		),
	}
}

type Keeper struct {
	cdc codec.BinaryCodec

	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper
	recordKeeper  RecordKeeper
	bondKeeper    bondkeeper.Keeper
	auctionKeeper auctionkeeper.Keeper

	// state management
	Schema      collections.Schema
	Params      collections.Item[registrytypes.Params]
	Records     *collections.IndexedMap[string, registrytypes.Record, RecordsIndexes]
	Authorities *collections.IndexedMap[string, registrytypes.NameAuthority, AuthoritiesIndexes]
	NameRecords *collections.IndexedMap[string, registrytypes.NameRecord, NameRecordsIndexes]
}

// NewKeeper creates a new Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	accountKeeper auth.AccountKeeper,
	bankKeeper bank.Keeper,
	recordKeeper RecordKeeper,
	bondKeeper bondkeeper.Keeper,
	auctionKeeper auctionkeeper.Keeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:           cdc,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		recordKeeper:  recordKeeper,
		bondKeeper:    bondKeeper,
		auctionKeeper: auctionKeeper,
		Params:        collections.NewItem(sb, registrytypes.ParamsPrefix, "params", codec.CollValue[registrytypes.Params](cdc)),
		Records: collections.NewIndexedMap(
			sb, registrytypes.RecordsPrefix, "records",
			collections.StringKey, codec.CollValue[registrytypes.Record](cdc),
			newRecordIndexes(sb),
		),
		Authorities: collections.NewIndexedMap(
			sb, registrytypes.AuthoritiesPrefix, "authorities",
			collections.StringKey, codec.CollValue[registrytypes.NameAuthority](cdc),
			newAuthorityIndexes(sb),
		),
		NameRecords: collections.NewIndexedMap(
			sb, registrytypes.NameRecordsPrefix, "name_records",
			collections.StringKey, codec.CollValue[registrytypes.NameRecord](cdc),
			newNameRecordIndexes(sb),
		),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

// HasRecord - checks if a record by the given id exists.
func (k Keeper) HasRecord(ctx sdk.Context, id string) (bool, error) {
	has, err := k.Records.Has(ctx, id)
	if err != nil {
		return false, err
	}

	return has, nil
}

// ListRecords - get all records.
func (k Keeper) ListRecords(ctx sdk.Context) ([]registrytypes.Record, error) {
	var records []registrytypes.Record

	err := k.Records.Walk(ctx, nil, func(key string, value registrytypes.Record) (bool, error) {
		if err := k.populateRecordNames(ctx, &value); err != nil {
			return true, err
		}
		records = append(records, value)

		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return records, nil
}

// GetRecordById - gets a record from the store.
func (k Keeper) GetRecordById(ctx sdk.Context, id string) (registrytypes.Record, error) {
	record, err := k.Records.Get(ctx, id)
	if err != nil {
		return registrytypes.Record{}, err
	}

	if err := k.populateRecordNames(ctx, &record); err != nil {
		return registrytypes.Record{}, err
	}

	return record, nil
}

// GetRecordsByBondId - gets a record from the store.
func (k Keeper) GetRecordsByBondId(ctx sdk.Context, bondId string) ([]registrytypes.Record, error) {
	var records []registrytypes.Record

	err := k.Records.Indexes.BondId.Walk(ctx, collections.NewPrefixedPairRange[string, string](bondId), func(bondId string, id string) (bool, error) {
		record, err := k.Records.Get(ctx, id)
		if err != nil {
			return true, err
		}

		if err := k.populateRecordNames(ctx, &record); err != nil {
			return true, err
		}
		records = append(records, record)

		return false, nil
	})
	if err != nil {
		return []registrytypes.Record{}, err
	}

	return records, nil
}

// RecordsFromAttributes gets a list of records whose attributes match all provided values
func (k Keeper) RecordsFromAttributes(ctx sdk.Context, attributes []*registrytypes.QueryRecordsRequest_KeyValueInput, all bool) ([]registrytypes.Record, error) {
	panic("unimplemented")
}

func (k Keeper) GetRecordExpiryQueue(ctx sdk.Context) []*registrytypes.ExpiryQueueRecord {
	panic("unimplemented")
}

// PutRecord - saves a record to the store.
func (k Keeper) SaveRecord(ctx sdk.Context, record registrytypes.Record) error {
	return k.Records.Set(ctx, record.Id, record)

	// TODO
	// k.updateBlockChangeSetForRecord(ctx, record.Id)
}

// ProcessSetRecord creates a record.
func (k Keeper) SetRecord(ctx sdk.Context, msg registrytypes.MsgSetRecord) (*registrytypes.ReadableRecord, error) {
	payload := msg.Payload.ToReadablePayload()
	record := registrytypes.ReadableRecord{Attributes: payload.RecordAttributes, BondId: msg.BondId}

	// Check signatures.
	resourceSignBytes, _ := record.GetSignBytes()
	cid, err := record.GetCid()
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Invalid record JSON")
	}

	record.Id = cid

	has, err := k.HasRecord(ctx, record.Id)
	if err != nil {
		return nil, err
	}
	if has {
		// Immutable record already exists. No-op.
		return &record, nil
	}

	record.Owners = []string{}
	for _, sig := range payload.Signatures {
		pubKey, err := legacy.PubKeyFromBytes(helpers.BytesFromBase64(sig.PubKey))
		if err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, fmt.Sprint("Error decoding pubKey from bytes: ", err))
		}

		sigOK := pubKey.VerifySignature(resourceSignBytes, helpers.BytesFromBase64(sig.Sig))
		if !sigOK {
			return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, fmt.Sprint("Signature mismatch: ", sig.PubKey))
		}
		record.Owners = append(record.Owners, pubKey.Address().String())
	}

	// Sort owners list.
	sort.Strings(record.Owners)
	sdkErr := k.processRecord(ctx, &record, false)
	if sdkErr != nil {
		return nil, sdkErr
	}

	return &record, nil
}

func (k Keeper) processRecord(ctx sdk.Context, record *registrytypes.ReadableRecord, isRenewal bool) error {
	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	rent := params.RecordRent
	if err = k.bondKeeper.TransferCoinsToModuleAccount(
		ctx, record.BondId, registrytypes.RecordRentModuleAccountName, sdk.NewCoins(rent),
	); err != nil {
		return err
	}

	record.CreateTime = ctx.BlockHeader().Time.Format(time.RFC3339)
	record.ExpiryTime = ctx.BlockHeader().Time.Add(params.RecordRentDuration).Format(time.RFC3339)
	record.Deleted = false

	recordObj, err := record.ToRecordObj()
	if err != nil {
		return err
	}

	// Save record in store.
	if err = k.SaveRecord(ctx, recordObj); err != nil {
		return err
	}

	// TODO look up/validate record type here

	if err := k.processAttributes(ctx, record.Attributes, record.Id, ""); err != nil {
		return err
	}

	// TODO
	// expiryTimeKey := GetAttributesIndexKey(ExpiryTimeAttributeName, []byte(record.ExpiryTime))
	// if err := k.SetAttributeMapping(ctx, expiryTimeKey, record.ID); err != nil {
	// 	return err
	// }

	// k.InsertRecordExpiryQueue(ctx, recordObj)

	return nil
}

func (k Keeper) processAttributes(ctx sdk.Context, attrs registrytypes.AttributeMap, id string, prefix string) error {
	np := basicnode.Prototype.Map
	nb := np.NewBuilder()
	encAttrs, err := canonicaljson.Marshal(attrs)
	if err != nil {
		return err
	}
	if len(attrs) == 0 {
		encAttrs = []byte("{}")
	}
	err = dagjson.Decode(nb, bytes.NewReader(encAttrs))
	if err != nil {
		return fmt.Errorf("failed to decode attributes: %w", err)
	}
	n := nb.Build()
	if n.Kind() != ipld.Kind_Map {
		return fmt.Errorf("record attributes must be a map, not %T", n.Kind())
	}

	return k.processAttributeMap(ctx, n, id, prefix)
}

func (k Keeper) processAttributeMap(ctx sdk.Context, n ipld.Node, id string, prefix string) error {
	for it := n.MapIterator(); !it.Done(); {
		//nolint:misspell
		keynode, valuenode, err := it.Next()
		if err != nil {
			return err
		}
		key, err := keynode.AsString()
		if err != nil {
			return err
		}

		if valuenode.Kind() == ipld.Kind_Map {
			err := k.processAttributeMap(ctx, valuenode, id, key)
			if err != nil {
				return err
			}
		} else {
			var buf bytes.Buffer
			if err := dagjson.Encode(valuenode, &buf); err != nil {
				return err
			}
			// TODO
			// value := buf.Bytes()
			// indexKey := GetAttributesIndexKey(prefix+key, value)
			// if err := k.SetAttributeMapping(ctx, indexKey, id); err != nil {
			// 	return err
			// }
		}
	}
	return nil
}

func (k Keeper) populateRecordNames(ctx sdk.Context, record *registrytypes.Record) error {
	iter, err := k.NameRecords.Indexes.Cid.MatchExact(ctx, record.Id)
	if err != nil {
		return err
	}

	names, err := iter.PrimaryKeys()
	if err != nil {
		return err
	}
	record.Names = names

	return nil
}

// GetModuleBalances gets the registry module account(s) balances.
func (k Keeper) GetModuleBalances(ctx sdk.Context) []*registrytypes.AccountBalance {
	var balances []*registrytypes.AccountBalance
	accountNames := []string{
		registrytypes.RecordRentModuleAccountName,
		registrytypes.AuthorityRentModuleAccountName,
	}

	for _, accountName := range accountNames {
		moduleAddress := k.accountKeeper.GetModuleAddress(accountName)

		moduleAccount := k.accountKeeper.GetAccount(ctx, moduleAddress)
		if moduleAccount != nil {
			accountBalance := k.bankKeeper.GetAllBalances(ctx, moduleAddress)
			balances = append(balances, &registrytypes.AccountBalance{
				AccountName: accountName,
				Balance:     accountBalance,
			})
		}
	}

	return balances
}
