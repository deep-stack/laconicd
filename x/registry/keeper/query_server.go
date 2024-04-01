package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	registrytypes "git.vdb.to/cerc-io/laconicd/x/registry"
)

var _ registrytypes.QueryServer = queryServer{}

type queryServer struct {
	k Keeper
}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k Keeper) registrytypes.QueryServer {
	return queryServer{k}
}

func (qs queryServer) Params(c context.Context, _ *registrytypes.QueryParamsRequest) (*registrytypes.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	params, err := qs.k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	return &registrytypes.QueryParamsResponse{Params: params}, nil
}

func (qs queryServer) Records(c context.Context, req *registrytypes.QueryRecordsRequest) (*registrytypes.QueryRecordsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	attributes := req.GetAttributes()
	all := req.GetAll()

	var records []registrytypes.Record
	var err error
	if len(attributes) > 0 {
		records, err = qs.k.RecordsFromAttributes(ctx, attributes, all)
		if err != nil {
			return nil, err
		}
	} else {
		records, err = qs.k.ListRecords(ctx)
		if err != nil {
			return nil, err
		}
	}

	return &registrytypes.QueryRecordsResponse{Records: records}, nil
}

func (qs queryServer) GetRecord(c context.Context, req *registrytypes.QueryGetRecordRequest) (*registrytypes.QueryGetRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	id := req.GetId()

	if has, err := qs.k.HasRecord(ctx, req.Id); !has {
		if err != nil {
			return nil, err
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "Record not found.")
	}

	record, err := qs.k.GetRecordById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &registrytypes.QueryGetRecordResponse{Record: record}, nil
}

func (qs queryServer) GetRecordsByBondId(c context.Context, req *registrytypes.QueryGetRecordsByBondIdRequest) (*registrytypes.QueryGetRecordsByBondIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	records, err := qs.k.GetRecordsByBondId(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &registrytypes.QueryGetRecordsByBondIdResponse{Records: records}, nil
}

func (qs queryServer) GetRegistryModuleBalance(c context.Context,
	_ *registrytypes.QueryGetRegistryModuleBalanceRequest,
) (*registrytypes.QueryGetRegistryModuleBalanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	balances := qs.k.GetModuleBalances(ctx)

	return &registrytypes.QueryGetRegistryModuleBalanceResponse{
		Balances: balances,
	}, nil
}

func (qs queryServer) NameRecords(c context.Context, _ *registrytypes.QueryNameRecordsRequest) (*registrytypes.QueryNameRecordsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	nameRecords, err := qs.k.ListNameRecords(ctx)
	if err != nil {
		return nil, err
	}

	return &registrytypes.QueryNameRecordsResponse{Names: nameRecords}, nil
}

func (qs queryServer) Whois(c context.Context, request *registrytypes.QueryWhoisRequest) (*registrytypes.QueryWhoisResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	nameAuthority, err := qs.k.GetNameAuthority(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	return &registrytypes.QueryWhoisResponse{NameAuthority: nameAuthority}, nil
}

func (qs queryServer) LookupLrn(c context.Context, req *registrytypes.QueryLookupLrnRequest) (*registrytypes.QueryLookupLrnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	lrn := req.GetLrn()

	lrnExists, err := qs.k.HasNameRecord(ctx, lrn)
	if err != nil {
		return nil, err
	}
	if !lrnExists {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "LRN not found.")
	}

	nameRecord, err := qs.k.LookupNameRecord(ctx, lrn)
	if nameRecord == nil {
		if err != nil {
			return nil, err
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "name record not found.")
	}

	return &registrytypes.QueryLookupLrnResponse{Name: nameRecord}, nil
}

func (qs queryServer) ResolveLrn(c context.Context, req *registrytypes.QueryResolveLrnRequest) (*registrytypes.QueryResolveLrnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	lrn := req.GetLrn()
	record, err := qs.k.ResolveLRN(ctx, lrn)
	if record == nil {
		if err != nil {
			return nil, err
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "record not found.")
	}

	return &registrytypes.QueryResolveLrnResponse{Record: record}, nil
}
