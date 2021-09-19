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
		Metadata struct {
			Name      string
			Path      string
			Templates Templates
		}
		Structure []interface{}
	}

	// StructureCreator interface defines the CreateStructure method.
	StructureCreator interface {
		CreateStructure(folders []string) error
	}

	// StructureCreatorFunc is an adapter type to allow use of ordinary functions as directory StructureCreator.
	StructureCreatorFunc func(folders []string) error

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

	// TemplateFile represents a template file.
	TemplateFile struct {
		Name   string
		Values map[string]interface{}
	}

	// Templates represents the template configs in the metadata section of the .fundi.yaml file.
	Templates struct {
		Path string
	}

	// TemplateParser interface defines ParseTemplates method.
	TemplateParser interface {
		ParseTemplates(data map[string]*TemplateFile, templatePath string) (map[string][]byte, error)
	}

	// TemplateParserFunc is an adapter type to allow use of ordinary functions as directory TemplateParser.
	TemplateParserFunc func(data map[string]*TemplateFile, templatePath string) (map[string][]byte, error)
)

// CreateStructure creates a directory structure.
func (fn StructureCreatorFunc) CreateStructure(folders []string) error {
	return fn(folders)
}

// Read wraps the reader function fn.
func (fn FundiFileReaderFunc) Read() (*FundiFile, error) {
	return fn()
}

// CreateFiles wraps the file creator function fn.
func (fn FileCreatorFunc) CreateFiles(files map[string][]byte) error {
	return fn(files)
}

// ParseTemplates wraps the template parser function fn.
func (fn TemplateParserFunc) ParseTemplates(
	data map[string]*TemplateFile,
	templatePath string,
) (map[string][]byte, error) {
	return fn(data, templatePath)
}

func generateHierarchy(root string, data interface{}) interface{} {
	switch reflect.ValueOf(data).Kind() {
	case reflect.Slice:
		dirs := cast.ToStringSlice(data)
		hierarchy := make([]string, 0, len(dirs))
		for _, dir := range dirs {
			hierarchy = append(hierarchy, root+string(os.PathSeparator)+cast.ToString(dir))
		}
		return hierarchy
	case reflect.Map:
		hierarchy := make(map[string][]byte)
		for name, byteData := range data.(map[string][]byte) {
			hierarchy[root+string(os.PathSeparator)+name] = byteData
		}
		return hierarchy
	}

	return nil
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

func getFilesAndTemplates(structure []interface{}) (map[string]*TemplateFile, error) {
	files := make(map[string]*TemplateFile)
	for _, item := range structure {
		kind := reflect.ValueOf(item).Kind()
		switch kind {
		case reflect.Map:
			dict := cast.ToStringMap(item)
			if isDirectory(cast.ToStringMap(item)) && hasContents(cast.ToStringMap(item)) {
				innerFiles, err := getFilesAndTemplates(cast.ToSlice(dict["contains"]))
				if err != nil {
					return innerFiles, err
				}

				parent := cast.ToString(dict["folder"])
				for name, tpl := range innerFiles {
					files[parent+string(os.PathSeparator)+name] = tpl
				}
			} else if isFile(cast.ToStringMap(item)) {
				files[cast.ToString(dict["file"])] = templateFile(cast.ToStringMap(dict["template"]))
			}
		default:
			return nil, errors.Errorf("unexpected kind: %s", kind)
		}
	}
	return files, nil
}

func templateFile(tpl map[string]interface{}) *TemplateFile {
	return &TemplateFile{
		Name:   cast.ToString(tpl["name"]),
		Values: cast.ToStringMap(tpl["values"]),
	}
}
