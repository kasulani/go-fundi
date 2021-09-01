package behaviour

import (
	"os"

	"github.com/cucumber/godog"
	"go.uber.org/zap"
)

func (specs *TestSpecifications) workingDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		specs.log.Fatal("failed to get working directory", zap.Error(err))
	}

	return dir
}

func (specs *TestSpecifications) loadFixtures(scenario *godog.Scenario) {
	specs.log.Info("load fixtures")
	switch scenario.Name {
	case "Scaffold command exits with code 0":
		specs.in.File = specs.workingDirectory() + "/pkg/behaviour/.bdd.test.fundi.yaml"
	case "Scaffold command exits with code 1":
		specs.in.File = specs.workingDirectory() + "/pkg/behaviour/.non-existing.fundi.yaml"
	default:
		return
	}
}
