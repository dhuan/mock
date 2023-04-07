package tests_e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	mocklib "github.com/dhuan/mock/pkg/mock"
	"github.com/dhuan/mock/internal/command_parse"
	"github.com/stretchr/testify/assert"
)

type E2eState struct {
	BinaryPath string
}

type Response struct {
	Body       []byte
	Headers    map[string]string
	StatusCode int
}

func GetTestPort() string {
	port := os.Getenv("MOCK_TEST_PORT")
	if port == "" {
		port = "4000"
	}

	return port
}

func NewState() *E2eState {
	state := &E2eState{
		BinaryPath: fmt.Sprintf("%s/bin/mock", pwd()),
	}

	return state
}

func RunMock(state *E2eState, command string) ([]byte, error) {
	replaceVars(&command)
	commandParameters := command_parse.ToCommandParameters(command)

	cmd := exec.Command(state.BinaryPath, commandParameters...)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return []byte(cleanupCommandOutput(string(result))), err
	}

	return []byte(cleanupCommandOutput(string(result))), nil
}

type KillMockFunc func()

func RunMockBg(state *E2eState, command string, env map[string]string) (KillMockFunc, *bytes.Buffer, *mocklib.MockConfig) {
	replaceVars(&command)
	commandParameters := command_parse.ToCommandParameters(command)

	cmd := exec.Command(state.BinaryPath, commandParameters...)
	buf := &bytes.Buffer{}
	cmd.Stdout = buf
	cmd.Stderr = buf
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, parseEnv(env)...)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	serverIsReady := waitForOutputInCommand("Starting server on port", 4, buf)
	if !serverIsReady {
		panic("Something went wrong while waiting for mock to start up.")
	}

	return func() {
		err := cmd.Process.Kill()
		if err != nil {
			panic(err)
		}
	}, buf, mocklib.Init(fmt.Sprintf("localhost:%s", GetTestPort()))
}

func MockAssert(assertConfig *mocklib.AssertConfig, serverOutput *bytes.Buffer) []mocklib.ValidationError {
	mockConfig := mocklib.Init(fmt.Sprintf("localhost:%s", GetTestPort()))
	validationErrors, err := mocklib.Assert(mockConfig, assertConfig)
	if err != nil {
		log.Println("An error occurred. Here's the server output:")
		fmt.Println(serverOutput)

		panic(err)
	}

	return validationErrors
}

func Request(config *mocklib.MockConfig, method, route, payload string, headers map[string]string) *Response {
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

	return &Response{
		Body:       responseBody,
		Headers:    parseHeaders(response.Header),
		StatusCode: response.StatusCode,
	}
}

func parseHeaders(headers http.Header) map[string]string {
	parsedHeaders := make(map[string]string)
	sortedKeys := getSortedKeys(headers)

	for _, key := range sortedKeys {
		parsedHeaders[key] = strings.Join(headers[key], ",")
	}

	return parsedHeaders
}

