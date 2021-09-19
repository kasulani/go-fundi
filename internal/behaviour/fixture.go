package behaviour

import (
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (specs *TestSpecifications) workingDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		specs.log.Fatal("failed to get working directory", zap.Error(err))
	}

	return dir
}

func (specs *TestSpecifications) setInitialContext(input string) error {
	specs.log.Info("set initial context")
	switch input {
	case "a good fundi file":
		specs.in.File = specs.workingDirectory() + string(os.PathSeparator) + "testdata/.test.fundi.yaml"
	case "a bad fundi file":
		specs.in.File = specs.workingDirectory() + string(os.PathSeparator) + "testdata/.non-existing.fundi.yaml"
	default:
		return errors.New("unknown input")
	}

	return nil
}
