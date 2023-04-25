package generate

import (
	"context"

	"github.com/goava/di"
	"github.com/pkg/errors"
)

type (
	// DirectoryStructureUseCase type.
	DirectoryStructureUseCase struct {
		fundiFile        FundiFileReader
		structureCreator StructureCreator
	}

	// EmptyFilesUseCase type.
	EmptyFilesUseCase struct {
		fileReader  FundiFileReader
		fileCreator FileCreator
	}

	// FilesUseCase type.
	FilesUseCase struct {
		fileReader  FundiFileReader
		fileCreator FileCreator
		parser      TemplateParser2
	}

	// InitialiseUseCase type.
	InitialiseUseCase struct {
		skipTemplates             bool
		directoryStructureUseCase *DirectoryStructureUseCase
		emptyFilesUseCase         *EmptyFilesUseCase
		filesUseCase              *FilesUseCase
	}
)

// ProvideUseCases returns a DI container option with use case types.
func ProvideUseCases() di.Option {
	return di.Options(
		di.Provide(NewDirectoryStructureUseCase),
		di.Provide(NewEmptyFilesUseCase),
		di.Provide(NewFilesUseCase),
		di.Provide(NewInitialiseUseCase),
	)
}

// NewDirectoryStructureUseCase returns an instance of DirectoryStructureUseCase.
func NewDirectoryStructureUseCase(
	reader FundiFileReader,
	creator StructureCreator) *DirectoryStructureUseCase {
	return &DirectoryStructureUseCase{
		fundiFile:        reader,
		structureCreator: creator,
	}
}

// NewEmptyFilesUseCase returns an instance of EmptyFilesUseCase.
func NewEmptyFilesUseCase(reader FundiFileReader, creator FileCreator) *EmptyFilesUseCase {
	return &EmptyFilesUseCase{
		fileReader:  reader,
		fileCreator: creator,
	}
}

// NewFilesUseCase returns an instance of FilesUseCase.
func NewFilesUseCase(reader FundiFileReader, creator FileCreator, parser TemplateParser2) *FilesUseCase {
	return &FilesUseCase{
		fileReader:  reader,
		fileCreator: creator,
		parser:      parser,
	}
}

// NewInitialiseUseCase returns an instance of InitialiseUseCase.
func NewInitialiseUseCase(
	dsUseCase *DirectoryStructureUseCase,
	efUseCase *EmptyFilesUseCase,
	fUseCase *FilesUseCase,
) *InitialiseUseCase {
	return &InitialiseUseCase{
		directoryStructureUseCase: dsUseCase,
		emptyFilesUseCase:         efUseCase,
		filesUseCase:              fUseCase,
	}
}

// Execute generates an empty directory structure.
func (usecase *DirectoryStructureUseCase) Execute(ctx context.Context) error {
	fundiFile, err := usecase.fundiFile.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read fundi file")
	}

	directories, err := getAllDirectories(fundiFile.ProjectStructure())
	if err != nil {
		return errors.Wrap(err, "failed to get directories")
	}

	if err := usecase.structureCreator.CreateStructure(
		generateHierarchy(fundiFile.ProjectPath(), directories).([]string),
	); err != nil {
		return errors.Wrap(err, "failed to create directory hierarchy")
	}

	return nil
}

// Execute adds empty files to an existing directory structure.
func (usecase *EmptyFilesUseCase) Execute(ctx context.Context) error {
	fundiFile, err := usecase.fileReader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read fundi file")
	}

	files, err := getFilesSkipTemplates(fundiFile.ProjectStructure())
	if err != nil {
		return errors.Wrap(err, "failed to get files")
	}

	if err := usecase.fileCreator.CreateFiles(
		generateEmptyFiles(
			generateHierarchy(fundiFile.ProjectPath(), files).([]string),
		),
	); err != nil {
		return errors.Wrap(err, "failed to create empty files")
	}

	return nil
}

// Execute generates files from templates.
func (usecase *FilesUseCase) Execute(ctx context.Context) error {
	fundiFile, err := usecase.fileReader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read fundi file")
	}

	filesAndTemplates, err := getFilesAndTemplates(fundiFile.ProjectStructure())
	if err != nil {
		return errors.Wrap(err, "failed to get files and their templates")
	}

	parsedFiles, err := usecase.parser.ParseTemplates(filesAndTemplates, fundiFile.TemplatesPath())
	if err != nil {
		return errors.Wrap(err, "failed to parse templates")
	}

	if err := usecase.fileCreator.CreateFiles(
		generateHierarchy(fundiFile.ProjectPath(), parsedFiles).(map[string][]byte),
	); err != nil {
		return errors.Wrap(err, "failed to create files")
	}

	return nil
}

// WithSkipTemplates sets skipTemplates attribute of InitialiseUseCase type.
func (usecase *InitialiseUseCase) WithSkipTemplates(skip bool) *InitialiseUseCase {
	usecase.skipTemplates = skip

	return usecase
}

// Execute creates a directory structure and generates files.
func (usecase *InitialiseUseCase) Execute(ctx context.Context) error {
	if err := usecase.directoryStructureUseCase.Execute(ctx); err != nil {
		return errors.Wrap(err, "failed to initialise")
	}

	if usecase.skipTemplates {
		return usecase.emptyFilesUseCase.Execute(ctx)
	}

	return usecase.filesUseCase.Execute(ctx)
}
