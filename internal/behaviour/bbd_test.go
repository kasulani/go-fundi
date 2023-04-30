//go:build behaviour
// +build behaviour

package behaviour

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
)

func TestBehaviour(t *testing.T) {
	specs := NewTestSpecifications()
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

func initializeSuite(test *Test) func(ts *godog.TestSuiteContext) {
	return func(ts *godog.TestSuiteContext) {
		ts.AfterSuite(func() {
			test.MustStop()
		})
	}
}

func initializeScenarios(test *Test) func(sc *godog.ScenarioContext) {
	return func(sc *godog.ScenarioContext) {
		test.Loader(sc)
		sc.BeforeScenario(func(s *godog.Scenario) {
			test.MustClearState(s)
		})
	}
}
