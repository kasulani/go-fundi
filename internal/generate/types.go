package generate

import "os"

type (
	// Metadata about the project.
	Metadata struct {
		output    string
		templates string
		values    string
	}

	// File in the project.
	File struct {
		name     string
		template string
	}

	// Files is a collection of File.
	Files []*File

	// Directory in the project structure.
	Directory struct {
		name           string
		files          Files
		subDirectories Directories
	}

	// Directories is a collection of Directory.
	Directories []*Directory

	// ConfigurationFile is the yaml file that specifies the project structure and the files that go into it.
	ConfigurationFile struct {
		metadata    *Metadata
		directories Directories
	}

	// FileTemplates is a map of file and its template.
	FileTemplates map[string]string
)

const (
	MetaDataOutputKey    = "output"
	MetaDataTemplatesKey = "templates"
	MetaDataValuesKey    = "values"
)

// GetDestinationPath returns destination path where the project will be created.
func (m *Metadata) GetDestinationPath() string {
	return m.output
}

// GetTemplatePath returns location of templates.
func (m *Metadata) GetTemplatePath() string {
	return m.templates
}

// GetValuesPath returns location of values.
func (m *Metadata) GetValuesPath() string {
	return m.values
}

func (d *Directory) hasSubDirectories() bool {
	return d.subDirectories != nil
}

func (cf *ConfigurationFile) getFilesAndTemplates() FileTemplates {
	fileTemplates := make(FileTemplates)

	for _, directory := range cf.directories {
		addFileAndTemplate(directory, fileTemplates, directory.name)
	}

	return fileTemplates
}

func addFileAndTemplate(directory *Directory, fileTemplates FileTemplates, prefix string) {
	for _, file := range directory.files {
		fileTemplates[prefix+string(os.PathSeparator)+file.name] = file.template
	}

	for _, subDirectory := range directory.subDirectories {
		addFileAndTemplate(subDirectory, fileTemplates, prefix+string(os.PathSeparator)+subDirectory.name)
	}
}
