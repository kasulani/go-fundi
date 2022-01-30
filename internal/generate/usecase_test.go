package generate

import (
	"context"
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
		hCreator          StructureCreator
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
				cfg.structure = []interface{}{"cmd", "pkg"}

				return cfg, nil
			}),
			hasError: true,
		},
		"CreateStructure fails": {
			reader: FundiFileReaderFunc(reader(t)),
			hCreator: StructureCreatorFunc(func(hierarchy []string) error {
				return errors.New("failed to create hierarchy")
			}),
			hasError: true,
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			generateStructure := DirectoryStructureUseCase{
				fundiFile:        tc.reader,
				structureCreator: tc.hCreator,
			}

			err := generateStructure.Execute(context.Background())
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

		return NewFundiFile(
			file.Metadata.Path,
			NewTemplates(file.Metadata.Templates.Path),
			file.Structure,
		), nil
	}
}

func (hc *mockHierarchyCreator) assertHierarchy(dirs []string) {
	hc.test.Helper()

	for _, dir := range dirs {
		info, err := hc.fileSystem.Stat(dir)
		checkError(hc.test, err)
		assert.True(hc.test, info.IsDir())
	}
}

func (hc *mockHierarchyCreator) CreateStructure(folders []string) error {
	hc.test.Helper()

	for _, h := range folders {
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
				cfg.structure = []interface{}{"docker-compose.yml", "README.md"}

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

			skipTemplates := EmptyFilesUseCase{
				fileReader:  tc.reader,
				fileCreator: tc.fCreator,
			}

			err := skipTemplates.Execute(context.Background())

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
				cfg.structure = []interface{}{"docker-compose.yml", "README.md"}

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

			filesFromTemplates := FilesUseCase{
				fileReader:  tc.reader,
				fileCreator: tc.fCreator,
				parser:      tc.parser,
			}

			err := filesFromTemplates.Execute(context.Background())

			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tc.fCreator.(*mockFileCreator).assertCreatedFiles(tc.want)
			}
		})
	}
}

func TestInitialiseUseCase(t *testing.T) {
	tests := map[string]struct {
		skipTemplates    bool
		fileReader       FundiFileReader
		fileCreator      FileCreator
		parser           TemplateParser
		structureCreator StructureCreator
		expectedError    error
	}{
		"returns no error when initialisation is successful ": {
			skipTemplates: true,
			fileReader:    FundiFileReaderFunc(reader(t)),
			structureCreator: StructureCreatorFunc(func(folders []string) error {
				return nil
			}),
			fileCreator: FileCreatorFunc(func(files map[string][]byte) error {
				return nil
			}),
		},
		"returns an error when creating the directory structure fails": {
			fileReader: FundiFileReaderFunc(func() (*FundiFile, error) {
				return &FundiFile{}, nil
			}),
			structureCreator: StructureCreatorFunc(func(folders []string) error {
				return errors.New("failed to create directory structure")
			}),
			expectedError: errors.New(
				"failed to initialise: failed to create directory hierarchy:" +
					" failed to create directory structure",
			),
		},
		"returns an error when creating empty files fails": {
			skipTemplates: true,
			fileReader: FundiFileReaderFunc(func() (*FundiFile, error) {
				return &FundiFile{}, nil
			}),
			structureCreator: StructureCreatorFunc(func(folders []string) error {
				return nil
			}),
			fileCreator: FileCreatorFunc(func(files map[string][]byte) error {
				return errors.New("failed to create files")
			}),
			expectedError: errors.New("failed to create empty files: failed to create files"),
		},
		"returns an error when creating files from templates fails": {
			fileReader: FundiFileReaderFunc(reader(t)),
			structureCreator: StructureCreatorFunc(func(folders []string) error {
				return nil
			}),
			fileCreator: FileCreatorFunc(func(files map[string][]byte) error {
				return nil
			}),
			parser: TemplateParserFunc(
				func(data map[string]*TemplateFile, templatePath string) (map[string][]byte, error) {
					return nil, errors.New("file error")
				},
			),
			expectedError: errors.New("failed to parse templates: file error"),
		},
	}

	for name, test := range tests {
		testcase := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			usecase := NewInitialiseUseCase(
				NewDirectoryStructureUseCase(testcase.fileReader, testcase.structureCreator),
				NewEmptyFilesUseCase(testcase.fileReader, testcase.fileCreator),
				NewFilesUseCase(testcase.fileReader, testcase.fileCreator, testcase.parser),
			)

			err := usecase.WithSkipTemplates(testcase.skipTemplates).Execute(context.Background())

			switch testcase.expectedError != nil {
			case true:
				assert.Error(t, err)
				assert.EqualError(t, err, testcase.expectedError.Error())
			case false:
				assert.NoError(t, err)
			}
		})
	}
}
