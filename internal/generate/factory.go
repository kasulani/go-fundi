package generate

import "github.com/spf13/cast"

// NewMetadata returns an instance of Metadata.
func NewMetadata(metadata map[string]any) *Metadata {
	return &Metadata{
		output:    cast.ToString(metadata[MetaDataOutputKey]),
		templates: cast.ToString(metadata[MetaDataTemplatesKey]),
		values:    cast.ToString(metadata[MetaDataValuesKey]),
		variables: cast.ToStringMap(metadata[MetaDataVariablesKey]),
	}
}

// NewFile returns an instance of File.
func NewFile(name, template string) *File {
	// tech-debt: convert  params (name, template) to value types
	return &File{
		name:     name,
		template: template,
	}
}

// NewDirectory returns an instance of Directory.
func NewDirectory(name string, files Files, directories Directories) *Directory {
	return &Directory{
		name:           name,
		files:          files,
		subDirectories: directories,
	}
}

// NewConfigurationFile returns an instance of ConfigurationFile.
func NewConfigurationFile(metadata *Metadata, directories Directories) *ConfigurationFile {
	return &ConfigurationFile{metadata: metadata, directories: directories}
}

// NewProjectUseCase returns an instance of ProjectUseCase.
func NewProjectUseCase(structureCreator DirectoryStructureCreator, fileCreator FilesCreator) *ProjectUseCase {
	return &ProjectUseCase{structureCreator: structureCreator, filesCreator: fileCreator}
}

// NewTestConfigurationFile returns a test ConfigurationFile you can use in unit tests.
func NewTestConfigurationFile() *ConfigurationFile {
	return NewConfigurationFile(
		&Metadata{output: ".", templates: "./testdata"},
		Directories{
			&Directory{
				name:  "project_root_directory",
				files: Files{&File{name: "README.md", template: "README.md.tmpl"}},
				subDirectories: Directories{
					&Directory{
						name:  "cmd",
						files: Files{&File{name: "main.go", template: "main.go.tmpl"}},
					},
					&Directory{
						name: "internal",
						subDirectories: Directories{
							&Directory{
								name:           "domain",
								files:          Files{&File{name: "domain.go", template: "domain.go.tmpl"}},
								subDirectories: nil,
							},
						},
					},
				},
			},
		},
	)
}
