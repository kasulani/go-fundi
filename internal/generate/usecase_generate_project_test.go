package generate

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGenerateDirectoriesOnly(t *testing.T) {
	tests := map[string]struct {
		expectedErr       error
		configFile        *ConfigurationFile
		structureCreator  DirectoryStructureCreator
		expectedStructure []string
	}{
		"when the directory structure creator fails, return an error": {
			expectedErr: errors.New("failed to create project directory structure: an-OS-error"),
			structureCreator: mockDirectoryStructureCreator(
				func(ctx context.Context, structure *ProjectDirectoryStructure) error {
					return errors.New("an-OS-error")
				},
			),
			configFile: NewTestConfigurationFile(),
		},
		"when the project directories are generated successfully, return no error": {
			configFile: NewTestConfigurationFile(),
			expectedStructure: []string{
				"./project_root_directory/cmd",
				"./project_root_directory/internal/domain",
			},
			structureCreator: &inMemoryDirectoryStructureCreator{test: t, fileSystem: afero.NewMemMapFs()},
		},
	}

	for name, testCase := range tests {
		testCase := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			useCase := NewProjectUseCase(testCase.structureCreator, nil)
			err := useCase.ScaffoldProject(context.Background(), DirectoriesOnly, testCase.configFile)

			switch testCase.expectedErr != nil {
			case true:
				assert.EqualError(t, err, testCase.expectedErr.Error())
			case false:
				assert.NoError(t, err)
				testCase.structureCreator.(*inMemoryDirectoryStructureCreator).assertDirectoryStructureExists(testCase.expectedStructure)
			}
		})
	}
}

func TestGetAllDirectoriesInTheConfigFile(t *testing.T) {
	tests := map[string]struct {
		expectedDirs []string
		configFile   *ConfigurationFile
	}{
		"returns all directories in the configurationFile": {
			expectedDirs: []string{
				"project_root_directory/cmd",
				"project_root_directory/internal/domain",
			},
			configFile: NewTestConfigurationFile(),
		},
		"returns an empty list of directories": {
			expectedDirs: make([]string, 0),
			configFile:   NewConfigurationFile(&Metadata{output: ".", templates: "./testdata"}, Directories{}),
		},
	}

	for name, testCase := range tests {
		testCase := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			useCase := NewProjectUseCase(mockDirectoryStructureCreator(
				func(_ context.Context, _ *ProjectDirectoryStructure) error {
					return nil
				},
			), nil)
			actualDirs := useCase.getAllDirectoriesInTheConfigFile(testCase.configFile.directories)

			assert.Equal(t, testCase.expectedDirs, actualDirs)
		})
	}
}

func TestGenerateEmptyFiles(t *testing.T) {
	fs := afero.NewMemMapFs()

	tests := map[string]struct {
		expectedErr      error
		configFile       *ConfigurationFile
		structureCreator DirectoryStructureCreator
		fileCreator      FileCreator
		expectedFiles    []string
	}{
		"when the file creator fails, return an error": {
			expectedErr: errors.New("failed to create project files: an-OS-error"),
			structureCreator: mockDirectoryStructureCreator(
				func(ctx context.Context, structure *ProjectDirectoryStructure) error {
					return nil
				},
			),
			fileCreator: mockFileCreator(func(files map[string][]byte) error {
				return errors.New("an-OS-error")
			}),
			configFile: NewTestConfigurationFile(),
		},
		"when the project directories are generated successfully, return no error": {
			configFile: NewTestConfigurationFile(),
			expectedFiles: []string{
				"./project_root_directory/README.md",
				"./project_root_directory/cmd/main.go",
				"./project_root_directory/internal/domain/domain.go",
			},
			structureCreator: &inMemoryDirectoryStructureCreator{test: t, fileSystem: fs},
			fileCreator:      &inMemoryFileCreator{test: t, fileSystem: fs},
		},
	}

	for name, testCase := range tests {
		testCase := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			useCase := NewProjectUseCase(testCase.structureCreator, testCase.fileCreator)
			err := useCase.ScaffoldProject(context.Background(), EmptyFiles, testCase.configFile)

			switch testCase.expectedErr != nil {
			case true:
				assert.EqualError(t, err, testCase.expectedErr.Error())
			case false:
				assert.NoError(t, err)
				testCase.fileCreator.(*inMemoryFileCreator).assertCreatedFiles(testCase.expectedFiles)
			}
		})
	}
}

func TestGetAllFilesInTheConfigFile(t *testing.T) {
	tests := map[string]struct {
		expectedFiles []string
		configFile    *ConfigurationFile
	}{
		"returns all files in the configurationFile": {
			expectedFiles: []string{
				"project_root_directory/README.md",
				"project_root_directory/cmd/main.go",
				"project_root_directory/internal/domain/domain.go",
			},
			configFile: NewTestConfigurationFile(),
		},
		"returns an empty list of files": {
			expectedFiles: make([]string, 0),
			configFile:    NewConfigurationFile(&Metadata{output: ".", templates: "./testdata"}, Directories{}),
		},
	}

	for name, testCase := range tests {
		testCase := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			useCase := NewProjectUseCase(mockDirectoryStructureCreator(
				func(_ context.Context, _ *ProjectDirectoryStructure) error {
					return nil
				},
			), nil)
			actualDirs := useCase.getAllFilesInTheConfigFile(testCase.configFile.directories)

			assert.Equal(t, testCase.expectedFiles, actualDirs)
		})
	}
}
