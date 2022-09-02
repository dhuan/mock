package tests_e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	mocklib "github.com/dhuan/mock/pkg/mock"
	"github.com/stretchr/testify/assert"
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
		return []byte(cleanupCommandOutput(string(result))), err
	}

	return []byte(cleanupCommandOutput(string(result))), nil
}

type KillMockFunc func()

func RunMockBg(state *E2eState, command string) KillMockFunc {
	parseCommandVars(&command)
	commandParameters := toCommandParameters(command)

	cmd := exec.Command(state.BinaryPath, commandParameters...)
	buf := &bytes.Buffer{}
	cmd.Stdout = buf
	cmd.Stderr = buf
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	serverIsReady := waitForOutputInCommand("Starting Mock server on port 4000.", 4, buf)
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

func MockAssert(assertConfig *mocklib.AssertConfig) []mocklib.ValidationError {
	mockConfig := mocklib.Init("localhost:4000")
	validationErrors, err := mocklib.Assert(mockConfig, assertConfig)
	if err != nil {
		panic(err)
	}

	return validationErrors
}

func Request(config *mocklib.MockConfig, method, route, payload string, headers map[string]string) []byte {
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
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return responseBody
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

func replaceRegex(subject string, find []string, replaceWith string) string {
	if len(find) == 0 {
		return subject
	}

	re := regexp.MustCompile(find[0])

	return replaceRegex(
		re.ReplaceAllString(subject, replaceWith),
		find[1:],
		replaceWith,
	)
}

func replaceRegexForEachLine(subject string, find []string, replaceWith string) string {
	lines := strings.Split(subject, "\n")

	for i, _ := range lines {
		lines[i] = replaceRegex(lines[i], find, replaceWith)
	}

	return strings.Join(lines, "\n")
}

func cleanupCommandOutput(str string) string {
	return replaceRegexForEachLine(trimCommandOutput(str), []string{`^[0-9\/]{1,} [0-9\:]{1,} `}, "")
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

func RunTest(
	t *testing.T,
	configurationFilePath,
	method,
	route string,
	assertionFunc func(t *testing.T, response []byte),
) {
	killMock := RunMockBg(NewState(), fmt.Sprintf("serve -c {{TEST_DATA_PATH}}/%s -p {{TEST_E2E_PORT}}", configurationFilePath))
	defer killMock()

	mockConfig := mocklib.Init("localhost:4000")
	responseBody := Request(mockConfig, method, route, "", map[string]string{})

	assertionFunc(t, responseBody)
}

func StringMatches(expected string) func(t *testing.T, response []byte) {
	return func(t *testing.T, responseBody []byte) {
		assert.Equal(t, expected, string(responseBody))
	}
}

func JsonMatches(expectedJson map[string]interface{}) func(t *testing.T, response []byte) {
	return func(t *testing.T, responseBody []byte) {
		jsonEncodedA, err := json.Marshal(expectedJson)
		if err != nil {
			t.Fatal("Failed to parse JSON from expected input!")
		}

		jsonEncodedB, err := encodeJsonAgain(responseBody)
		if err != nil {
			t.Fatal("Failed to parse JSON from response!")
		}

		assert.Equal(t, string(jsonEncodedA), string(jsonEncodedB))
	}
}

func encodeJsonAgain(encodedJson []byte) ([]byte, error) {
	var jsonTarget map[string]interface{}
	err := json.Unmarshal(encodedJson, &jsonTarget)
	if err != nil {
		return []byte(""), err
	}

	return json.Marshal(jsonTarget)
}

var ContentTypeJsonHeaders map[string]string = map[string]string{
	"Content-type": "application/json",
}
