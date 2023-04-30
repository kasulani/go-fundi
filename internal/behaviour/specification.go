package behaviour

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cucumber/godog"
	assert "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	cmd struct {
		output []byte
		error  error
	}

	// Test used in behaviour tests.
	Test struct {
		ctx        context.Context
		log        *zap.Logger
		failures   []string
		cmd        *cmd
		ConfigFile string
	}
)

const testDir = "./funditest"

// NewTestSpecifications provides a new instance of Test.
func NewTestSpecifications() *Test {
	l, err := loadLogger()
	if err != nil {
		log.Fatalf("failed to create logger: %q", err)
	}

	specs := &Test{
		ctx: context.Background(),
		log: l,
		cmd: &cmd{},
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
func (test *Test) Loader(sc *godog.ScenarioContext) {
	test.registerAllSteps(sc)

	sc.After(func(ctx context.Context, s *godog.Scenario, err error) (context.Context, error) {
		if err == nil {
			return test.ctx, nil
		}

		for _, failure := range test.failures {
			fmt.Printf("scenario has failed with error: \n%s\n", failure)
		}

		test.failures = []string{}

		return test.ctx, nil
	})
}

// MustStop frees up all test resources.
func (test *Test) MustStop() {
	test.log.Info("clean up directories created during the testing")
	if err := afero.NewOsFs().RemoveAll(testDir); err != nil {
		test.log.Error("failed to remove test directory hierarchy", zap.Error(err))
	}
}

// MustClearState resets the state of the test.
func (test *Test) MustClearState(scenario *godog.Scenario) {
	test.log.Info(fmt.Sprintf("clear any previous state before scenario: %s", scenario.Name))
	if err := afero.NewOsFs().RemoveAll(testDir); err != nil {
		test.log.Fatal("failed to remove test directory hierarchy", zap.Error(err))
	}
}

func (test *Test) commandOutput() string {
	return strings.TrimSpace(string(test.cmd.output))
}

func (test *Test) workingDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		test.log.Fatal("failed to get working directory", zap.Error(err))
	}

	return dir
}
