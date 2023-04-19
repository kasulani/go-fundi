package generate

import (
	"context"
	"os"

	"github.com/pkg/errors"
)

type (
	// ProjectDirectoryStructureUseCase creates the project directory structure.
	ProjectDirectoryStructureUseCase struct {
		structureCreator DirectoryStructureCreator
	}
)

func (useCase *ProjectDirectoryStructureUseCase) getAllDirectoriesInTheConfigFile(directories Directories) []string {
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

// GenerateProjectStructure specified in the configuration file.
func (useCase *ProjectDirectoryStructureUseCase) GenerateProjectStructure(
	ctx context.Context,
	configFile *ConfigurationFile,
) error {
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
