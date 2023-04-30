package behaviour

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/cucumber/godog"
	assert "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (test *Test) registerAllSteps(sc *godog.ScenarioContext) {
	test.log.Info("register all steps")
	sc.Step(`^I execute the cli command$`, test.iExecuteTheCliCommand)
	sc.Step(`^I must get an exit code (\d+)$`, test.iMustGetAnExitCode)
	sc.Step(`^I must get a command output$`, test.iMustGetACommandOutput)
	sc.Step(`^I have the following configuration$`, test.iHaveTheFollowingConfiguration)
	sc.Step(`^a "([^"]*)" file with the following contents$`, test.aFileWithTheFollowingContents)
}

func (test *Test) iExecuteTheCliCommand(command *godog.DocString) error {
	if len(command.Content) == 0 {
		return errors.New("command string can't be empty")
	}

	parts := strings.Split(test.parseCommand(command.Content), " ")
	test.cmd.output, test.cmd.error = exec.Command(parts[0], parts[1:]...).Output() //nolint:gosec

	return nil
}

func (test *Test) iMustGetAnExitCode(exitCode int) error {
	switch exitCode {
	case 0:
		if !assert.Expect(test.cmd.error).To(assert.BeNil()) {
			return errors.New("expected error to be nil")
		}
	case 1:
		if !assert.Expect(test.cmd.error != nil).To(assert.BeTrue()) {
			return errors.New("expected error not to be nil")
		}
	default:
		return errors.New("unknown exit code")
	}

	return nil
}

func (test *Test) iMustGetACommandOutput(expected *godog.DocString) error {
	if !assert.Expect(test.commandOutput()).To(assert.Equal(expected.Content)) {
		return errors.New("actual command output does not match the expected command output")
	}

	return nil
}

func (test *Test) parseCommand(cmd string) string {
	buf := new(bytes.Buffer)
	tmpl := template.Must(template.New("cmd").Parse(cmd))

	if err := tmpl.Execute(buf, test); err != nil {
		test.log.Fatal("failed to execute template", zap.Error(err))
	}

	return buf.String()
}

func (test *Test) iHaveTheFollowingConfiguration(config *godog.DocString) error {
	data := []byte(config.Content)

	test.ConfigFile = fmt.Sprintf("%s/testdata/.fundi.yaml", test.workingDirectory())

	if err := os.WriteFile(test.ConfigFile, data, 0600); err != nil {
		return err
	}

	return nil
}

func (test *Test) aFileWithTheFollowingContents(fileName string, fileData *godog.DocString) error {
	data := []byte(fileData.Content)

	if err := os.WriteFile(
		fmt.Sprintf("%s/testdata/%s", test.workingDirectory(), fileName),
		data,
		0600,
	); err != nil {
		return err
	}

	return nil
}
