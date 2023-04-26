package app

import "github.com/spf13/afero"

func newFileReader(fs afero.Fs) *fileReader {
	return &fileReader{fs: fs}
}

func newDirectoryCreator(fs afero.Fs) *directoryCreator {
	return &directoryCreator{fs: fs}
}

func newFilesCreator(fs afero.Fs) *filesCreator {
	return &filesCreator{fs: fs}
}
