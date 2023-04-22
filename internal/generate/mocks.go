package generate

import (
	"context"
	"os"
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
	mockFileCreator     func(files map[string][]byte) error
	inMemoryFileCreator struct {
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

// CreateFiles is a mock.
func (m mockFileCreator) CreateFiles(files map[string][]byte) error {
	return m(files)
}

func (mf *inMemoryFileCreator) CreateFiles(files map[string][]byte) error {
	mf.test.Helper()

	for name, data := range files {
		mf.test.Logf("creating file: %s...", name)

		if err := afero.WriteFile(mf.fileSystem, name, data, 0644); err != nil {
			return err
		}
	}

	return nil
}

func (mf *inMemoryFileCreator) assertCreatedFiles(filenames []string) {
	mf.test.Helper()

	for _, name := range filenames {
		info, err := mf.fileSystem.Stat(name)
		assert.False(mf.test, info.IsDir())
		assert.False(mf.test, os.IsNotExist(err))
	}
}
