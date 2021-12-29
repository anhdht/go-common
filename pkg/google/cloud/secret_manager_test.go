package cloud_test

import (
	"github.com/anhdht/go-common/pkg/google/cloud"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type secretManagerSuite struct {
	suite.Suite

	projectId string
}

func TestSecretManagerSuite(t *testing.T) {
	suite.Run(t, &secretManagerSuite{
		projectId: "PROJECT_ID",
	})
}

func (suite *secretManagerSuite) TearDownTest() {
	suite.Nil(os.Unsetenv(suite.projectId))
}

func (suite *secretManagerSuite) TestSecretManager_Load_UnsetProjectId() {
	suite.Nil(os.Unsetenv(suite.projectId))

	s := cloud.NewSecretManager()
	data, err := s.Load()
	
	suite.Nil(err)
	suite.Nil(data)
}
