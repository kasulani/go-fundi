package generate

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestScaffoldProject(t *testing.T) {
	fs := afero.NewMemMapFs()

	tests := map[string]struct {
		expectedErr       error
		configFile        *ConfigurationFile
		structureCreator  DirectoryStructureCreator
		fileCreator       FilesCreator
		expectedStructure []string
		expectedFiles     []string
	}{
		"when the directory structure creator fails, return an error": {
			expectedErr: errors.New("failed to create project directory structure: an-OS-error"),
			structureCreator: mockDirectoryStructureCreator(
				func(ctx context.Context, output string, directories []string) error {
					return errors.New("an-OS-error")
				},
			),
			configFile: NewTestConfigurationFile(),
		},
		"when the file creator fails, return an error": {
			expectedErr: errors.New("failed to create project files: an-OS-error"),
			structureCreator: mockDirectoryStructureCreator(
				func(ctx context.Context, output string, directories []string) error {
					return nil
				},
			),
			fileCreator: mockFilesCreator(func(_ context.Context, _ *Metadata, _ FileTemplates) error {
				return errors.New("an-OS-error")
			}),
			configFile: NewTestConfigurationFile(),
		},
		"when scaffolding is successful, return no error": {
			configFile: NewTestConfigurationFile(),
			expectedStructure: []string{
				"./project_root_directory/cmd",
				"./project_root_directory/internal/domain",
			},
			expectedFiles: []string{
				"./project_root_directory/README.md",
				"./project_root_directory/cmd/main.go",
				"./project_root_directory/internal/domain/domain.go",
			},
			structureCreator: &inMemoryDirectoryStructureCreator{test: t, fileSystem: fs},
			fileCreator:      &inMemoryFilesCreator{test: t, fileSystem: fs},
		},
	}

	for name, testCase := range tests {
		testCase := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			useCase := NewProjectUseCase(testCase.structureCreator, testCase.fileCreator)
			err := useCase.ScaffoldProject(context.Background(), testCase.configFile)

			switch testCase.expectedErr != nil {
			case true:
				assert.EqualError(t, err, testCase.expectedErr.Error())
			case false:
				assert.NoError(t, err)
				testCase.structureCreator.(*inMemoryDirectoryStructureCreator).assertDirectoryStructureExists(testCase.expectedStructure)
				testCase.fileCreator.(*inMemoryFilesCreator).assertCreatedFiles(testCase.expectedFiles)
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
				func(ctx context.Context, output string, directories []string) error {
					return nil
				},
			), nil)
			actualDirs := useCase.getAllDirectoriesInTheConfigFile(testCase.configFile.directories)

			assert.Equal(t, testCase.expectedDirs, actualDirs)
		})
	}
}

func TestGetFilesAndTemplates(t *testing.T) {
	tests := map[string]struct {
		expectedFileTemplates FileTemplates
		configFile            *ConfigurationFile
	}{
		"returns all files templates in the configurationFile": {
			expectedFileTemplates: FileTemplates{
				"project_root_directory/README.md":                 "README.md.tmpl",
				"project_root_directory/cmd/main.go":               "main.go.tmpl",
				"project_root_directory/internal/domain/domain.go": "domain.go.tmpl",
			},
			configFile: NewTestConfigurationFile(),
		},
		"returns an empty list of file templates": {
			expectedFileTemplates: FileTemplates{},
			configFile:            NewConfigurationFile(&Metadata{output: ".", templates: "./testdata"}, Directories{}),
		},
	}

	for name, testCase := range tests {
		testCase := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, testCase.expectedFileTemplates, testCase.configFile.getFilesAndTemplates())
		})
	}
}
