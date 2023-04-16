package app

import (
	"errors"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestReadYAMLFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	fileName := "test.yml"

	tests := map[string]struct {
		expectedErr error
		fileData    []byte
		fileName    string
	}{
		"when the file does not exist, return an error": {
			expectedErr: errors.New("failed to read file unknown-file.yml: open unknown-file.yml: file does not exist"),
			fileData:    []byte(`"key": "value"`),
			fileName:    "unknown-file.yml",
		},
		"when the reader fails to unmarshal the YAML file, return an error": {
			expectedErr: errors.New("failed to unmarshal YAML data: yaml: did not find expected whitespace or line break"),
			fileData:    []byte(`!*#$%`),
			fileName:    fileName,
		},
		"when the reader successfully reads the YAML file, return no error": {
			fileName: fileName,
			fileData: []byte(`
metadata:
  output: "."
  templates: "./templates"
  values: "./values.yml"
directories:
  - name: project_name
    files:
      - name: README.md
        template: readme.md.tmpl
`),
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			err := afero.WriteFile(fs, fileName, testCase.fileData, 0644)
			assert.NoError(t, err)

			reader := fileReader{fs: fs}

			cfg, err := reader.readYAMLFile(testCase.fileName)

			switch testCase.expectedErr != nil {
			case true:
				assert.EqualError(t, err, testCase.expectedErr.Error())
			case false:
				assert.NoError(t, err)
				assert.Equal(t, ".", cfg.Metadata.Output)
				assert.Equal(t, "./templates", cfg.Metadata.Templates)
				assert.Equal(t, "./values.yml", cfg.Metadata.Values)
				assert.Len(t, cfg.Directories, 1)
				assert.Equal(t, "project_name", cfg.Directories[0].Name)
				assert.Len(t, cfg.Directories[0].Files, 1)
				assert.Equal(t, "README.md", cfg.Directories[0].Files[0].Name)
				assert.Equal(t, "readme.md.tmpl", cfg.Directories[0].Files[0].Template)
			}
		})
	}
}
