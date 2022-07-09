package app

import (
	"context"
	"log"
	"os"

	"github.com/goava/di"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/afero"

	"github.com/kasulani/go-fundi/internal/generate"
)

type (
	config struct {
		LogLevel string `envconfig:"LOG_LEVEL" default:"debug"`
	}
)

// Container is a dependency injection container.
func Container() *di.Container {
	if os.Getenv("LOG_LEVEL") == "debug" {
		di.SetTracer(&di.StdTracer{})
	}

	container, err := di.New(
		di.Provide(context.Background),
		di.Provide(newConfig),
		di.Provide(afero.NewOsFs),
		di.Provide(newSpinner),
		di.Provide(newStructureCreator, di.As(new(generate.StructureCreator))),
		di.Provide(newFilesCreator, di.As(new(generate.FileCreator))),
		di.Provide(newYmlConfig, di.As(new(generate.FundiFileReader))),
		di.Provide(newTemplateParser, di.As(new(generate.TemplateParser))),
		generate.ProvideUseCases(),
		provideCLICommands(),
		di.Invoke(registerSubCommands),
	)

	if err != nil {
		log.Fatalf("failed to create DI container: %q", err)
	}

	return container
}

// Run is a high level app entry function.
func Run(root *rootCommand) error {
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

// registerSubCommands adds all the sub commands to the root command.
func registerSubCommands(root *rootCommand, commands subCommands) {
	for _, command := range commands {
		command.AddTo(root)
	}
}
