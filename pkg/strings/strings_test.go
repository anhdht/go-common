package strings_test

import (
	"github.com/anhdht/go-common/pkg/strings"
	"github.com/stretchr/testify/suite"
	"testing"
)

type stringsSuite struct {
	suite.Suite
}

func TestStringsSuite(t *testing.T) {
	suite.Run(t, &stringsSuite{})
}

func (suite *stringsSuite) TestTrimSlash() {
	suite.Equal("hello", strings.TrimSlash("hello/"))
	suite.Equal("hello", strings.TrimSlash("hello//"))
	suite.Equal("hello", strings.TrimSlash("/hello"))
	suite.Equal("hello", strings.TrimSlash("//hello"))
	suite.Equal("hello", strings.TrimSlash("/hello/"))
	suite.Equal("hello", strings.TrimSlash("\\hello"))
	suite.Equal("hello", strings.TrimSlash("\\\\hello"))
	suite.Equal("hello", strings.TrimSlash("hello\\"))
	suite.Equal("hello", strings.TrimSlash("hello\\\\"))
	suite.Equal("hello", strings.TrimSlash("hello\\/"))
	suite.Equal("hello", strings.TrimSlash("hello/\\"))
	suite.Equal("hello", strings.TrimSlash("hello/\\/"))
	suite.Equal("hello", strings.TrimSlash("\\/hello"))
	suite.Equal("hello", strings.TrimSlash("/\\hello"))
	suite.Equal("hello", strings.TrimSlash("/\\/hello"))
	suite.Equal("hello", strings.TrimSlash("/\\/hello/\\/"))
}
