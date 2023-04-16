package app

import (
	"context"

	"github.com/pkg/errors"
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
	structure *generate.ProjectDirectoryStructure,
) error {
	dirs := structure.Directories()

	for _, dir := range dirs {
		if err := creator.fs.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}
