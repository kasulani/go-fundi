package fundi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/goava/di"
	"github.com/pterm/pterm"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/kasulani/go-fundi/internal/generate"
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

	// spinner type provides a way to track progress of background tasks.
	spinner struct {
		printer *pterm.SpinnerPrinter
	}

	// rootCommand of the cli application.
	rootCommand     Command
	scaffoldCommand Command
	generateCommand Command
	filesCommand    Command

	configFile struct {
		Version  int `yaml:"version"`
		Metadata struct {
			Name      string `yaml:"name"`
			Path      string `yaml:"path"`
			Templates struct {
				Path string `yaml:"path"`
			} `yaml:"templates"`
		} `yaml:"metadata"`
		Structure []interface{} `yaml:"structure"`
	}

	ymlConfig struct {
		fs      afero.Fs
		file    *configFile
		spinner *spinner
		Flag    struct{ file string }
	}

	structureCreator struct {
		spinner *spinner
		fs      afero.Fs
	}

	filesCreator struct {
		spinner *spinner
		fs      afero.Fs
	}

	templateParser struct {
		spinner *spinner
		fs      afero.Fs
	}
)

func provideCliCommands() di.Option {
	return di.Options(
		di.Provide(afero.NewOsFs),
		di.Provide(newSpinner),
		di.Provide(newStructureCreator, di.As(new(generate.StructureCreator))),
		di.Provide(newFilesCreator, di.As(new(generate.FileCreator))),
		di.Provide(newYmlConfig, di.As(new(generate.FundiFileReader))),
		di.Provide(newTemplateParser, di.As(new(generate.TemplateParser))),
		generate.ProvideUseCases(),
		di.Provide(newRootCommand),
		di.Provide(newFilesCommand),
		di.Provide(newScaffoldCommand),
		di.Provide(newGenerateCommand, di.As(new(Cmd))),
	)
}

func newRootCommand() *rootCommand {
	return &rootCommand{
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
	directoryStructure *generate.DirectoryStructure,
) *scaffoldCommand {
	cmd := &scaffoldCommand{
		Command: &cobra.Command{
			Use:     "scaffold",
			Aliases: []string{"scaf", "sca"},
			Short:   "scaffold a new project directory structure only",
			Long:    `use this command to generate a directory structure for a new project.`,
			Run: func(cmd *cobra.Command, args []string) {
				if err := directoryStructure.UseCase(); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				os.Exit(0)
			},
		},
		ctx: ctx,
	}

	return cmd
}

func newFilesCommand(
	ctx context.Context,
	emptyFiles *generate.EmptyFiles,
	filesFromTemplates *generate.FilesFromTemplates) *filesCommand {
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
					if err := filesFromTemplates.UseCase(); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				}
				os.Exit(0)
			},
		},
		ctx: ctx,
	}

	cmd.Flags().Bool("skip-templates", false, "generate empty files")

	return cmd
}

func newGenerateCommand(
	ctx context.Context,
	reader generate.FundiFileReader,
	filesCmd *filesCommand,
	scaffoldCmd *scaffoldCommand,
) *generateCommand {
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
	genCmd.AddCommand(scaffoldCmd.Command)
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

func (gc *generateCommand) AddTo(root *rootCommand) {
	root.AddCommand(gc.Command)
}

func (file *configFile) asFundiFile() *generate.FundiFile {
	ff := new(generate.FundiFile)

	ff.Metadata.Name = file.Metadata.Name
	ff.Metadata.Path = file.Metadata.Path
	ff.Metadata.Templates.Path = file.Metadata.Templates.Path
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

func newStructureCreator(fs afero.Fs, tracker *spinner) *structureCreator {
	return &structureCreator{spinner: tracker, fs: fs}
}

func (creator *structureCreator) CreateStructure(folders []string) error {
	spin := creator.spinner.start("Creating directory hierarchy...")
	for _, h := range folders {
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

func newTemplateParser(fs afero.Fs, tracker *spinner) *templateParser {
	return &templateParser{spinner: tracker, fs: fs}
}

func (tp *templateParser) ParseTemplates(
	data map[string]*generate.TemplateFile,
	templatePath string,
) (map[string][]byte, error) {
	if templatePath == "" {
		return nil, errors.New("template path missing")
	}

	buffer := new(bytes.Buffer)
	parsedFiles := make(map[string][]byte)
	spin := tp.spinner.start("Parsing templates...")

	for name, tpl := range data {
		buffer.Reset()
		if tpl.Name == "" {
			parsedFiles[name] = buffer.Bytes()
			continue
		}

		spin.message(fmt.Sprintf("Parsing templates: reading file %s...", tpl.Name))
		contents, err := afero.ReadFile(tp.fs, templatePath+string(os.PathSeparator)+tpl.Name)
		if err != nil {
			spin.message("Parsing templates: failed ✗").asFailure()

			return nil, err
		}

		spin.message(fmt.Sprintf("Parsing templates: processing %s...", tpl.Name))
		tmpl := template.Must(template.New("tmpl").Parse(string(contents)))
		if err := tmpl.Execute(buffer, tpl.Values); err != nil {
			spin.message("Parsing templates: failed ✗").asFailure()

			return nil, err
		}

		parsedFiles[name] = buffer.Bytes()
	}
	spin.message("Parsing templates: finished ✓").asSuccessful()

	return parsedFiles, nil
}
