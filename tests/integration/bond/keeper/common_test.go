package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	integrationTest "git.vdb.to/cerc-io/laconic2d/tests/integration"
	types "git.vdb.to/cerc-io/laconic2d/x/bond"
)

type KeeperTestSuite struct {
	suite.Suite
	integrationTest.TestFixture

	queryClient types.QueryClient
}

func (kts *KeeperTestSuite) SetupTest() {
	err := kts.TestFixture.Setup()
	assert.Nil(kts.T(), err)

	qr := kts.App.QueryHelper()
	kts.queryClient = types.NewQueryClient(qr)
}

func TestBondKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
