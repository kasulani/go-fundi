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
		root              string
		directories       []string
		expectedHierarchy []string
	}{
		"non nil parameters": {
			root:              "./testing",
			directories:       []string{"cmd", "pkg"},
			expectedHierarchy: []string{"./testing/cmd", "./testing/pkg"},
		},
		"nil parameters": {
			root:              "",
			directories:       []string{},
			expectedHierarchy: []string{},
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actualHierarchy := generateHierarchy(tc.root, tc.directories)
			assert.Equal(t, tc.expectedHierarchy, actualHierarchy)
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
