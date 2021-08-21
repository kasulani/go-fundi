package fundi

import (
	"context"
	"fmt"
	"os"

	"github.com/goava/di"
	"github.com/spf13/cobra"

	"github.com/kasulani/go-fundi/pkg/scaffold"
)

type (
	// Command is cli command type.
	Command struct {
		*cobra.Command

		ctx context.Context
	}

	// Cmd interface defines AddTo method.
	Cmd interface {
		AddTo(root *rootCommand)
	}

	// Commands is a slice of Cmd.
	Commands []Cmd

	rootCommand Command
	initCommand Command
	// todo: add command to generate fundi config file
)

func provideCliCommands() di.Option {
	return di.Options(
		di.Provide(context.Background),
		di.Provide(newConfig),
		provideRepository(),
		scaffold.ProvideUseCases(),
		di.Provide(NewInputValidator, di.As(new(InputValidator))),
		di.Provide(newRootCommand),
		di.Provide(newInitCommand, di.As(new(Cmd))),
	)
}

func newRootCommand() *rootCommand {
	return &rootCommand{
		Command: &cobra.Command{
			Use:     "fundi",
			Short:   "fundi is a cli tool that scaffolds a go project for you",
			Long:    "fundi is a cli tool that scaffolds a go project for you",
			Version: "1.0.0",
		},
	}
}

func newInitCommand(ctx context.Context) *initCommand {
	// todo: flags - path to config manifest
	return &initCommand{
		Command: &cobra.Command{
			Use:   "init",
			Short: "initialize a new go project",
			Long:  "This command initializes a new go project",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Todo: scaffolding...")
				// todo: get input
				// todo: validate input e.g - add tag to validate file/path exits
				// todo: invoke use case passing the valid input
				// todo: return appropriate response
				os.Exit(0)
			},
		},
		ctx: ctx,
	}
}

func (c *initCommand) AddTo(root *rootCommand) {
	root.AddCommand(c.Command)
}
