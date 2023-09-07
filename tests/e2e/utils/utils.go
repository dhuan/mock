package tests_e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/dhuan/mock/internal/command_parse"
	mocklib "github.com/dhuan/mock/pkg/mock"
	"github.com/stretchr/testify/assert"
)

type E2eState struct {
	BinaryPath string
	Port       int
}

type Response struct {
	Body       []byte
	Headers    map[string]string
	StatusCode int
}

func getFreePort() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	listener.Close()

	return listener.Addr().(*net.TCPAddr).Port
}

func NewState() *E2eState {
	state := &E2eState{
		BinaryPath: fmt.Sprintf("%s/bin/mock", pwd()),
		Port:       getFreePort(),
	}

	return state
}

func RunMock(state *E2eState, command string) ([]byte, error) {
	replaceVars(&command, state)
	commandParameters := command_parse.ToCommandParameters(command)

	cmd := exec.Command(state.BinaryPath, commandParameters...)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return []byte(cleanupCommandOutput(string(result))), err
	}

	return []byte(cleanupCommandOutput(string(result))), nil
}

type KillMockFunc func()

func RunMockBg(state *E2eState, command string, env map[string]string, panicIfServerDidNotStart bool) (KillMockFunc, *bytes.Buffer, *mocklib.MockConfig, bool) {
	replaceVars(&command, state)
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
	if !serverIsReady && panicIfServerDidNotStart {
		panic(fmt.Sprintf("Something went wrong while waiting for mock to start up:\n\n%s", buf))
	}

	return func() {
		err := cmd.Process.Kill()
		if err != nil {
			panic(err)
		}
	}, buf, mocklib.Init(fmt.Sprintf("localhost:%d", state.Port)), serverIsReady
}

func MockAssert(assertOptions *mocklib.AssertOptions, serverOutput *bytes.Buffer, state *E2eState) []mocklib.ValidationError {
	mockConfig := mocklib.Init(fmt.Sprintf("localhost:%d", state.Port))
	validationErrors, err := mocklib.Assert(mockConfig, assertOptions)
	if err != nil {
		log.Println("An error occurred. Here's the server output:")
		fmt.Println(serverOutput)

		panic(err)
	}

	return validationErrors
}

