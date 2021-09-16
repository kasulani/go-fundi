package fundi

import (
	"context"
	"log"
	"os"

	"github.com/goava/di"
	"github.com/kelseyhightower/envconfig"
)

type (
	config struct {
		LogLevel string `envconfig:"LOG_LEVEL" default:"debug"`
	}
)

// Container is a dependency injection container.
func Container(selector string) *di.Container {
	if os.Getenv("LOG_LEVEL") == "debug" {
		di.SetTracer(&di.StdTracer{})
	}

	var c *di.Container
	var err error

	switch selector {
	case "cli":
		c, err = di.New(
			di.Provide(context.Background),
			di.Provide(newConfig),
			provideCliCommands(),
			di.Invoke(RegisterCliCommands),
		)
	default:
		log.Fatalf("unknown container selector: %s", selector)
	}

	if err != nil {
		log.Fatalf("failed to create DI container: %q", err)
	}

	return c
}

// StartCLI is a high level cli entry function.
func StartCLI(root *rootCommand) error {
	return root.Execute()
}

// newConfig returns config.
func newConfig() *config {
	cfg := new(config)
	err := envconfig.Process("", cfg)

	if err != nil {
		log.Fatalf("failed to load configuration: %q", err)
	}

	return cfg
}

// RegisterCliCommands adds all the sub commands to the root command.
func RegisterCliCommands(root *rootCommand, commands Commands) {
	for _, command := range commands {
		command.AddTo(root)
	}
}
