package generate

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

	// ProjectDirectoryStructure represents the project directory tree.
	ProjectDirectoryStructure struct {
		output      string
		directories []string
	}
)

func (d *Directory) hasSubDirectories() bool {
	return d.subDirectories != nil
}
