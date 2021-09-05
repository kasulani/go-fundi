package fundi

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/goava/di"
	"github.com/kasulani/go-fundi/pkg/generate"
	"github.com/pterm/pterm"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type (
	// Command is cli command type.
	Command struct {
		*cobra.Command

		ctx context.Context
	}

	// Cmd interface defines AddTo method.
	Cmd interface {
		AddTo(root *RootCommand)
	}

	// Commands is a slice of Cmd.
	Commands []Cmd

	// spinner type provides a way to track progress of background tasks.
	spinner struct {
		printer *pterm.SpinnerPrinter
	}

	// RootCommand of the cli application.
	RootCommand     Command
	scaffoldCommand Command
	generateCommand Command
	filesCommand    Command

	configFile struct {
		Version  int `yaml:"version"`
		Metadata struct {
			Name string `yaml:"name"`
			Path string `yaml:"path"`
		} `yaml:"metadata"`
		Structure []interface{} `yaml:"structure"`
	}

	ymlConfig struct {
		fs      afero.Fs
		file    *configFile
		spinner *spinner
		Flag    struct{ file string }
	}

	directoryCreator struct {
		spinner *spinner
		fs      afero.Fs
	}

	filesCreator struct {
		spinner *spinner
		fs      afero.Fs
	}
)

func provideCliCommands() di.Option {
	return di.Options(
		di.Provide(afero.NewOsFs),
		di.Provide(newSpinner),
		di.Provide(newDirectoryCreator, di.As(new(generate.HierarchyCreator))),
		di.Provide(newFilesCreator, di.As(new(generate.FileCreator))),
		di.Provide(newYmlConfig, di.As(new(generate.FundiFileReader))),
		generate.ProvideUseCases(),
		di.Provide(newRootCommand),
		di.Provide(newFilesCommand),
		di.Provide(newScaffoldCommand, di.As(new(Cmd))),
		di.Provide(newGenerateCommand, di.As(new(Cmd))),
	)
}

func newRootCommand() *RootCommand {
	return &RootCommand{
		Command: &cobra.Command{
			Use:     "fundi",
			Short:   "fundi is a scaffolding and code generation cli tool",
			Long:    `fundi is a scaffolding and code generation cli tool`,
			Version: "1.0.0",
		},
	}
}

func newScaffoldCommand(
	ctx context.Context,
	generateStructure *generate.DirectoryStructure,
	reader generate.FundiFileReader,
) *scaffoldCommand {
	cmd := &scaffoldCommand{
		Command: &cobra.Command{
			Use:     "scaffold",
			Aliases: []string{"scaf", "sca"},
			Short:   "scaffold a new project directory structure only",
			Long:    `use this command to generate a directory structure for a new project.`,
			Run: func(cmd *cobra.Command, args []string) {
				if err := generateStructure.UseCase(); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				os.Exit(0)
			},
		},
		ctx: ctx,
	}

	cmd.Flags().StringVarP(
		&reader.(*ymlConfig).Flag.file,
		"file",
		"f",
		"./.fundi.yaml",
		"file path to .fundi.yaml",
	)

	return cmd
}

func newFilesCommand(ctx context.Context, emptyFiles *generate.EmptyFiles) *filesCommand {
	cmd := &filesCommand{
		Command: &cobra.Command{
			Use:     "files",
			Aliases: []string{"file", "fil"},
			Short:   "add generated files to the project directory structure",
			Long:    `use this command to add generated files to the project directory as specified in the .fundi.yml file`,
			Run: func(cmd *cobra.Command, args []string) {
				skipTemplates, err := cmd.Flags().GetBool("skip-templates")
				if err != nil {
					println(err)
				}

				switch skipTemplates {
				case true:
					if err := emptyFiles.UseCase(); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				case false:
					fmt.Println("using templates")
				}
				os.Exit(0)
			},
		},
		ctx: ctx,
	}

	cmd.Flags().Bool("skip-templates", false, "generate empty files")

	return cmd
}

