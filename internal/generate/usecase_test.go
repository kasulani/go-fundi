package generate

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type (
	testConfigFile struct {
		Metadata struct {
			Name      string `yaml:"name"`
			Path      string `yaml:"path"`
			Templates struct {
				Path string `yaml:"path"`
			} `yaml:"templates"`
		} `yaml:"metadata"`
		Structure []interface{} `yaml:"structure"`
	}

	mockHierarchyCreator struct {
		test       *testing.T
		fileSystem afero.Fs
	}
	mockFileCreator struct {
		test       *testing.T
		fileSystem afero.Fs
	}
)

func TestProjectStructure_UseCase(t *testing.T) {
	tests := map[string]struct {
		reader            FundiFileReader
		hCreator          HierarchyCreator
		hasError          bool
		expectedHierarchy []string
	}{
		"happy path": {
			reader: FundiFileReaderFunc(reader(t)),
			hCreator: &mockHierarchyCreator{
				test:       t,
				fileSystem: afero.NewMemMapFs(),
			},
			hasError: false,
			expectedHierarchy: []string{
				"./funditest/cmd",
				"./funditest/docker-files",
				"./funditest/docs",
				"./funditest/features",
				"./funditest/pkg/app",
				"./funditest/pkg/behaviour",
				"./funditest/pkg/domain",
			},
		},
		"fundi file reader fails": {
			reader: FundiFileReaderFunc(func() (*FundiFile, error) {
				t.Helper()

				return nil, errors.New("read-error")
			}),
			hasError: true,
		},
		"getAllDirectories fails": {
			reader: FundiFileReaderFunc(func() (*FundiFile, error) {
				t.Helper()
				cfg := new(FundiFile)
				cfg.Structure = []interface{}{"cmd", "pkg"}

				return cfg, nil
			}),
			hasError: true,
		},
		"CreateHierarchy fails": {
			reader: FundiFileReaderFunc(reader(t)),
			hCreator: HierarchyCreatorFunc(func(hierarchy []string) error {
				return errors.New("failed to create hierarchy")
			}),
			hasError: true,
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			generateStructure := DirectoryStructure{
				fundiFile: tc.reader,
				hCreator:  tc.hCreator,
			}

			err := generateStructure.UseCase()
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tc.hCreator.(*mockHierarchyCreator).assertHierarchy(tc.expectedHierarchy)
			}
		})
	}
}

func reader(t *testing.T) func() (*FundiFile, error) {
	t.Helper()
	return func() (*FundiFile, error) {
		data, err := ioutil.ReadFile("../../testdata/.test.fundi.yaml")
		checkError(t, err)

		file := new(testConfigFile)

		err = yaml.Unmarshal(data, file)
		if err != nil {
			return nil, err
		}

		return file.toFundiFile(), nil
	}
}

func (ts *testConfigFile) toFundiFile() *FundiFile {
	ff := new(FundiFile)

	ff.Metadata.Name = ts.Metadata.Name
	ff.Metadata.Path = ts.Metadata.Path
	ff.Metadata.Templates.Path = ts.Metadata.Templates.Path
	ff.Structure = ts.Structure

	return ff
}

func (hc *mockHierarchyCreator) assertHierarchy(dirs []string) {
	hc.test.Helper()

	for _, dir := range dirs {
		info, err := hc.fileSystem.Stat(dir)
		checkError(hc.test, err)
		assert.True(hc.test, info.IsDir())
	}
}

func (hc *mockHierarchyCreator) CreateHierarchy(hierarchy []string) error {
	hc.test.Helper()

	for _, h := range hierarchy {
		hc.test.Logf("creating directory hierarchy: %s...", h)
		if err := hc.fileSystem.MkdirAll(h, 0755); err != nil {
			return err
		}
	}

	return nil
}

