package generate

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

type (
	// ProjectUseCase creates the project directory structure and files.
	ProjectUseCase struct {
		structureCreator DirectoryStructureCreator
		fileCreator      FileCreator
	}
)

const (
	// DirectoriesOnly is an output selector.
	DirectoriesOnly = "directories-only"
	// EmptyFiles is an output selector.
	EmptyFiles = "empty-files"
	// All is an output selector.
	All = "all"
)

func (useCase *ProjectUseCase) getAllDirectoriesInTheConfigFile(directories Directories) []string {
	dirs := make([]string, 0)

	for _, directory := range directories {
		if !directory.hasSubDirectories() {
			dirs = append(dirs, directory.name)
			continue
		}

		subDirectories := useCase.getAllDirectoriesInTheConfigFile(directory.subDirectories)

		for _, dir := range subDirectories {
			dirs = append(dirs, directory.name+string(os.PathSeparator)+dir)
		}
	}

	return dirs
}

func (useCase *ProjectUseCase) getAllFilesInTheConfigFile(directories Directories) []string {
	files := make([]string, 0)

	for _, directory := range directories {
		for _, file := range directory.files {
			files = append(files, directory.name+string(os.PathSeparator)+file.name)
		}

		if directory.hasSubDirectories() {
			otherFiles := useCase.getAllFilesInTheConfigFile(directory.subDirectories)
			for _, file := range otherFiles {
				files = append(files, directory.name+string(os.PathSeparator)+file)
			}
		}
	}

	return files
}

func (useCase *ProjectUseCase) generateProjectStructure(ctx context.Context, configFile *ConfigurationFile) error {
	err := useCase.structureCreator.CreateDirectoryStructure(
		ctx,
		&ProjectDirectoryStructure{
			output:      configFile.metadata.output,
			directories: useCase.getAllDirectoriesInTheConfigFile(configFile.directories),
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to create project directory structure")
	}

	return nil
}

func (useCase *ProjectUseCase) generateEmptyFiles(_ context.Context, configFile *ConfigurationFile) error {
	files := make(map[string][]byte)

	allFiles := useCase.getAllFilesInTheConfigFile(configFile.directories)
	for _, file := range allFiles {
		files[file] = []byte("")
	}

	if err := useCase.fileCreator.CreateFiles(files); err != nil {
		return errors.Wrap(err, "failed to create project files")
	}

	return nil
}

// ScaffoldProject using the provided ConfigurationFile.
func (useCase *ProjectUseCase) ScaffoldProject(
	ctx context.Context,
	output string,
	configFile *ConfigurationFile,
) error {
	switch output {
	case DirectoriesOnly:
		return useCase.generateProjectStructure(ctx, configFile)
	case EmptyFiles:
		if err := useCase.generateProjectStructure(ctx, configFile); err != nil {
			return err
		}

		return useCase.generateEmptyFiles(ctx, configFile)
	case All:
		return errors.New("not-implemented")
	default:
		return fmt.Errorf("unknown output selector %s", output)
	}
}
