package generate

import (
	"io/ioutil"
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type (
	testFile struct {
		Version  int `yaml:"version"`
		Metadata struct {
			Name string `yaml:"name"`
			Path string `yaml:"path"`
		} `yaml:"metadata"`
		Structure []interface{} `yaml:"structure"`
	}

	mockHierarchyCreator struct {
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
		data, err := ioutil.ReadFile(".test.fundi.yaml")
		checkError(t, err)

		file := new(testFile)

		err = yaml.Unmarshal(data, file)
		if err != nil {
			return nil, err
		}

		return file.toFundiFile(), nil
	}
}

func (ts *testFile) toFundiFile() *FundiFile {
	ff := new(FundiFile)

	ff.Metadata.Name = ts.Metadata.Name
	ff.Metadata.Path = ts.Metadata.Path
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
