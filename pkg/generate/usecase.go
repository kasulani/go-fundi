package generate

import (
	"github.com/goava/di"
	"github.com/pkg/errors"
)

type (
	// Generator interface defines the UseCase method.
	Generator interface {
		UseCase() error
	}

	// DirectoryStructure use case type.
	DirectoryStructure struct {
		fundiFile FundiFileReader
		hCreator  HierarchyCreator
	}
)

// ProvideUseCases returns a DI container option with use case types.
func ProvideUseCases() di.Option {
	return di.Options(
		di.Provide(NewDirectoryStructure),
	)
}

// NewDirectoryStructure returns an instance of DirectoryStructure.
func NewDirectoryStructure(
	reader FundiFileReader,
	creator HierarchyCreator) *DirectoryStructure {
	return &DirectoryStructure{
		fundiFile: reader,
		hCreator:  creator,
	}
}

// UseCase to generate an empty directory structure.
func (ps *DirectoryStructure) UseCase() error {
	project, err := ps.fundiFile.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read fundi file")
	}

	directories, err := getAllDirectories(project.Structure)
	if err != nil {
		return errors.Wrap(err, "failed to get directories")
	}

	if err := ps.hCreator.CreateHierarchy(generateHierarchy(project.Metadata.Path, directories)); err != nil {
		return errors.Wrap(err, "failed to create directory hierarchy")
	}

	return nil
}
