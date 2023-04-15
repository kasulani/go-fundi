package generate

// NewConfigurationFile returns an instance of ConfigurationFile.
func NewConfigurationFile(metadata *Metadata, directories Directories) *ConfigurationFile {
	return &ConfigurationFile{metadata: metadata, directories: directories}
}

// NewProjectDirectoryStructureUseCase returns an instance of ProjectDirectoryStructureUseCase.
func NewProjectDirectoryStructureUseCase(
	structureCreator DirectoryStructureCreator,
) *ProjectDirectoryStructureUseCase {
	return &ProjectDirectoryStructureUseCase{structureCreator: structureCreator}
}