func RequestApiReset(config *mocklib.MockConfig) {
	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://%s/__mock__/reset", config.Url),
		nil,
	)
	if err != nil {
		panic(err)
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

func replaceVars(command *string) {
	vars := map[string]string{
		"TEST_DATA_PATH": fmt.Sprintf("%s/tests/e2e/data", pwd()),
		"TEST_E2E_PORT":  GetTestPort(),
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

	for i := range lines {
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
	headers map[string]string,
	body string,
	assertionFunc ...func(t *testing.T, response *Response),
) {
	RunTestBase(t, configurationFilePath, "", method, route, headers, body, nil, assertionFunc...)
}

func RunTestWithEnv(
	t *testing.T,
	configurationFilePath,
	method,
	route string,
	headers map[string]string,
	body string,
	env map[string]string,
	assertionFunc ...func(t *testing.T, response *Response),
) {
	RunTestBase(t, configurationFilePath, "", method, route, headers, body, env, assertionFunc...)
}

func RunTestWithArgs(
	t *testing.T,
	configurationFilePath string,
	args []string,
	method,
	route string,
	headers map[string]string,
	body string,
	assertionFunc ...func(t *testing.T, response *Response),
) {
	RunTestBase(t, configurationFilePath, strings.Join(args, " "), method, route, headers, body, map[string]string{}, assertionFunc...)
}

func RunTestWithNoConfigAndWithArgs(
	t *testing.T,
	args []string,
	method,
	route string,
	headers map[string]string,
	body string,
	assertionFunc ...func(t *testing.T, response *Response),
) {
	RunTestBase(t, "", strings.Join(args, " "), method, route, headers, body, map[string]string{}, assertionFunc...)
}

func resolveCommand(configurationFilePath string) string {
	if configurationFilePath == "" {
		return "serve -p {{TEST_E2E_PORT}}"
	}

	return fmt.Sprintf("serve -c {{TEST_DATA_PATH}}/%s -p {{TEST_E2E_PORT}}", configurationFilePath)
}

func RunTestBase(
	t *testing.T,
	configurationFilePath,
	extraArgs,
	method,
	route string,
	headers map[string]string,
	body string,
	env map[string]string,
	assertionFunc ...func(t *testing.T, response *Response),
) {
	command := resolveCommand(configurationFilePath)
	if extraArgs != "" {
		command = fmt.Sprintf("%s %s", command, extraArgs)
	}

	killMock, _, mockConfig := RunMockBg(NewState(), command, env)
	defer killMock()

	response := Request(mockConfig, method, route, body, headers)

	for i := range assertionFunc {
		assertionFunc[i](t, response)
	}
}

func StringMatches(expected string) func(t *testing.T, response *Response) {
	return func(t *testing.T, response *Response) {
		replaceVars(&expected)

		assert.Equal(t, expected, string(response.Body))
	}
}

func LineEquals(lineNumber int, expectedLine string) func(t *testing.T, response *Response) {
	return func(t *testing.T, response *Response) {
		replaceVars(&expectedLine)

		assert.Equal(t, expectedLine, getLineFromString(lineNumber-1, string(response.Body)))
	}
}

func LineRegexMatches(lineNumber int, regex string) func(t *testing.T, response *Response) {
	return func(t *testing.T, response *Response) {
		assert.Regexp(
			t,
			regexp.MustCompile(regex),
			getLineFromString(lineNumber-1, string(response.Body)),
		)
	}
}

func StatusCodeMatches(expectedStatusCode int) func(t *testing.T, response *Response) {
	return func(t *testing.T, response *Response) {
		assert.Equal(t, expectedStatusCode, response.StatusCode)
	}
}

func HeadersMatch(expectedHeaders map[string]string) func(t *testing.T, response *Response) {
	return func(t *testing.T, response *Response) {
		expectedHeadersKeys := getSortedKeys(expectedHeaders)

		for _, expectedHeaderKey := range expectedHeadersKeys {
			headerValue, ok := response.Headers[expectedHeaderKey]
			if !ok {
				t.Error(
					fmt.Sprintf("Header key does not exist in the resulting request: %s", expectedHeaderKey),
				)

				return
			}

			assert.Equal(t, expectedHeaders[expectedHeaderKey], headerValue)
		}
	}
}

func HeaderKeysNotIncluded(headerKeys []string) func(t *testing.T, response *Response) {
	return func(t *testing.T, response *Response) {
		for _, headerKey := range headerKeys {
			_, exists := response.Headers[headerKey]

			if exists {
				t.Error(
					fmt.Sprintf("Expected header key to not exist, but it does: %s", headerKey),
				)
			}
		}
	}
}

func JsonMatches(expectedJson map[string]interface{}) func(t *testing.T, response *Response) {
	return func(t *testing.T, response *Response) {
		jsonEncodedA, err := json.Marshal(expectedJson)
		if err != nil {
			t.Fatal("Failed to parse JSON from expected input!")
		}

		jsonEncodedB, err := encodeJsonAgain(response.Body)
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

func getSortedKeys[T interface{}](subject map[string]T) []string {
	keys := GetKeys(subject)
	sort.Strings(keys)

	return keys
}

func GetKeys[T_Key comparable, T_Value interface{}](subject map[T_Key]T_Value) []T_Key {
	keys := make([]T_Key, 0, len(subject))

	for key := range subject {
		keys = append(keys, key)
	}

	return keys
}

func AssertMapHasValues[T_Key comparable, T_Value comparable](
	t *testing.T,
	subject map[T_Key]T_Value,
	values map[T_Key]T_Value,
) {
	for key, value := range values {
		valueb, ok := subject[key]

		if !ok {
			t.Error(fmt.Sprintf("Key '%+v' does not exist in given map.", key))
		}

		assert.Equal(t, value, valueb)
	}
}

func IndexOf[T comparable](list []T, value T) int {
	for i, _ := range list {
		if list[i] == value {
			return i
		}
	}

	return -1
}

func getLineFromString(lineNumber int, str string) string {
	lines := strings.Split(str, "\n")

	if (len(lines) - 1) < lineNumber {
		return ""
	}

	return lines[lineNumber]
}

func AssertTimeDifferenceLessThanSeconds(t *testing.T, timeA, timeB time.Time, seconds int) {
	a := timeA.Unix()
	b := timeB.Unix()

	diffSeconds := int(b - a)

	assert.Less(t, diffSeconds, seconds)
}

func AssertTimeDifferenceEqualOrMoreThanSeconds(t *testing.T, timeA, timeB time.Time, seconds int) {
	a := timeA.Unix()
	b := timeB.Unix()

	diffSeconds := int(b - a)

	assert.GreaterOrEqual(t, diffSeconds, seconds)
}

func parseEnv(env map[string]string) []string {
	result := make([]string, len(env))
	keys := GetKeys(env)

	for i, key := range keys {
		result[i] = fmt.Sprintf("%s=%s", key, env[key])
	}

	return result
}
