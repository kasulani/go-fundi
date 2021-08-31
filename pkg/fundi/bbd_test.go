// +build behaviour

package fundi

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/kasulani/go-fundi/pkg/behaviour"
)

func TestBehaviour(t *testing.T) {
	specs := behaviour.NewTestSpecifications()
	suite := godog.TestSuite{
		Name:                 "fundi",
		TestSuiteInitializer: initializeSuite(specs),
		ScenarioInitializer:  initializeScenarios(specs),
		Options: &godog.Options{
			Randomize:     1,
			StopOnFailure: false,
			Format:        "pretty",
			Paths:         featuresFiles(t),
			Tags:          "~@notYetImplemented",
		},
	}

	if suite.Run() != 0 {
		t.Fatal("failed to run behaviour tests")
	}
}

func featuresFiles(t *testing.T) []string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %q", err)
	}

	parentDir := strings.Split(filepath.Dir(wd), string(os.PathSeparator))[1]

	return []string{"/" + parentDir + "/features"}
}

func initializeSuite(specs *behaviour.TestSpecifications) func(ts *godog.TestSuiteContext) {
	return func(ts *godog.TestSuiteContext) {
		ts.AfterSuite(func() {
			specs.MustStop()
		})
	}
}

func initializeScenarios(specs *behaviour.TestSpecifications) func(sc *godog.ScenarioContext) {
	return func(sc *godog.ScenarioContext) {
		specs.Loader(sc)
		sc.BeforeScenario(func(s *godog.Scenario) {
			specs.MustClearState(s)
		})
	}
}
