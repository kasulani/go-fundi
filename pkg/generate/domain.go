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

	// HierarchyCreatorFunc is an adapter type to allow use of ordinary functions as directory HierarchyCreator.
	HierarchyCreatorFunc func(hierarchy []string) error

	// FundiFileReader interface defines the Read method.
	FundiFileReader interface {
		Read() (*FundiFile, error)
	}

	// FundiFileReaderFunc is an adapter type to allow use of ordinary functions as fundi file readers.
	FundiFileReaderFunc func() (*FundiFile, error)

	// FileCreator interface define the CreateFiles method.
	FileCreator interface {
		CreateFiles(files map[string][]byte) error
	}

	// FileCreatorFunc is an adapter type to allow use of ordinary functions as directory FileCreator.
	FileCreatorFunc func(files map[string][]byte) error
)

// CreateHierarchy creates a directory hierarchy.
func (maker HierarchyCreatorFunc) CreateHierarchy(hierarchy []string) error {
	return maker(hierarchy)
}

// Read wraps the reader function fn.
func (fn FundiFileReaderFunc) Read() (*FundiFile, error) {
	return fn()
}

// CreateFiles wraps the file creator function fn.
func (fn FileCreatorFunc) CreateFiles(files map[string][]byte) error {
	return fn(files)
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
			if isDirectory(cast.ToStringMap(item)) && hasContents(cast.ToStringMap(item)) {
				dirs, err := getAllDirectories(cast.ToSlice(dict["contains"]))
				if err != nil {
					return directories, err
				}

				parent := cast.ToString(dict["folder"])
				if len(dirs) > 0 {
					for _, dir := range dirs {
						d := parent + string(os.PathSeparator) + dir
						directories = append(directories, d)
					}
				} else {
					directories = append(directories, parent)
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

func hasContents(item map[string]interface{}) bool {
	dict := cast.ToStringMap(item)

	return len(cast.ToSlice(dict["contains"])) > 0
}

func isDirectory(item map[string]interface{}) bool {
	_, yes := item["folder"]

	return yes
}

func getFilesSkipTemplates(structure []interface{}) ([]string, error) {
	var files []string
	for _, item := range structure {
		kind := reflect.ValueOf(item).Kind()
		switch kind {
		case reflect.Map:
			dict := cast.ToStringMap(item)
			if isDirectory(cast.ToStringMap(item)) && hasContents(cast.ToStringMap(item)) {
				allFiles, err := getFilesSkipTemplates(cast.ToSlice(dict["contains"]))
				if err != nil {
					return allFiles, err
				}

				parent := cast.ToString(dict["folder"])
				for _, file := range allFiles {
					files = append(files, parent+string(os.PathSeparator)+file)
				}
			} else if isFile(cast.ToStringMap(item)) {
				files = append(files, cast.ToString(dict["file"]))
			}
		default:
			return nil, errors.Errorf("unexpected kind: %s", kind)
		}
	}
	return files, nil
}

func isFile(item map[string]interface{}) bool {
	_, yes := item["file"]

	return yes
}

func generateEmptyFiles(paths []string) map[string][]byte {
	files := make(map[string][]byte)

	for _, path := range paths {
		files[path] = []byte("")
	}

	return files
}
