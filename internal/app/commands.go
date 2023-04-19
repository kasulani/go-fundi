package app

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kasulani/go-fundi/internal/generate"
)

type (
	generateDirectoriesCommand Command
)

func newGenerateDirectoriesCommand(
	ctx context.Context,
	reader *fileReader,
	useCase *generate.ProjectDirectoryStructureUseCase,
) *generateDirectoriesCommand {
	var filePath string

	cmd := &generateDirectoriesCommand{
		&cobra.Command{
			Use:   "generate-directories",
			Short: "generate project directory structure",
			Long:  `use this subcommand to generate your project directory structure.`,
			Run: func(cmd *cobra.Command, args []string) {
				yamlFile, err := reader.readYAMLFile(filePath)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				if err := useCase.GenerateProjectStructure(ctx, yamlFile.toConfigurationFile()); err != nil {
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
func (cmd *generateDirectoriesCommand) AddTo(root *rootCommand) {
	root.AddCommand(cmd.Command)
}