func newGenerateCommand(ctx context.Context, reader generate.FundiFileReader, filesCmd *filesCommand) *generateCommand {
	genCmd := &generateCommand{
		Command: &cobra.Command{
			Use:     "generate",
			Example: "generate files",
			Aliases: []string{"gene", "gen"},
			Short:   "generate project assets",
			Long:    `generate project assets`,
		},
		ctx: ctx,
	}
	genCmd.AddCommand(filesCmd.Command)
	genCmd.PersistentFlags().StringVarP(
		&reader.(*ymlConfig).Flag.file,
		"use-config",
		"c",
		"./.fundi.yaml",
		"path to fundi config file",
	)

	return genCmd
}

func (sc *scaffoldCommand) AddTo(root *RootCommand) {
	root.AddCommand(sc.Command)
}

func (gc *generateCommand) AddTo(root *RootCommand) {
	root.AddCommand(gc.Command)
}

func (file *configFile) asFundiFile() *generate.FundiFile {
	ff := new(generate.FundiFile)

	ff.Metadata.Name = file.Metadata.Name
	ff.Metadata.Path = file.Metadata.Path
	ff.Structure = file.Structure

	return ff
}

func newYmlConfig(fs afero.Fs, tracker *spinner) *ymlConfig {
	return &ymlConfig{
		fs:      fs,
		file:    new(configFile),
		spinner: tracker,
		Flag: struct {
			file string
		}{file: "."},
	}
}

func (reader *ymlConfig) Read() (*generate.FundiFile, error) {
	spin := reader.spinner.start("Reading fundi file...")
	data, err := afero.ReadFile(reader.fs, reader.Flag.file)
	if err != nil {
		spin.message("Reading fundi file: failed ✗").asFailure()

		return nil, err
	}
	spin.message("Reading fundi file: finished ✓").asSuccessful()

	spin = reader.spinner.start("Unmarshalling file data...")
	if err := yaml.Unmarshal(data, reader.file); err != nil {
		spin.message("Unmarshalling file data: failed ✗").asFailure()

		return nil, err
	}
	spin.message("Unmarshalling file data: finished ✓").asSuccessful()

	return reader.file.asFundiFile(), err
}

func newDirectoryCreator(fs afero.Fs, tracker *spinner) *directoryCreator {
	return &directoryCreator{spinner: tracker, fs: fs}
}

func (creator *directoryCreator) CreateHierarchy(hierarchy []string) error {
	spin := creator.spinner.start("Creating directory hierarchy...")
	for _, h := range hierarchy {
		spin.message(fmt.Sprintf("Creating directory hierarchy: %s...", h))
		if err := creator.fs.MkdirAll(h, 0755); err != nil {
			spin.message("Creating directory hierarchy: failed ✗").asFailure()

			return err
		}
	}
	spin.message("Creating directory hierarchy: finished ✓").asSuccessful()

	return nil
}

func newSpinner() *spinner {
	return new(spinner)
}

func (sp *spinner) start(msg string) *spinner {
	printer, err := pterm.DefaultSpinner.Start(msg)
	if err != nil {
		log.Fatalf("failed to initialise spinner printer: %s", err)
	}

	sp.printer = printer

	return sp
}

func (sp *spinner) message(msg string) *spinner {
	sp.printer.UpdateText(msg)

	return sp
}

func (sp *spinner) asSuccessful() {
	sp.printer.Success()
}

func (sp *spinner) asFailure() {
	sp.printer.Fail()
}

func newFilesCreator(fs afero.Fs, tracker *spinner) *filesCreator {
	return &filesCreator{spinner: tracker, fs: fs}
}

func (creator *filesCreator) CreateFiles(files map[string][]byte) error {
	spin := creator.spinner.start("Creating files...")
	for name, data := range files {
		spin.message(fmt.Sprintf("Creating files: %s...", name))
		if err := afero.WriteFile(creator.fs, name, data, 0644); err != nil {
			spin.message("Creating files: failed ✗").asFailure()

			return err
		}
	}
	spin.message("Creating files: finished ✓").asSuccessful()

	return nil
}
