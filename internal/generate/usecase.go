package generate

import (
	"github.com/goava/di"
	"github.com/pkg/errors"
)

type (
	// UseCase interface defines the Execute method on a use case.
	UseCase interface {
		Execute() error
	}

	// DirectoryStructureUseCase use case type.
	DirectoryStructureUseCase struct {
		fundiFile        FundiFileReader
		structureCreator StructureCreator
	}

	// EmptyFilesUseCase type.
	EmptyFilesUseCase struct {
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
		di.Provide(NewDirectoryStructureUseCase),
		di.Provide(NewEmptyFilesUseCase),
		di.Provide(NewFilesFromTemplates),
	)
}

// NewDirectoryStructureUseCase returns an instance of DirectoryStructureUseCase use case.
func NewDirectoryStructureUseCase(
	reader FundiFileReader,
	creator StructureCreator) *DirectoryStructureUseCase {
	return &DirectoryStructureUseCase{
		fundiFile:        reader,
		structureCreator: creator,
	}
}

// NewEmptyFilesUseCase returns an instance of EmptyFilesUseCase use case.
func NewEmptyFilesUseCase(reader FundiFileReader, creator FileCreator) *EmptyFilesUseCase {
	return &EmptyFilesUseCase{
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

// Execute generates an empty directory structure.
func (ps *DirectoryStructureUseCase) Execute() error {
	fundiFile, err := ps.fundiFile.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read fundi file")
	}

	directories, err := getAllDirectories(fundiFile.ProjectStructure())
	if err != nil {
		return errors.Wrap(err, "failed to get directories")
	}

	if err := ps.structureCreator.CreateStructure(
		generateHierarchy(fundiFile.ProjectPath(), directories).([]string),
	); err != nil {
		return errors.Wrap(err, "failed to create directory hierarchy")
	}

	return nil
}

// Execute adds empty files to an existing directory structure.
func (ef *EmptyFilesUseCase) Execute() error {
	fundiFile, err := ef.fileReader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read fundi file")
	}

	files, err := getFilesSkipTemplates(fundiFile.ProjectStructure())
	if err != nil {
		return errors.Wrap(err, "failed to get files")
	}

	if err := ef.fCreator.CreateFiles(
		generateEmptyFiles(
			generateHierarchy(fundiFile.ProjectPath(), files).([]string),
		),
	); err != nil {
		return errors.Wrap(err, "failed to create empty files")
	}

	return nil
}

// Execute generates files from templates.
func (f *FilesFromTemplates) Execute() error {
	fundiFile, err := f.fileReader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read fundi file")
	}

	filesAndTemplates, err := getFilesAndTemplates(fundiFile.ProjectStructure())
	if err != nil {
		return errors.Wrap(err, "failed to get files and their templates")
	}

	parsedFiles, err := f.parser.ParseTemplates(filesAndTemplates, fundiFile.TemplatesPath())
	if err != nil {
		return errors.Wrap(err, "failed to parse templates")
	}

	if err := f.fCreator.CreateFiles(
		generateHierarchy(fundiFile.ProjectPath(), parsedFiles).(map[string][]byte),
	); err != nil {
		return errors.Wrap(err, "failed to create files")
	}

	return nil
}
