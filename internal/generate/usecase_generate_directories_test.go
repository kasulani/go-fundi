package generate

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestProjectDirectoryStructureUseCase(t *testing.T) {
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
			configFile: NewConfigurationFile(
				&Metadata{output: ".", templates: "./testdata"},
				Directories{
					&Directory{
						name: "project_root_directory",
						subDirectories: Directories{
							&Directory{
								name:  "cmd",
								files: Files{&File{name: "main.go", template: "main.go.tmpl"}},
							},
							&Directory{
								name: "internal",
								subDirectories: Directories{
									&Directory{
										name:           "domain",
										files:          Files{&File{name: "domain.go", template: "domain.go.tmpl"}},
										subDirectories: nil,
									},
								},
							},
						},
					},
				},
			),
		},
		"when the project structure is successfully created, return no error": {
			configFile: NewConfigurationFile(
				&Metadata{output: ".", templates: "./testdata"},
				Directories{
					&Directory{
						name: "project_root_directory",
						subDirectories: Directories{
							&Directory{
								name:  "cmd",
								files: Files{&File{name: "main.go", template: "main.go.tmpl"}},
							},
							&Directory{
								name: "internal",
								subDirectories: Directories{
									&Directory{
										name:           "domain",
										files:          Files{&File{name: "domain.go", template: "domain.go.tmpl"}},
										subDirectories: nil,
									},
								},
							},
						},
					},
				},
			),
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

			useCase := NewProjectDirectoryStructureUseCase(testCase.structureCreator)
			err := useCase.GenerateProjectStructure(context.Background(), testCase.configFile)

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
		expectedDirs     []string
		configFile       *ConfigurationFile
		structureCreator DirectoryStructureCreator
	}{
		"returns all files in the configurationFile": {
			expectedDirs: []string{
				"project_root_directory/cmd",
				"project_root_directory/internal/domain",
			},
			structureCreator: mockDirectoryStructureCreator(
				func(_ context.Context, _ *ProjectDirectoryStructure) error {
					return nil
				},
			),
			configFile: NewConfigurationFile(
				&Metadata{output: ".", templates: "./testdata"},
				Directories{
					&Directory{
						name: "project_root_directory",
						subDirectories: Directories{
							&Directory{
								name:  "cmd",
								files: Files{&File{name: "main.go", template: "main.go.tmpl"}},
							},
							&Directory{
								name: "internal",
								subDirectories: Directories{
									&Directory{
										name:           "domain",
										files:          Files{&File{name: "domain.go", template: "domain.go.tmpl"}},
										subDirectories: nil,
									},
								},
							},
						},
					},
				},
			),
		},
		"returns an empty list of files": {
			expectedDirs: make([]string, 0),
			structureCreator: mockDirectoryStructureCreator(
				func(_ context.Context, _ *ProjectDirectoryStructure) error {
					return nil
				},
			),
			configFile: NewConfigurationFile(&Metadata{output: ".", templates: "./testdata"}, Directories{}),
		},
	}

	for name, testCase := range tests {
		testCase := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			useCase := NewProjectDirectoryStructureUseCase(testCase.structureCreator)
			actualDirs := useCase.getAllDirectoriesInTheConfigFile(testCase.configFile.directories)

			assert.Equal(t, testCase.expectedDirs, actualDirs)
		})
	}
}
