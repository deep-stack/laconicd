package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	registrytypes "git.vdb.to/cerc-io/laconic2d/x/registry"
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

func (qs queryServer) GetRecord(c context.Context, req *registrytypes.QueryRecordByIdRequest) (*registrytypes.QueryRecordByIdResponse, error) {
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

	return &registrytypes.QueryRecordByIdResponse{Record: record}, nil
}

func (qs queryServer) GetRecordsByBondId(c context.Context, req *registrytypes.QueryRecordsByBondIdRequest) (*registrytypes.QueryRecordsByBondIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	records, err := qs.k.GetRecordsByBondId(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &registrytypes.QueryRecordsByBondIdResponse{Records: records}, nil
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

func (qs queryServer) LookupCrn(c context.Context, req *registrytypes.QueryLookupCrnRequest) (*registrytypes.QueryLookupCrnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	crn := req.GetCrn()

	crnExists, err := qs.k.HasNameRecord(ctx, crn)
	if err != nil {
		return nil, err
	}
	if !crnExists {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "CRN not found.")
	}

	nameRecord, err := qs.k.LookupNameRecord(ctx, crn)
	if nameRecord == nil {
		if err != nil {
			return nil, err
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "name record not found.")
	}

	return &registrytypes.QueryLookupCrnResponse{Name: nameRecord}, nil
}

func (qs queryServer) ResolveCrn(c context.Context, req *registrytypes.QueryResolveCrnRequest) (*registrytypes.QueryResolveCrnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	crn := req.GetCrn()
	record, err := qs.k.ResolveCRN(ctx, crn)
	if record == nil {
		if err != nil {
			return nil, err
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "record not found.")
	}

	return &registrytypes.QueryResolveCrnResponse{Record: record}, nil
}

func (qs queryServer) GetRecordExpiryQueue(c context.Context, _ *registrytypes.QueryGetRecordExpiryQueueRequest) (*registrytypes.QueryGetRecordExpiryQueueResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	records := qs.k.GetRecordExpiryQueue(ctx)
	return &registrytypes.QueryGetRecordExpiryQueueResponse{Records: records}, nil
}

func (qs queryServer) GetAuthorityExpiryQueue(c context.Context,
	_ *registrytypes.QueryGetAuthorityExpiryQueueRequest,
) (*registrytypes.QueryGetAuthorityExpiryQueueResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	authorities := qs.k.GetAuthorityExpiryQueue(ctx)
	return &registrytypes.QueryGetAuthorityExpiryQueueResponse{Authorities: authorities}, nil
}
