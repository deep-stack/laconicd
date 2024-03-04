package registry

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/stretchr/testify/suite"

	"git.vdb.to/cerc-io/laconic2d/tests/e2e"
)

func TestRegistryE2ETestSuite(t *testing.T) {
	cfg := network.DefaultConfig(e2e.NewTestNetworkFixture)
	cfg.NumValidators = 1

	suite.Run(t, NewE2ETestSuite(cfg))
}
