package app

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"

	"github.com/kasulani/go-fundi/internal/generate"
)

type (
	metadata struct {
		Output    string `yaml:"output"`
		Templates string `yaml:"templates"`
		Values    string `yaml:"values"`
	}

	file struct {
		Name     string `yaml:"name"`
		Template string `yaml:"template"`
	}

	files []*file

	directory struct {
		Name           string      `yaml:"name"`
		Files          files       `yaml:"files"`
		SubDirectories directories `yaml:"directories"`
	}

	directories []*directory

	yamlFile struct {
		Metadata    *metadata   `yaml:"metadata"`
		Directories directories `yaml:"directories"`
	}

	fileReader struct{ fs afero.Fs }

	directoryCreator struct{ fs afero.Fs }

	filesCreator struct{ fs afero.Fs }
)

// readYAMLFile returns an instance of yamlFile.
func (fr *fileReader) readYAMLFile(filepath string) (*yamlFile, error) {
	data, err := afero.ReadFile(fr.fs, filepath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file %s", filepath)
	}

	var cfg yamlFile
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal YAML data")
	}

	return &cfg, nil
}

func (yf *yamlFile) toConfigurationFile() *generate.ConfigurationFile {
	dirs := make(generate.Directories, len(yf.Directories))
	for i, dir := range yf.Directories {
		dirs[i] = generate.NewDirectory(dir.Name, yf.convertFiles(dir.Files), yf.convertDirectories(dir.SubDirectories))
	}

	return generate.NewConfigurationFile(
		generate.NewMetadata(yf.Metadata.Output, yf.Metadata.Templates, yf.Metadata.Values),
		dirs,
	)
}

func (yf *yamlFile) convertFiles(fs files) generate.Files {
	if len(fs) == 0 {
		return nil
	}

	files := make(generate.Files, len(fs))
	for i, f := range fs {
		files[i] = generate.NewFile(f.Name, f.Template)
	}

	return files
}

func (yf *yamlFile) convertDirectories(ds directories) generate.Directories {
	if len(ds) == 0 {
		return nil
	}

	dirs := make(generate.Directories, len(ds))
	for i, d := range ds {
		dirs[i] = generate.NewDirectory(d.Name, yf.convertFiles(d.Files), yf.convertDirectories(d.SubDirectories))
	}

	return dirs
}

func (creator *directoryCreator) CreateDirectoryStructure(
	_ context.Context,
	output string,
	directories []string,
) error {
	dirs := directories
	if len(dirs) == 0 {
		fmt.Println("no files to create")

		return nil
	}

	bar, err := pterm.DefaultProgressbar.WithTotal(len(dirs)).WithTitle("Generating directories").Start()
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if err := creator.fs.MkdirAll(output+string(os.PathSeparator)+dir, 0755); err != nil {
			return errors.Wrapf(err, "failed to create directory %s", dir)
		}
		bar.Increment()
	}

	_, err = bar.Stop()
	if err != nil {
		return err
	}

	return nil
}

func (fc *filesCreator) CreateFiles(
	_ context.Context,
	metadata *generate.Metadata,
	templateFiles generate.FileTemplates,
) error {
	if len(templateFiles) == 0 {
		fmt.Println("no files to create")

		return nil
	}

	bar, err := pterm.DefaultProgressbar.WithTotal(len(templateFiles)).WithTitle("Generating files").Start()
	if err != nil {
		return err
	}

	templateValues, err := fc.getTemplateValues(metadata.GetValuesPath())
	if err != nil {
		return err
	}

	output := metadata.GetDestinationPath()
	templatePath := metadata.GetTemplatePath()

	for name, templateFile := range templateFiles {
		data, err := fc.parseTemplate(templatePath, templateFile, templateValues)
		if err != nil {
			return errors.Wrapf(err, "failed to parse template %s", templateFile)
		}

		destinationPath := output + string(os.PathSeparator) + name
		if err := afero.WriteFile(fc.fs, destinationPath, data, 0644); err != nil {
			_, _ = bar.Stop()

			return errors.Wrapf(err, "failed to create file %s", destinationPath)
		}
		bar.Increment()
	}

	_, err = bar.Stop()
	if err != nil {
		return err
	}

	return nil
}

func (fc *filesCreator) parseTemplate(
	templatePath,
	templateName string,
	values map[string]interface{},
) ([]byte, error) {
	buffer := new(bytes.Buffer)
	if templateName == "" {
		return buffer.Bytes(), nil
	}

	contents, err := afero.ReadFile(fc.fs, templatePath+string(os.PathSeparator)+templateName)
	if err != nil {
		return buffer.Bytes(), err
	}

	tmpl, err := template.New(templateName).Parse(string(contents))
	if err != nil {
		return nil, err
	}

	if err := tmpl.Execute(buffer, values[templateName]); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (fc *filesCreator) getTemplateValues(path string) (map[string]interface{}, error) {
	values := make(map[string]interface{})

	valuesFile, err := afero.ReadFile(fc.fs, path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read values file %s", path)
	}

	err = yaml.Unmarshal(valuesFile, &values)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal values file %s", path)
	}

	return values, nil
}
