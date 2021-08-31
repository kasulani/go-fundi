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

	ymlFile struct {
		Version  int `yaml:"version"`
		Metadata struct {
			Name string `yaml:"name"`
			Path string `yaml:"path"`
		} `yaml:"metadata"`
		Structure []interface{} `yaml:"structure"`
	}

	fileReader struct {
		fs      afero.Fs
		ymlFile *ymlFile
		spinner *spinner
		Flag    struct{ file string }
	}

	directoryCreator struct {
		spinner *spinner
		fs      afero.Fs
	}
)

func provideCliCommands() di.Option {
	return di.Options(
		di.Provide(afero.NewOsFs),
		di.Provide(newSpinner),
		di.Provide(newDirectoryCreator, di.As(new(generate.HierarchyCreator))),
		di.Provide(newYmlFileReader, di.As(new(generate.FundiFileReader))),
		generate.ProvideUseCases(),
		di.Provide(newRootCommand),
		di.Provide(newScaffoldCommand, di.As(new(Cmd))),
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
		&reader.(*fileReader).Flag.file,
		"file",
		"f",
		"./.fundi.yaml",
		"file path to .fundi.yaml",
	)

	return cmd
}

func (sc *scaffoldCommand) AddTo(root *RootCommand) {
	root.AddCommand(sc.Command)
}

func (file *ymlFile) toFundiFile() *generate.FundiFile {
	ff := new(generate.FundiFile)

	ff.Metadata.Name = file.Metadata.Name
	ff.Metadata.Path = file.Metadata.Path
	ff.Structure = file.Structure

	return ff
}

func newYmlFileReader(fs afero.Fs, tracker *spinner) *fileReader {
	return &fileReader{
		fs:      fs,
		ymlFile: new(ymlFile),
		spinner: tracker,
		Flag: struct {
			file string
		}{file: "."},
	}
}

func (reader *fileReader) Read() (*generate.FundiFile, error) {
	spin := reader.spinner.start("Reading fundi file...")
	data, err := afero.ReadFile(reader.fs, reader.Flag.file)
	if err != nil {
		spin.message("Reading fundi file: failed ✗").asFailure()

		return nil, err
	}
	spin.message("Reading fundi file: finished ✓").asSuccessful()

	spin = reader.spinner.start("Unmarshalling file data...")
	if err := yaml.Unmarshal(data, reader.ymlFile); err != nil {
		spin.message("Unmarshalling file data: failed ✗").asFailure()

		return nil, err
	}
	spin.message("Unmarshalling file data: finished ✓").asSuccessful()

	return reader.ymlFile.toFundiFile(), err
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
