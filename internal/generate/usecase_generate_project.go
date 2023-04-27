package generate

import (
	"context"
	"os"

	"github.com/pkg/errors"
)

type (
	// ProjectUseCase creates the project directory structure and files.
	ProjectUseCase struct {
		structureCreator DirectoryStructureCreator
		filesCreator     FilesCreator
	}
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

func (useCase *ProjectUseCase) generateFilesFromTemplates(ctx context.Context, configFile *ConfigurationFile) error {
	if err := useCase.filesCreator.CreateFiles(
		ctx,
		configFile.metadata,
		configFile.getFilesAndTemplates(),
	); err != nil {
		return errors.Wrap(err, "failed to create project files")
	}

	return nil
}

// ScaffoldProject using the provided ConfigurationFile.
func (useCase *ProjectUseCase) ScaffoldProject(ctx context.Context, configFile *ConfigurationFile) error {
	if err := useCase.generateProjectStructure(ctx, configFile); err != nil {
		return err
	}
	if err := useCase.generateFilesFromTemplates(ctx, configFile); err != nil {
		return err
	}

	return nil
}
