package env_test

import (
	"github.com/anhdht/go-common/pkg/env"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type envSuite struct {
	suite.Suite

	key string
}

func TestEnvSuite(t *testing.T) {
	suite.Run(t, &envSuite{
		key: "TEST",
	})
}

func (suite *envSuite) TearDownTest() {
	suite.Nil(os.Unsetenv(suite.key))
}

func (suite *envSuite) TestGetEnv_Exists() {
	value := "test"

	err := os.Setenv(suite.key, value)
	suite.Nil(err, "os.Setenv should returns nil")

	v := env.GetEnv(suite.key, "default")
	suite.Equal(value, v, "env.GetEnv should returns `test`")
}

func (suite *envSuite) TestGetEnv_NonExists() {
	value := "default"
	v := env.GetEnv(suite.key, value)
	suite.Equal(value, v, "env.GetEnv should returns `default`")
}
