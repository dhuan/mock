package tests_e2e

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	mocklib "github.com/dhuan/mock/pkg/mock"
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

type KillMockFunc func()

func RunMockBg(state *E2eState, command string) KillMockFunc {
	parseCommandVars(&command)
	commandParameters := toCommandParameters(command)

	cmd := exec.Command(state.BinaryPath, commandParameters...)
	buf := &bytes.Buffer{}
	cmd.Stdout = buf
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	serverIsReady := waitForOutputInCommand("Mock server is listening on port 4000.", 4, buf)
	if !serverIsReady {
		panic("Something went wrong while waiting for mock to start up.")
	}

	return func() {
		err := cmd.Process.Kill()
		if err != nil {
			panic(err)
		}
	}
}

func Request(config *mocklib.MockConfig, method, route, payload string, headers map[string]string) {
	request, err := http.NewRequest(
		method,
		fmt.Sprintf("http://%s/%s", config.Url, route),
		bytes.NewBuffer([]byte(payload)),
	)
	if err != nil {
		panic(err)
	}

	for headerKey, headerValue := range headers {
		request.Header.Set(headerKey, headerValue)
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		panic(err)
	}
}

func waitForOutputInCommand(expectedOutput string, attempts int, buffer *bytes.Buffer) bool {
	for attempts > 0 {
		if strings.Contains(buffer.String(), expectedOutput) {
			return true
		}

		time.Sleep(500 * time.Millisecond)

		attempts--
	}

	return false
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
