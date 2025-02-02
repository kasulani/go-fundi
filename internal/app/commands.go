package app

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kasulani/go-fundi/internal/generate"
)

type (
	// Command is cli command type.
	Command struct {
		*cobra.Command
	}

	rootCommand Command

	// SubCommand interface defines AddTo method.
	SubCommand interface {
		AddTo(root *rootCommand)
	}

	// subCommands is a slice of SubCommand.
	subCommands []SubCommand

	generateProjectCommand Command
)

func newRootCommand() *rootCommand {
	return &rootCommand{
		Command: &cobra.Command{
			Use:     "fundi",
			Short:   "fundi is a scaffolding and code generation cli tool",
			Long:    `fundi is a scaffolding and code generation cli tool`,
			Version: "1.1.0",
		},
	}
}

func newGenerateProjectCommand(
	ctx context.Context,
	reader *fileReader,
	useCase *generate.ProjectUseCase,
) *generateProjectCommand {
	var filePath string

	cmd := &generateProjectCommand{
		&cobra.Command{
			Use:   "generate",
			Short: "generate your project directory structure and files",
			Long:  `use this subcommand to generate your project directory structure and files.`,
			Run: func(cmd *cobra.Command, args []string) {
				yamlFile, err := reader.readYAMLFile(filePath)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				if err := useCase.ScaffoldProject(ctx, yamlFile.toConfigurationFile()); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				os.Exit(0)
			},
		},
	}

	cmd.PersistentFlags().StringVarP(
		&filePath,
		"config-file",
		"f",
		"./.fundi.yaml",
		"path to your config file",
	)

	return cmd
}

func (cmd *generateProjectCommand) AddTo(root *rootCommand) {
	root.AddCommand(cmd.Command)
}
