package behaviour

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cucumber/godog"
	assert "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	out struct {
		cmdOutput []byte
		error     error
	}

	// TestSpecifications used in behaviour tests.
	TestSpecifications struct {
		ctx      context.Context
		log      *zap.Logger
		failures []string
		out      *out
	}
)

// NewTestSpecifications provides a new instance of TestSpecifications.
func NewTestSpecifications() *TestSpecifications {
	l, err := loadLogger()
	if err != nil {
		log.Fatalf("failed to create logger: %q", err)
	}

	specs := &TestSpecifications{
		ctx: context.Background(),
		log: l,
		out: &out{},
	}

	specs.log.Info("register fail handler")
	assert.RegisterFailHandler(func(message string, _ ...int) {
		specs.failures = append(specs.failures, message)
	})

	return specs
}

func loadLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.OutputPaths = []string{"stdout"}
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.MessageKey = "message"
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true

	return cfg.Build()
}

// Loader bootstraps the tests.
func (specs *TestSpecifications) Loader(sc *godog.ScenarioContext) {
	specs.registerAllSteps(sc)

	sc.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if err == nil {
			return specs.ctx, nil
		}

		for _, failure := range specs.failures {
			fmt.Printf("scenario has failed with error: \n%s\n", failure)
		}

		specs.failures = []string{}

		return specs.ctx, nil
	})
}

// MustStop frees up all test resources.
func (specs *TestSpecifications) MustStop() {
	specs.log.Info("clean up directories created during the testing")
	if err := afero.NewOsFs().RemoveAll("./funditest"); err != nil {
		specs.log.Error("failed to remove test directory hierarchy", zap.Error(err))
	}
}

// MustClearState resets the state of the test.
func (specs *TestSpecifications) MustClearState(scenario *godog.Scenario) {
	specs.log.Info(fmt.Sprintf("clear any previous state before scenario: %s", scenario.Name))
	specs.out = &out{}
}

func (specs TestSpecifications) commandOutput() string {
	return strings.TrimSpace(string(specs.out.cmdOutput))
}
