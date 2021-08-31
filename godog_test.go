package bdd

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/kasulani/go-fundi/pkg/behaviour"
)

var specs *behaviour.TestSpecifications

func InitializeSuite(ts *godog.TestSuiteContext) {
	ts.BeforeSuite(func() {
		specs = behaviour.NewTestSpecifications()
	})
	ts.AfterSuite(func() {
		specs.MustStop()
	})
}

func InitializeScenario(sc *godog.ScenarioContext) {
	specs.Loader(sc)
	sc.Before(func(ctx context.Context, s *godog.Scenario) (context.Context, error) {
		specs.MustClearState(s)

		return ctx, nil
	})
}
