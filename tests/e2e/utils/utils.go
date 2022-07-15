package tests_e2e

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type E2eState struct {
	BinaryPath string
}

func NewState() *E2eState {
	state := &E2eState{
		BinaryPath: fmt.Sprintf("%s/bin/mock", pwd()),
	}

	return state
}

func RunMock(state *E2eState, command string) ([]byte, error) {
	parseCommandVars(&command)
	commandParameters := toCommandParameters(command)

	cmd := exec.Command(state.BinaryPath, commandParameters...)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return []byte(trimCommandOutput(string(result))), err
	}

	return []byte(trimCommandOutput(string(result))), nil
}

func parseCommandVars(command *string) {
	vars := map[string]string{
		"TEST_DATA_PATH": fmt.Sprintf("%s/tests/e2e/data", pwd()),
		"TEST_E2E_PORT":  "4000",
	}

	for key, value := range vars {
		*command = strings.Replace(
			*command,
			fmt.Sprintf("{{%s}}", key),
			value,
			-1,
		)
	}
}

func trimCommandOutput(str string) string {
	lines := strings.Split(str, "\n")

	if strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}

	for strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[0 : len(lines)-2]
	}

	return strings.Join(lines, "\n")
}

func toCommandParameters(command string) []string {
	splitResult := strings.Split(command, " ")

	return splitResult
}

func pwd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s/../..", wd)
}
