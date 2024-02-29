package keeper_test

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"

	integrationTest "git.vdb.to/cerc-io/laconic2d/tests/integration"
	types "git.vdb.to/cerc-io/laconic2d/x/auction"
)

const testCommitHash = "71D8CF34026E32A3A34C2C2D4ADF25ABC8D7943A4619761BE27F196603D91B9D"

func (kts *KeeperTestSuite) TestGrpcQueryParams() {
	testCases := []struct {
		msg string
		req *types.QueryParamsRequest
	}{
		{
			"fetch params",
			&types.QueryParamsRequest{},
		},
	}
	for _, test := range testCases {
		kts.Run(fmt.Sprintf("Case %s", test.msg), func() {
			resp, err := kts.queryClient.Params(context.Background(), test.req)
			kts.Require().Nil(err)
			kts.Require().Equal(*(resp.Params), types.DefaultParams())
		})
	}
}

func (kts *KeeperTestSuite) TestGrpcGetAuction() {
	testCases := []struct {
		msg           string
		req           *types.QueryAuctionRequest
		createAuction bool
	}{
		{
			"fetch auction with empty auction ID",
			&types.QueryAuctionRequest{},
			false,
		},
		{
			"fetch auction with valid auction ID",
			&types.QueryAuctionRequest{},
			true,
		},
	}

	for _, test := range testCases {
		kts.Run(fmt.Sprintf("Case %s", test.msg), func() {
			var expectedAuction types.Auction
			if test.createAuction {
				auction, _, err := kts.createAuctionAndCommitBid(false)
				kts.Require().Nil(err)
				test.req.Id = auction.Id
				expectedAuction = *auction
			}

			resp, err := kts.queryClient.GetAuction(context.Background(), test.req)
			if test.createAuction {
				kts.Require().Nil(err)
				kts.Require().NotNil(resp.GetAuction())
				kts.Require().EqualExportedValues(expectedAuction, *(resp.GetAuction()))
			} else {
				kts.Require().NotNil(err)
				kts.Require().Error(err)
			}
		})
	}
}

func (kts *KeeperTestSuite) TestGrpcGetAllAuctions() {
	testCases := []struct {
		msg            string
		req            *types.QueryAuctionsRequest
		createAuctions bool
		auctionCount   int
	}{
		{
			"fetch auctions when no auctions exist",
			&types.QueryAuctionsRequest{},
			false,
			0,
		},

		{
			"fetch auctions with one auction created",
			&types.QueryAuctionsRequest{},
			true,
			1,
		},
	}

	for _, test := range testCases {
		kts.Run(fmt.Sprintf("Case %s", test.msg), func() {
			if test.createAuctions {
				_, _, err := kts.createAuctionAndCommitBid(false)
				kts.Require().Nil(err)
			}

			resp, _ := kts.queryClient.Auctions(context.Background(), test.req)
			kts.Require().Equal(test.auctionCount, len(resp.GetAuctions().Auctions))
		})
	}
}

func (kts *KeeperTestSuite) TestGrpcGetBids() {
	testCases := []struct {
		msg           string
		req           *types.QueryBidsRequest
		createAuction bool
		commitBid     bool
		bidCount      int
	}{
		{
			"fetch all bids when no auction exists",
			&types.QueryBidsRequest{},
			false,
			false,
			0,
		},
		{
			"fetch all bids for valid auction but no added bids",
			&types.QueryBidsRequest{},
			true,
			false,
			0,
		},
		{
			"fetch all bids for valid auction and valid bid",
			&types.QueryBidsRequest{},
			true,
			true,
			1,
		},
	}

	for _, test := range testCases {
		kts.Run(fmt.Sprintf("Case %s", test.msg), func() {
			if test.createAuction {
				auction, _, err := kts.createAuctionAndCommitBid(test.commitBid)
				kts.Require().NoError(err)
				test.req.AuctionId = auction.Id
			}

			resp, err := kts.queryClient.GetBids(context.Background(), test.req)
			if test.createAuction {
				kts.Require().Nil(err)
				kts.Require().Equal(test.bidCount, len(resp.GetBids()))
			} else {
				kts.Require().NotNil(err)
				kts.Require().Error(err)
			}
		})
	}
}

func (kts *KeeperTestSuite) TestGrpcGetBid() {
	testCases := []struct {
		msg                 string
		req                 *types.QueryBidRequest
		createAuctionAndBid bool
	}{
		{
			"fetch bid when bid does not exist",
			&types.QueryBidRequest{},
			false,
		},
		{
			"fetch bid when valid bid exists",
			&types.QueryBidRequest{},
			true,
		},
	}

	for _, test := range testCases {
		kts.Run(fmt.Sprintf("Case %s", test.msg), func() {
			if test.createAuctionAndBid {
				auction, bid, err := kts.createAuctionAndCommitBid(test.createAuctionAndBid)
				kts.Require().NoError(err)
				test.req.AuctionId = auction.Id
				test.req.Bidder = bid.BidderAddress
			}

			resp, err := kts.queryClient.GetBid(context.Background(), test.req)
			if test.createAuctionAndBid {
				kts.Require().NoError(err)
				kts.Require().NotNil(resp.Bid)
				kts.Require().Equal(test.req.Bidder, resp.Bid.BidderAddress)
			} else {
				kts.Require().NotNil(err)
				kts.Require().Error(err)
			}
		})
	}
}

