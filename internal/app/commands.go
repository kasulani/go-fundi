package app

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kasulani/go-fundi/internal/generate"
)

type (
	generateProjectCommand Command
)

func newGenerateProjectCommand(
	ctx context.Context,
	reader *fileReader,
	useCase *generate.ProjectUseCase,
) *generateProjectCommand {
	var filePath string

	cmd := &generateProjectCommand{
		&cobra.Command{
			Use:   "generate-cmd",
			Short: "generate your project directory structure and files",
			Long:  `use this subcommand to generate your project directory structure and files.`,
			Run: func(cmd *cobra.Command, args []string) {
				yamlFile, err := reader.readYAMLFile(filePath)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				output := generate.All
				flags := []string{generate.DirectoriesOnly, generate.EmptyFiles}

				for _, flag := range flags {
					isSet, err := cmd.Flags().GetBool(flag)
					if err != nil {
						println(err)
					}

					if isSet {
						output = flag
						break
					}
				}

				if err := useCase.ScaffoldProject(ctx, output, yamlFile.toConfigurationFile()); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				os.Exit(0)
			},
		},
	}

	cmd.Flags().Bool("directories-only", false, "generate project directories")
	cmd.Flags().Bool("empty-files", false, "generate project directories and empty files")

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
