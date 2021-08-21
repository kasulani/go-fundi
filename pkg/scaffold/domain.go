package scaffold

type (
	// Project represents the .fundi.yaml manifest file.
	Project struct {
		// Name of the project directory
		Name string
		// Kind of project; api or cli
		Kind string
		// MetaData of the project
		MetaData map[string]string
		// Specifications of the project
		Specifications interface{}
	}

	// ManifestReader interface defines the Read method.
	ManifestReader interface {
		Read() error
	}
)

// NewProject returns a new Project.
func NewProject() *Project {
	return &Project{}
}