func Request(config *mocklib.MockConfig, method, route string, payload io.Reader, headers map[string]string, serverOutput *bytes.Buffer) *Response {
	url := fmt.Sprintf("http://%s/%s", config.Url, route)
	request, err := http.NewRequest(
		method,
		url,
		payload,
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
		fmt.Println(serverOutput)
		fmt.Printf("Request failed (%s).\n\nServer output:\n\n%s", url, serverOutput)
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

func replaceVars(command *string, state *E2eState) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	vars := map[string]string{
		"TEST_DATA_PATH": fmt.Sprintf("%s/tests/e2e/data", pwd()),
		"TEST_E2E_PORT":  fmt.Sprint(state.Port),
		"WD":             wd,
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

type TestRequest struct {
	Method  string
	Route   string
	Headers map[string]string
	Body    io.Reader
}

func NewGetTestRequest(route string) *TestRequest {
	return &TestRequest{
		Route: route,
	}
}

func NewPostTestRequest(route string) *TestRequest {
	return &TestRequest{
		Route:  route,
		Method: "post",
	}
}

func RunTest(
	t *testing.T,
	configurationFilePath,
	method,
	route string,
	headers map[string]string,
	body io.Reader,
	assertionFunc ...func(t *testing.T, response *Response, serverOutput []byte, state *E2eState),
) {
	request := TestRequest{
		Method:  method,
		Route:   route,
		Headers: headers,
		Body:    body,
	}

	RunTestBase(t, true, configurationFilePath, "", []TestRequest{request}, nil, assertionFunc...)
}

func RunTestWithMultipleRequests(
	t *testing.T,
	configurationFilePath string,
	requests []TestRequest,
	assertionFunc ...func(t *testing.T, response *Response, serverOutput []byte, state *E2eState),
) {
	RunTestBase(t, true, configurationFilePath, "", requests, nil, assertionFunc...)
}

func RunTestWithEnv(
	t *testing.T,
	configurationFilePath,
	method,
	route string,
	headers map[string]string,
	body io.Reader,
	env map[string]string,
	assertionFunc ...func(t *testing.T, response *Response, serverOutput []byte, state *E2eState),
) {
	request := TestRequest{
		Method:  method,
		Route:   route,
		Headers: headers,
		Body:    body,
	}

	RunTestBase(t, true, configurationFilePath, "", []TestRequest{request}, env, assertionFunc...)
}

func RunTestWithArgs(
	t *testing.T,
	configurationFilePath string,
	args []string,
	method,
	route string,
	headers map[string]string,
	body io.Reader,
	assertionFunc ...func(t *testing.T, response *Response, serverOutput []byte, state *E2eState),
) {
	request := TestRequest{
		Method:  method,
		Route:   route,
		Headers: headers,
		Body:    body,
	}

	RunTestBase(t, true, configurationFilePath, strings.Join(args, " "), []TestRequest{request}, map[string]string{}, assertionFunc...)
}

func RunTestWithNoConfigAndWithArgs(
	t *testing.T,
	args []string,
	method,
	route string,
	headers map[string]string,
	body io.Reader,
	assertionFunc ...func(t *testing.T, response *Response, serverOutput []byte, state *E2eState),
) {
	request := TestRequest{
		Method:  method,
		Route:   route,
		Headers: headers,
		Body:    body,
	}

	RunTestBase(t, true, "", strings.Join(args, " "), []TestRequest{request}, map[string]string{}, assertionFunc...)
}

func RunTestWithNoConfigAndWithArgsFailing(
	t *testing.T,
	args []string,
	method,
	route string,
	headers map[string]string,
	body io.Reader,
	assertionFunc ...func(t *testing.T, response *Response, serverOutput []byte, state *E2eState),
) {
	request := TestRequest{
		Method:  method,
		Route:   route,
		Headers: headers,
		Body:    body,
	}

	RunTestBase(t, false, "", strings.Join(args, " "), []TestRequest{request}, map[string]string{}, assertionFunc...)
}

func RunTestWithJsonConfig(
	t *testing.T,
	jsonStr string,
	args []string,
	method,
	route string,
	headers map[string]string,
	body io.Reader,
	assertionFunc ...func(t *testing.T, response *Response, serverOutput []byte, state *E2eState),
) {
	request := TestRequest{
		Method:  method,
		Route:   route,
		Headers: headers,
		Body:    body,
	}

	configFile := mkTemp([]byte(jsonStr))

	RunTestBase(t, true, configFile, strings.Join(args, " "), []TestRequest{request}, map[string]string{}, assertionFunc...)
}

func RunTestWithArgsAndEnv(
	t *testing.T,
	args []string,
	method,
	route string,
	headers map[string]string,
	body io.Reader,
	env map[string]string,
	assertionFunc ...func(t *testing.T, response *Response, serverOutput []byte, state *E2eState),
) {
	request := TestRequest{
		Method:  method,
		Route:   route,
		Headers: headers,
		Body:    body,
	}

	RunTestBase(t, true, "", strings.Join(args, " "), []TestRequest{request}, env, assertionFunc...)
}

func resolveCommand(configurationFilePath string) string {
	if configurationFilePath == "" {
		return "serve -p {{TEST_E2E_PORT}}"
	}

	configPath := fmt.Sprintf("{{TEST_DATA_PATH}}/%s", configurationFilePath)
	if isAbsolutePath(configurationFilePath) {
		configPath = configurationFilePath
	}

	return fmt.Sprintf("serve -c %s -p {{TEST_E2E_PORT}}", configPath)
}

func afterTest(t *testing.T, output *bytes.Buffer) {
	if t.Failed() {
		fmt.Println("\n===========================================================================")
		fmt.Println("Server output for failed test:")
		fmt.Println("===========================================================================")
		fmt.Println(output.String())
		fmt.Println("===========================================================================")
	}
}

func RunTestBase(
	t *testing.T,
	panicIfServerDidNotStart bool,
	configurationFilePath,
	extraArgs string,
	requests []TestRequest,
	env map[string]string,
	assertionFunc ...func(t *testing.T, response *Response, serverOutput []byte, state *E2eState),
) {
	command := resolveCommand(configurationFilePath)
	if extraArgs != "" {
		command = fmt.Sprintf("%s %s", command, extraArgs)
	}

	state := NewState()
	killMock, output, mockConfig, started := RunMockBg(state, command, env, panicIfServerDidNotStart)
	defer killMock()
	defer afterTest(t, output)

	response := &Response{}

	if started {
		for i := range requests {
			response = Request(mockConfig, requests[i].Method, requests[i].Route, requests[i].Body, requests[i].Headers, output)
		}
	}

	for i := range assertionFunc {
		assertionFunc[i](t, response, output.Bytes(), state)
	}
}

func StringMatches(expected string) func(t *testing.T, response *Response, serverOutput []byte, state *E2eState) {
	return func(t *testing.T, response *Response, serverOutput []byte, state *E2eState) {
		replaceVars(&expected, state)

		assert.Equal(t, expected, string(response.Body))
	}
}

func LineEquals(lineNumber int, expectedLine string) func(t *testing.T, response *Response, serverOutput []byte, state *E2eState) {
	return func(t *testing.T, response *Response, serverOutput []byte, state *E2eState) {
		replaceVars(&expectedLine, state)

		assert.Equal(t, expectedLine, getLineFromString(lineNumber-1, string(response.Body)))
	}
}

func ApplicationOutputHasLines(expectedLines []string) func(t *testing.T, response *Response, serverOutput []byte, state *E2eState) {
	return func(t *testing.T, response *Response, serverOutput []byte, state *E2eState) {
		if len(expectedLines) == 0 {
			return
		}

		serverOutputLines := breakLines(string(serverOutput))
		lineMatch := -1

		for i := range expectedLines {
			replaceVars(&expectedLines[i], state)
		}

		for i := range serverOutputLines {
			serverOutputLines[i] = removeLogDatePrefix(serverOutputLines[i])

			if serverOutputLines[i] == expectedLines[0] {
				lineMatch = i
			}
		}

		if lineMatch == -1 {
			t.Fatal(fmt.Sprintf("There is no line matching: %s", expectedLines[0]))
		}

		i := 0
		for {
			expectedLine := expectedLines[i]

			if expectedLine != serverOutputLines[lineMatch] {
				fmt.Println(fmt.Sprintf("Line expected: %s\nLine actual:   %s", expectedLine, serverOutputLines[lineMatch]))
				t.Fail()
			}

			lineMatch = lineMatch + 1
			i = i + 1

			if i > len(expectedLines)-1 {
				break
			}
		}
	}
}

func ApplicationOutputMatches(expectedLines []string) func(t *testing.T, response *Response, serverOutput []byte, state *E2eState) {
	return func(t *testing.T, response *Response, serverOutput []byte, state *E2eState) {
		if len(expectedLines) == 0 {
			return
		}

		serverOutputLines := breakLines(string(serverOutput))

		i := 0
		for {
			replaceVars(&expectedLines[i], state)
			expectedLine := expectedLines[i]

			if expectedLine != removeLogDatePrefix(serverOutputLines[i]) {
				fmt.Println(fmt.Sprintf("Line expected: %s\nLine actual:   %s", expectedLine, serverOutputLines[i]))
				t.Fail()
			}

			i = i + 1

			if i > len(expectedLines)-1 {
				break
			}
		}
	}
}

func LineRegexMatches(lineNumber int, regex string) func(t *testing.T, response *Response, serverOutput []byte, state *E2eState) {
	return func(t *testing.T, response *Response, serverOutput []byte, state *E2eState) {
		assert.Regexp(
			t,
			regexp.MustCompile(regex),
			getLineFromString(lineNumber-1, string(response.Body)),
		)
	}
}

func StatusCodeMatches(expectedStatusCode int) func(t *testing.T, response *Response, serverOutput []byte, state *E2eState) {
	return func(t *testing.T, response *Response, serverOutput []byte, state *E2eState) {
		assert.Equal(t, expectedStatusCode, response.StatusCode)
	}
}

func HeadersMatch(expectedHeaders map[string]string) func(t *testing.T, response *Response, serverOutput []byte, state *E2eState) {
	return func(t *testing.T, response *Response, serverOutput []byte, state *E2eState) {
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

func HeaderKeysNotIncluded(headerKeys []string) func(t *testing.T, response *Response, serverOutput []byte, state *E2eState) {
	return func(t *testing.T, response *Response, serverOutput []byte, state *E2eState) {
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

func JsonMatches(expectedJson map[string]interface{}) func(t *testing.T, response *Response, serverOutput []byte, state *E2eState) {
	return func(t *testing.T, response *Response, serverOutput []byte, state *E2eState) {
		jsonEncodedA, err := json.Marshal(expectedJson)
		if err != nil {
			t.Fatal("Failed to parse JSON from expected input!")
		}

		jsonEncodedB, err := encodeJsonAgain(response.Body)
		if err != nil {
			t.Fatal("Failed to parse JSON from response!")
		}

		jsonA := string(jsonEncodedA)
		jsonB := string(jsonEncodedB)

		replaceVars(&jsonA, state)
		replaceVars(&jsonB, state)

		assert.Equal(t, jsonA, jsonB)
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
	for i := range list {
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

func BuildFormPayload(data map[string]string) io.Reader {
	dataParsed := make(map[string][]string)

	for key, value := range data {
		dataParsed[key] = []string{value}
	}

	return strings.NewReader(url.Values(dataParsed).Encode())
}

func breakLines(text string) []string {
	return strings.Split(text, "\n")
}

func removeLogDatePrefix(text string) string {
	return replaceRegex(text, []string{`^[0-9/]{1,} [0-9\:]{1,} `}, "")
}

func EnvVarExists(varName string) bool {
	_, exists := os.LookupEnv(varName)

	return exists
}

func mkTemp(content []byte) string {
	result, err := exec.Command("mktemp").Output()
	if err != nil {
		panic(err)
	}

	filePath := strings.TrimSuffix(string(result), "\n")

	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		panic(err)
	}

	return filePath
}

func beginsWith(subject, find string) bool {
	return strings.Index(subject, find) == 0
}

func isAbsolutePath(path string) bool {
	return beginsWith(path, "/")
}