func (kts *KeeperTestSuite) TestGrpcGetAuctionsByBidder() {
	testCases := []struct {
		msg                       string
		req                       *types.QueryAuctionsByBidderRequest
		createAuctionAndCommitBid bool
		auctionCount              int
	}{
		{
			"get auctions by bidder with invalid bidder address",
			&types.QueryAuctionsByBidderRequest{},
			false,
			0,
		},
		{
			"get auctions by bidder with valid auction and bid",
			&types.QueryAuctionsByBidderRequest{},
			true,
			1,
		},
	}

	for _, test := range testCases {
		kts.Run(fmt.Sprintf("Case %s", test.msg), func() {
			if test.createAuctionAndCommitBid {
				_, bid, err := kts.createAuctionAndCommitBid(test.createAuctionAndCommitBid)
				kts.Require().NoError(err)
				test.req.BidderAddress = bid.BidderAddress
			}

			resp, err := kts.queryClient.AuctionsByBidder(context.Background(), test.req)
			if test.createAuctionAndCommitBid {
				kts.Require().NoError(err)
				kts.Require().NotNil(resp.Auctions)
				kts.Require().Equal(test.auctionCount, len(resp.Auctions.Auctions))
			} else {
				kts.Require().NotNil(err)
				kts.Require().Error(err)
			}
		})
	}
}

func (kts *KeeperTestSuite) TestGrpcGetAuctionsByOwner() {
	testCases := []struct {
		msg           string
		req           *types.QueryAuctionsByOwnerRequest
		createAuction bool
		auctionCount  int
	}{
		{
			"get auctions by owner with invalid owner address",
			&types.QueryAuctionsByOwnerRequest{},
			false,
			0,
		},
		{
			"get auctions by owner with valid auction",
			&types.QueryAuctionsByOwnerRequest{},
			true,
			1,
		},
	}

	for _, test := range testCases {
		kts.Run(fmt.Sprintf("Case %s", test.msg), func() {
			if test.createAuction {
				auction, _, err := kts.createAuctionAndCommitBid(false)
				kts.Require().NoError(err)
				test.req.OwnerAddress = auction.OwnerAddress
			}

			resp, err := kts.queryClient.AuctionsByOwner(context.Background(), test.req)
			if test.createAuction {
				kts.Require().NoError(err)
				kts.Require().NotNil(resp.Auctions)
				kts.Require().Equal(test.auctionCount, len(resp.Auctions.Auctions))
			} else {
				kts.Require().NotNil(err)
				kts.Require().Error(err)
			}
		})
	}
}

func (kts *KeeperTestSuite) TestGrpcQueryBalance() {
	testCases := []struct {
		msg           string
		req           *types.QueryGetAuctionModuleBalanceRequest
		createAuction bool
		auctionCount  int
	}{
		{
			"get balance with no auctions created",
			&types.QueryGetAuctionModuleBalanceRequest{},
			false,
			0,
		},
		{
			"get balance with single auction created",
			&types.QueryGetAuctionModuleBalanceRequest{},
			true,
			1,
		},
	}

	for _, test := range testCases {
		if test.createAuction {
			_, _, err := kts.createAuctionAndCommitBid(true)
			kts.Require().NoError(err)
		}

		resp, err := kts.queryClient.GetAuctionModuleBalance(context.Background(), test.req)
		kts.Require().NoError(err)
		kts.Require().Equal(test.auctionCount, len(resp.GetBalance()))
	}
}

func (kts *KeeperTestSuite) createAuctionAndCommitBid(commitBid bool) (*types.Auction, *types.Bid, error) {
	ctx, k := kts.SdkCtx, kts.AuctionKeeper
	accCount := 1
	if commitBid {
		accCount++
	}

	// Create funded account(s)
	accounts := simtestutil.AddTestAddrs(kts.BankKeeper, integrationTest.BondDenomProvider{}, ctx, accCount, math.NewInt(100))

	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, nil, err
	}

	auction, err := k.CreateAuction(ctx, types.NewMsgCreateAuction(*params, accounts[0]))
	if err != nil {
		return nil, nil, err
	}

	if commitBid {
		bid, err := k.CommitBid(ctx, types.NewMsgCommitBid(auction.Id, testCommitHash, accounts[1]))
		if err != nil {
			return nil, nil, err
		}

		return auction, bid, nil
	}

	return auction, nil, nil
}
