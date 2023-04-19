package generate

import "context"

type (
	// DirectoryStructureCreator defines CreateDirectoryStructure.
	DirectoryStructureCreator interface {
		CreateDirectoryStructure(ctx context.Context, structure *ProjectDirectoryStructure) error
	}
)
