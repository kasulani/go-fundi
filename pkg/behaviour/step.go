package behaviour

import (
	"bytes"
	"os/exec"
	"strings"
	"text/template"

	"github.com/cucumber/godog"
	assert "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (specs *TestSpecifications) registerAllSteps(sc *godog.ScenarioContext) {
	specs.log.Info("register all steps")
	sc.Step(`^I execute the cli command$`, specs.iExecuteTheCliCommand)
	sc.Step(`^I must get an exit code (\d+)$`, specs.iMustGetAnExitCode)
	sc.Step(`^I must get a command output$`, specs.iMustGetACommandOutput)
}

func (specs *TestSpecifications) iExecuteTheCliCommand(command *godog.DocString) error {
	if len(command.Content) == 0 {
		return errors.New("command string can't be empty")
	}

	parts := strings.Split(specs.parseCommand(command.Content), " ")
	specs.out.cmdOutput, specs.out.error = exec.Command(parts[0], parts[1:]...).Output()

	return nil
}

func (specs *TestSpecifications) iMustGetAnExitCode(exitCode int) error {
	switch exitCode {
	case 0:
		if !assert.Expect(specs.out.error).To(assert.BeNil()) {
			return errors.New("expected error to be nil")
		}
	case 1:
		if !assert.Expect(specs.out.error != nil).To(assert.BeTrue()) {
			return errors.New("expected error not to be nil")
		}
	default:
		return errors.New("unknown exit code")
	}

	return nil
}

func (specs *TestSpecifications) iMustGetACommandOutput(expected *godog.DocString) error {
	if !assert.Expect(specs.commandOutput()).To(assert.Equal(expected.Content)) {
		return errors.New("actual command output does not match the expected command output")
	}

	return nil
}

func (specs *TestSpecifications) parseCommand(cmd string) string {
	buf := new(bytes.Buffer)
	tmpl := template.Must(template.New("cmd").Parse(cmd))

	if err := tmpl.Execute(buf, specs.in); err != nil {
		specs.log.Fatal("failed to execute template", zap.Error(err))
	}

	return buf.String()
}