func TestEmptyFiles_UseCase(t *testing.T) {
	tests := map[string]struct {
		reader        FundiFileReader
		fCreator      FileCreator
		hasError      bool
		expectedFiles []string
	}{
		"happy path": {
			reader: FundiFileReaderFunc(reader(t)),
			fCreator: &mockFileCreator{
				test:       t,
				fileSystem: afero.NewMemMapFs(),
			},
			hasError: false,
			expectedFiles: []string{
				"funditest/docker-compose.yml",
				"funditest/README.md",
				"funditest/docs/index.html",
				"funditest/pkg/app/doc.go",
			},
		},
		"fundi file reader fails": {
			reader: FundiFileReaderFunc(func() (*FundiFile, error) {
				t.Helper()

				return nil, errors.New("read-error")
			}),
			hasError: true,
		},
		"getFilesSkipTemplates fails": {
			reader: FundiFileReaderFunc(func() (*FundiFile, error) {
				t.Helper()
				cfg := new(FundiFile)
				cfg.Structure = []interface{}{"docker-compose.yml", "README.md"}

				return cfg, nil
			}),
			hasError: true,
		},
		"CreateFiles fails": {
			reader: FundiFileReaderFunc(reader(t)),
			fCreator: FileCreatorFunc(func(files map[string][]byte) error {
				return errors.New("failed to add files to directory structure")
			}),
			hasError: true,
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			emptyFiles := EmptyFiles{
				fileReader: tc.reader,
				fCreator:   tc.fCreator,
			}

			err := emptyFiles.UseCase()

			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tc.fCreator.(*mockFileCreator).assertCreatedFiles(tc.expectedFiles)
			}
		})
	}
}

func (mf *mockFileCreator) CreateFiles(files map[string][]byte) error {
	mf.test.Helper()

	for name, data := range files {
		mf.test.Logf("creating file: %s...", name)

		if err := afero.WriteFile(mf.fileSystem, name, data, 0644); err != nil {
			return err
		}
	}

	return nil
}

func (mf *mockFileCreator) assertCreatedFiles(filenames []string) {
	mf.test.Helper()

	for _, name := range filenames {
		info, err := mf.fileSystem.Stat(name)
		assert.False(mf.test, info.IsDir())
		assert.False(mf.test, os.IsNotExist(err))
	}
}

func TestNewFilesFromTemplates(t *testing.T) {
	tests := map[string]struct {
		reader   FundiFileReader
		fCreator FileCreator
		parser   TemplateParser
		hasError bool
		want     []string
	}{
		"happy path": {
			reader: FundiFileReaderFunc(reader(t)),
			fCreator: &mockFileCreator{
				test:       t,
				fileSystem: afero.NewMemMapFs(),
			},
			parser: TemplateParserFunc(func(data map[string]*TemplateFile, templatePath string) (map[string][]byte, error) {
				parsedFiles := make(map[string][]byte)
				for fileName, _ := range data {
					parsedFiles[fileName] = []byte("test")
				}

				return parsedFiles, nil
			}),
			hasError: false,
			want: []string{
				"funditest/docker-compose.yml",
				"funditest/README.md",
				"funditest/docs/index.html",
				"funditest/pkg/app/doc.go",
			},
		},
		"fundi file reader fails": {
			reader: FundiFileReaderFunc(func() (*FundiFile, error) {
				t.Helper()

				return nil, errors.New("read-error")
			}),
			hasError: true,
		},
		"getFilesAndTemplates fails": {
			reader: FundiFileReaderFunc(func() (*FundiFile, error) {
				t.Helper()
				cfg := new(FundiFile)
				cfg.Structure = []interface{}{"docker-compose.yml", "README.md"}

				return cfg, nil
			}),
			hasError: true,
		},
		"ParseTemplates fails": {
			reader: FundiFileReaderFunc(reader(t)),
			parser: TemplateParserFunc(func(data map[string]*TemplateFile, templatePath string) (map[string][]byte, error) {
				return nil, errors.New("parse-error")
			}),
			hasError: true,
		},
		"CreateFiles fails": {
			reader: FundiFileReaderFunc(reader(t)),
			fCreator: FileCreatorFunc(func(files map[string][]byte) error {
				return errors.New("failed to add files to directory structure")
			}),
			parser: TemplateParserFunc(func(data map[string]*TemplateFile, templatePath string) (map[string][]byte, error) {
				parsedFiles := make(map[string][]byte)
				for fileName, _ := range data {
					parsedFiles[fileName] = []byte("test")
				}

				return parsedFiles, nil
			}),
			hasError: true,
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			filesFromTemplates := FilesFromTemplates{
				fileReader: tc.reader,
				fCreator:   tc.fCreator,
				parser:     tc.parser,
			}

			err := filesFromTemplates.UseCase()

			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tc.fCreator.(*mockFileCreator).assertCreatedFiles(tc.want)
			}
		})
	}
}
