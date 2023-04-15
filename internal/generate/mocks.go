package generate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spf13/afero"
)

type (
	mockDirectoryStructureCreator     func(ctx context.Context, structure *ProjectDirectoryStructure) error
	inMemoryDirectoryStructureCreator struct {
		test       *testing.T
		fileSystem afero.Fs
	}
)

// CreateDirectoryStructure is a mock.
func (m mockDirectoryStructureCreator) CreateDirectoryStructure(
	ctx context.Context,
	structure *ProjectDirectoryStructure,
) error {
	return m(ctx, structure)
}

// CreateDirectoryStructure is implemented by an in memory file system.
func (m *inMemoryDirectoryStructureCreator) CreateDirectoryStructure(
	_ context.Context,
	structure *ProjectDirectoryStructure,
) error {
	m.test.Helper()

	for _, dir := range structure.directories {
		m.test.Logf("creating directory hierarchy: %s...", dir)
		if err := m.fileSystem.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (m *inMemoryDirectoryStructureCreator) assertDirectoryStructureExists(dirs []string) {
	m.test.Helper()

	for _, dir := range dirs {
		info, err := m.fileSystem.Stat(dir)
		if err != nil {
			m.test.Fatalf("unexpected error: %s", err)
		}

		assert.True(m.test, info.IsDir())
	}
}
