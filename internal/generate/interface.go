package generate

import "context"

type (
	// DirectoryStructureCreator defines CreateDirectoryStructure.
	DirectoryStructureCreator interface {
		CreateDirectoryStructure(ctx context.Context, output string, directories []string) error
	}

	// FilesCreator interface defines CreateFiles.
	FilesCreator interface {
		CreateFiles(ctx context.Context, metadata *Metadata, files FileTemplates) error
	}
)
