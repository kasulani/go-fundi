package generate

// NewMetadata returns an instance of Metadata.
func NewMetadata(output, templates, values string) *Metadata {
	// tech-debt: convert  params (output, templates, values) to value types
	return &Metadata{
		output:    output,
		templates: templates,
		values:    values,
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

// NewProjectDirectoryStructureUseCase returns an instance of ProjectDirectoryStructureUseCase.
func NewProjectDirectoryStructureUseCase(
	structureCreator DirectoryStructureCreator,
) *ProjectDirectoryStructureUseCase {
	return &ProjectDirectoryStructureUseCase{structureCreator: structureCreator}
}

// NewTestConfigurationFile returns a test ConfigurationFile you can use in unit tests.
func NewTestConfigurationFile() *ConfigurationFile {
	return NewConfigurationFile(
		&Metadata{output: ".", templates: "./testdata"},
		Directories{
			&Directory{
				name: "project_root_directory",
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
