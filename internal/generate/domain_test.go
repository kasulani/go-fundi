package generate

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestGenerateHierarchy(t *testing.T) {
	tests := map[string]struct {
		root string
		data interface{}
		want interface{}
	}{
		"data as a slice": {
			root: "./testing",
			data: []string{"cmd", "pkg"},
			want: []string{"./testing/cmd", "./testing/pkg"},
		},
		"data as a map": {
			root: "./testing",
			data: map[string][]byte{
				"app.go": []byte("app-file-data"),
				"doc.go": []byte("doc-file-data"),
			},
			want: map[string][]byte{
				"./testing/app.go": []byte("app-file-data"),
				"./testing/doc.go": []byte("doc-file-data"),
			},
		},
		"nil parameters": {
			root: "",
			data: nil,
			want: nil,
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual := generateHierarchy(tc.root, tc.data)
			assert.Equal(t, tc.want, actual)
		})
	}
}

func TestGetAllDirectories(t *testing.T) {
	tests := map[string]struct {
		expectedDirs []string
		hasError     bool
		targetErr    error
		reader       FundiFileReader
	}{
		"has valid structure": {
			expectedDirs: []string{
				"funditest/cmd",
				"funditest/docker-files",
				"funditest/docs",
				"funditest/features",
				"funditest/pkg/app",
				"funditest/pkg/behaviour",
				"funditest/pkg/domain",
			},
			hasError: false,
			reader:   FundiFileReaderFunc(reader(t)),
		},
		"has empty structure": {
			expectedDirs: nil,
			hasError:     false,
			reader: FundiFileReaderFunc(func() (*FundiFile, error) {
				t.Helper()

				return new(FundiFile), nil
			}),
		},
		"structure is a slice of strings": {
			expectedDirs: nil,
			hasError:     true,
			targetErr:    errors.New("unexpected kind: string"),
			reader: FundiFileReaderFunc(func() (*FundiFile, error) {
				t.Helper()
				cfg := new(FundiFile)
				cfg.Structure = []interface{}{"cmd", "pkg"}

				return cfg, nil
			}),
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := tc.reader.Read()
			checkError(t, err)

			actualDirs, err := getAllDirectories(cfg.Structure)
			assert.Equal(t, tc.expectedDirs, actualDirs)

			if tc.hasError {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.targetErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetFilesSkipTemplates(t *testing.T) {
	tests := map[string]struct {
		expectedFiles []string
		hasError      bool
		targetErr     error
		reader        FundiFileReader
	}{
		"has valid structure": {
			expectedFiles: []string{
				"funditest/docker-compose.yml",
				"funditest/README.md",
				"funditest/docs/index.html",
				"funditest/pkg/app/doc.go",
			},
			hasError: false,
			reader:   FundiFileReaderFunc(reader(t)),
		},
		"has empty structure": {
			expectedFiles: nil,
			hasError:      false,
			reader: FundiFileReaderFunc(func() (*FundiFile, error) {
				t.Helper()

				return new(FundiFile), nil
			}),
		},
		"structure is a slice of strings": {
			expectedFiles: nil,
			hasError:      true,
			targetErr:     errors.New("unexpected kind: string"),
			reader: FundiFileReaderFunc(func() (*FundiFile, error) {
				t.Helper()
				cfg := new(FundiFile)
				cfg.Structure = []interface{}{"docker-compose.yml", "README.md"}

				return cfg, nil
			}),
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := tc.reader.Read()
			checkError(t, err)

			actualFiles, err := getFilesSkipTemplates(cfg.Structure)
			assert.Equal(t, tc.expectedFiles, actualFiles)

			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateEmptyFiles(t *testing.T) {
	tests := map[string]struct {
		paths         []string
		expectedFiles map[string][]byte
	}{
		"paths provided": {
			paths: []string{
				"funditest/README.md",
				"funditest/docker-compose.yml",
			},
			expectedFiles: map[string][]byte{
				"funditest/README.md":          []byte(""),
				"funditest/docker-compose.yml": []byte(""),
			},
		},
		"no paths provided": {
			paths:         []string{},
			expectedFiles: make(map[string][]byte),
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actualFiles := generateEmptyFiles(tc.paths)
			assert.Equal(t, tc.expectedFiles, actualFiles)
		})
	}
}

func TestGetFilesAndTemplates(t *testing.T) {
	tests := map[string]struct {
		want      map[string]*TemplateFile
		hasError  bool
		targetErr error
		reader    FundiFileReader
	}{
		"has valid structure": {
			want: map[string]*TemplateFile{
				"funditest/docker-compose.yml": {
					Name:   "",
					Values: map[string]interface{}{},
				},
				"funditest/README.md": {
					Name:   "",
					Values: map[string]interface{}{},
				},
				"funditest/docs/index.html": {
					Name:   "",
					Values: map[string]interface{}{},
				},
				"funditest/pkg/app/doc.go": {
					Name: "doc.go.tmpl",
					Values: map[string]interface{}{
						"package": "app",
					},
				},
			},
			hasError: false,
			reader:   FundiFileReaderFunc(reader(t)),
		},
		"has empty structure": {
			want:     map[string]*TemplateFile{},
			hasError: false,
			reader: FundiFileReaderFunc(func() (*FundiFile, error) {
				t.Helper()

				return new(FundiFile), nil
			}),
		},
		"structure is a slice of strings": {
			want:      nil,
			hasError:  true,
			targetErr: errors.New("unexpected kind: string"),
			reader: FundiFileReaderFunc(func() (*FundiFile, error) {
				t.Helper()
				cfg := new(FundiFile)
				cfg.Structure = []interface{}{"docker-compose.yml", "README.md"}

				return cfg, nil
			}),
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := tc.reader.Read()
			checkError(t, err)

			actual, err := getFilesAndTemplates(cfg.Structure)
			assert.Equal(t, tc.want, actual)

			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
