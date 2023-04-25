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

func (specs *TestSpecifications) registerAllSteps(sc *godog.ScenarioContext) {
	specs.log.Info("register all steps")
	sc.Step(`^I have "([^"]*)"$`, specs.iHave)
	sc.Step(`^I execute the cli command$`, specs.iExecuteTheCliCommand)
	sc.Step(`^I must get an exit code (\d+)$`, specs.iMustGetAnExitCode)
	sc.Step(`^I must get a command output$`, specs.iMustGetACommandOutput)
	sc.Step(`^file "([^"]*)" has contents$`, specs.fileHasContents)
	sc.Step(`^I have the following configuration$`, specs.iHaveTheFollowingConfiguration)
	sc.Step(`^a "([^"]*)" file with the following contents$`, specs.aFileWithTheFollowingContents)
}

func (specs *TestSpecifications) iHave(input string) error {
	return specs.setInitialContext(input)
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

func (specs *TestSpecifications) fileHasContents(filename string, expected *godog.DocString) error {
	var data []byte
	var err error

	switch filename {
	case "doc.go":
		data, err = os.ReadFile("./funditest/pkg/app/doc.go")
		if err != nil {
			specs.log.Fatal("failed to open file", zap.Error(err))

			return err
		}
	default:
		return errors.New("unknown file name")
	}

	if !assert.Expect(string(data)).To(assert.BeIdenticalTo(expected.Content)) {
		return errors.New("actual command output does not match the expected command output")
	}

	return nil
}

func (specs *TestSpecifications) iHaveTheFollowingConfiguration(config *godog.DocString) error {
	data := []byte(config.Content)

	specs.in.File = fmt.Sprintf("%s/testdata/.fundi.yaml", specs.workingDirectory())

	if err := os.WriteFile(specs.in.File, data, 0600); err != nil {
		return err
	}

	return nil
}

func (specs *TestSpecifications) aFileWithTheFollowingContents(fileName string, fileData *godog.DocString) error {
	data := []byte(fileData.Content)

	if err := os.WriteFile(
		fmt.Sprintf("%s/testdata/%s", specs.workingDirectory(), fileName),
		data,
		0600,
	); err != nil {
		return err
	}

	return nil
}
