package generate

import (
	"os"
	"reflect"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type (
	// FundiFile is a model of the .fundi.yaml file.
	FundiFile struct {
		// Metadata of the project.
		Metadata struct {
			// Name of the project.
			Name string
			// Path is the location of the root directory.
			Path string
		}
		// Structure of the project.
		Structure []interface{}
	}

	// HierarchyCreator interface defines the CreateHierarchy method.
	HierarchyCreator interface {
		CreateHierarchy(hierarchy []string) error
	}

	// HierarchyCreatorFunc is an adapter type to allow use of ordinary functions as directory hCreator.
	HierarchyCreatorFunc func(hierarchy []string) error

	// FundiFileReader interface defines the Read method.
	FundiFileReader interface {
		Read() (*FundiFile, error)
	}

	// FundiFileReaderFunc is an adapter type to allow use of ordinary functions as fundiFile readers.
	FundiFileReaderFunc func() (*FundiFile, error)
)

// CreateHierarchy creates a directory hierarchy.
func (maker HierarchyCreatorFunc) CreateHierarchy(hierarchy []string) error {
	return maker(hierarchy)
}

// Read wraps the reader function fn.
func (fn FundiFileReaderFunc) Read() (*FundiFile, error) {
	return fn()
}

func generateHierarchy(root string, dirs []string) []string {
	hierarchy := make([]string, 0, len(dirs))

	for _, dir := range dirs {
		hierarchy = append(hierarchy, root+string(os.PathSeparator)+dir)
	}

	return hierarchy
}

func getAllDirectories(structure []interface{}) ([]string, error) {
	var directories []string
	for _, item := range structure {
		kind := reflect.ValueOf(item).Kind()
		switch kind {
		case reflect.Map:
			dict := cast.ToStringMap(item)
			if isDirectory(cast.ToStringMap(item)) && hasContents(cast.ToSlice(dict["contains"])) {
				parent := cast.ToString(dict["folder"])

				dirs, err := getAllDirectories(cast.ToSlice(dict["contains"]))
				if err != nil {
					return directories, err
				}

				for _, dir := range dirs {
					d := parent + string(os.PathSeparator) + dir
					directories = append(directories, d)
				}
			} else if isDirectory(cast.ToStringMap(item)) {
				directories = append(directories, cast.ToString(dict["folder"]))
			}
		default:
			return nil, errors.Errorf("unexpected kind: %s", kind)
		}
	}

	return directories, nil
}

func hasContents(contains []interface{}) bool {
	return len(contains) > 0
}

func isDirectory(item map[string]interface{}) bool {
	_, yes := item["folder"]

	return yes
}
