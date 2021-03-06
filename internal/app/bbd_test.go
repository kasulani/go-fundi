// +build behaviour

package app

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/kasulani/go-fundi/internal/behaviour"
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
	os.Chdir("../../")
	parentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %q", err)
	}

	return []string{parentDir + "/features"}
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
