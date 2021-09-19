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
		fundiFile        FundiFileReader
		structureCreator StructureCreator
	}

	// FilesSkipTemplates use case type.
	FilesSkipTemplates struct {
		fileReader FundiFileReader
		fCreator   FileCreator
	}

	// FilesFromTemplates use case type.
	FilesFromTemplates struct {
		fileReader FundiFileReader
		fCreator   FileCreator
		parser     TemplateParser
	}
)

// ProvideUseCases returns a DI container option with use case types.
func ProvideUseCases() di.Option {
	return di.Options(
		di.Provide(NewDirectoryStructure),
		di.Provide(NewFilesSkipTemplates),
		di.Provide(NewFilesFromTemplates),
	)
}

// NewDirectoryStructure returns an instance of DirectoryStructure use case.
func NewDirectoryStructure(
	reader FundiFileReader,
	creator StructureCreator) *DirectoryStructure {
	return &DirectoryStructure{
		fundiFile:        reader,
		structureCreator: creator,
	}
}

// NewFilesSkipTemplates returns an instance of FilesSkipTemplates use case.
func NewFilesSkipTemplates(reader FundiFileReader, creator FileCreator) *FilesSkipTemplates {
	return &FilesSkipTemplates{
		fileReader: reader,
		fCreator:   creator,
	}
}

// NewFilesFromTemplates returns an instance of FilesFromTemplates use case.
func NewFilesFromTemplates(reader FundiFileReader, creator FileCreator, parser TemplateParser) *FilesFromTemplates {
	return &FilesFromTemplates{
		fileReader: reader,
		fCreator:   creator,
		parser:     parser,
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

	if err := ps.structureCreator.CreateStructure(
		generateHierarchy(fundiFile.Metadata.Path, directories).([]string),
	); err != nil {
		return errors.Wrap(err, "failed to create directory hierarchy")
	}

	return nil
}

// UseCase to add empty files to an existing directory structure.
func (ef *FilesSkipTemplates) UseCase() error {
	fundiFile, err := ef.fileReader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read fundi file")
	}

	files, err := getFilesSkipTemplates(fundiFile.Structure)
	if err != nil {
		return errors.Wrap(err, "failed to get files")
	}

	if err := ef.fCreator.CreateFiles(
		generateEmptyFiles(
			generateHierarchy(fundiFile.Metadata.Path, files).([]string),
		),
	); err != nil {
		return errors.Wrap(err, "failed to create empty files")
	}

	return nil
}

// UseCase to generate files from templates.
func (f *FilesFromTemplates) UseCase() error {
	fundiFile, err := f.fileReader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read fundi file")
	}

	filesAndTemplates, err := getFilesAndTemplates(fundiFile.Structure)
	if err != nil {
		return errors.Wrap(err, "failed to get files and their templates")
	}

	parsedFiles, err := f.parser.ParseTemplates(filesAndTemplates, fundiFile.Metadata.Templates.Path)
	if err != nil {
		return errors.Wrap(err, "failed to parse templates")
	}

	if err := f.fCreator.CreateFiles(
		generateHierarchy(fundiFile.Metadata.Path, parsedFiles).(map[string][]byte),
	); err != nil {
		return errors.Wrap(err, "failed to create files")
	}

	return nil
}
