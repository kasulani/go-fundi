package generate

import "context"

type (
	// DirectoryStructureCreator defines CreateDirectoryStructure.
	DirectoryStructureCreator interface {
		CreateDirectoryStructure(ctx context.Context, structure *ProjectDirectoryStructure) error
	}

	// FileCreator interface define the CreateFiles method.
	FileCreator interface {
		CreateFiles(files map[string][]byte) error
	}

	// FilesCreator interface defines CreateFiles.
	FilesCreator interface {
		CreateFiles(ctx context.Context, metadata *Metadata, files FileTemplates) error
	}
)
