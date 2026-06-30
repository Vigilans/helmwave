package helper_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/stretchr/testify/suite"
)

type DotenvTestSuite struct {
	suite.Suite
}

func TestDotenvTestSuite(t *testing.T) {
	suite.Run(t, new(DotenvTestSuite))
}

func (s *DotenvTestSuite) writeEnvFile(dir, name, key, val string) string {
	s.T().Helper()
	path := filepath.Join(dir, name)
	s.Require().NoError(os.WriteFile(path, []byte(key+"="+val+"\n"), 0o600))

	return path
}

func (s *DotenvTestSuite) TestExplicitPath() {
	const key = "HELMWAVE_TEST_DOTENV_EXPLICIT"
	s.T().Setenv(key, "marker")
	s.Require().NoError(os.Unsetenv(key))

	dir := s.T().TempDir()
	path := s.writeEnvFile(dir, "explicit.env", key, "from_explicit")

	helper.Dotenv(path)
	s.Require().Equal("from_explicit", os.Getenv(key))
}

func (s *DotenvTestSuite) TestCwdFallback() {
	const key = "HELMWAVE_TEST_DOTENV_CWD"
	s.T().Setenv(key, "marker")
	s.Require().NoError(os.Unsetenv(key))

	dir := s.T().TempDir()
	s.writeEnvFile(dir, ".env", key, "from_cwd")
	s.T().Chdir(dir)

	helper.Dotenv("")
	s.Require().Equal("from_cwd", os.Getenv(key))
}

func (s *DotenvTestSuite) TestNoFileIsNoop() {
	s.T().Chdir(s.T().TempDir())

	helper.Dotenv("")
}
