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

	// EmptyFiles use case type.
	EmptyFiles struct {
		fileReader FundiFileReader
		fCreator   FileCreator
	}
)

// ProvideUseCases returns a DI container option with use case types.
func ProvideUseCases() di.Option {
	return di.Options(
		di.Provide(NewDirectoryStructure),
		di.Provide(NewEmptyFiles),
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

// NewEmptyFiles returns an instance of EmptyFiles.
func NewEmptyFiles(reader FundiFileReader, creator FileCreator) *EmptyFiles {
	return &EmptyFiles{
		fileReader: reader,
		fCreator:   creator,
	}
}

// UseCase to generate an empty directory structure.
func (ps *DirectoryStructure) UseCase() error {
	fundiFile, err := ps.fundiFile.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read fundi file")
	}

	directories, err := getAllDirectories(fundiFile.Structure)
	if err != nil {
		return errors.Wrap(err, "failed to get directories")
	}

	if err := ps.hCreator.CreateHierarchy(generateHierarchy(fundiFile.Metadata.Path, directories)); err != nil {
		return errors.Wrap(err, "failed to create directory hierarchy")
	}

	return nil
}

// UseCase to add empty files to an existing directory structure.
func (ef *EmptyFiles) UseCase() error {
	fundiFile, err := ef.fileReader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read fundi file")
	}

	files, err := getFilesSkipTemplates(fundiFile.Structure)
	if err != nil {
		return errors.Wrap(err, "failed to get directories")
	}

	if err := ef.fCreator.CreateFiles(
		generateEmptyFiles(generateHierarchy(fundiFile.Metadata.Path, files)),
	); err != nil {
		return errors.Wrap(err, "failed to add empty files to directory structure")
	}

	return nil
}
