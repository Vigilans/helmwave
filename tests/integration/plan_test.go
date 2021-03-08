package integration

import (
	"github.com/stretchr/testify/suite"
	"github.com/zhilyaev/helmwave/pkg/helmwave"
	"github.com/zhilyaev/helmwave/pkg/kubedog"
	"github.com/zhilyaev/helmwave/pkg/template"
	"github.com/zhilyaev/helmwave/pkg/yml"
	helm "helm.sh/helm/v3/pkg/cli"
	"testing"
)

const (
	PlanPath = ".helmwave/"
	PlanFile = "planfile"
)

type PlanTestSuite struct {
	suite.Suite
	app *helmwave.Config
}

func (s *PlanTestSuite) SetupTest() {
	s.app = &helmwave.Config{
		Helm: &helm.EnvSettings{},
		Tpl: template.Tpl{
			From: "../fixtures/helmwave.yml.tpl",
			To:   "helmwave.yml",
		},
		Yml:      yml.Config{},
		PlanPath: PlanPath,
		Logger: &helmwave.Log{
			Level:  "DEBUG",
			Format: "flat",
			Color:  false,
		},
		Parallel: false,
		Kubedog:  &kubedog.Config{},
	}

	err := s.app.InitLogger()
	s.Require().NoError(err)

	err = s.app.Tpl.Render()
	s.Require().NoError(err)

	err = yml.Read(s.app.Tpl.To, &s.app.Yml)
	s.Require().NoError(err)
}

func (s *PlanTestSuite) TestPlanReleases() {
	opts := &yml.SavePlanOptions{}

	opts.File(s.app.PlanPath + PlanFile).Dir(s.app.PlanPath)

	opts.PlanReleases()
	opts.Tags(s.app.Tags.Value())

	err := s.app.Yml.Plan(opts, s.app.Helm)
	s.Require().NoError(err)

	s.FileExists(PlanPath + ".manifest/nginx@test-nginx.yml")
}

func TestPlanTestSuite(t *testing.T) {
	suite.Run(t, new(PlanTestSuite))
}
