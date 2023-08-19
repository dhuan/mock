package tests_e2e

import (
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_Response_FileResponse(t *testing.T) {
	RunTest(
		t,
		"config_with_file_responses/config.json",
		"POST",
		"foo/bar/1",
		nil,
		strings.NewReader(""),
		StringMatches("Hello world! This is response A.\n"),
	)
}

func Test_E2E_Response_FileResponse_WithResponseFileFlag(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route hello/world",
			"--response-file data/response.txt",
		},
		"GET",
		"hello/world",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(200),
		StringMatches("Hello world!\n"),
	)
}

func Test_E2E_Response_ResponseInsideFolder(t *testing.T) {
	RunTest(
		t,
		"config_with_file_responses/config.json",
		"POST",
		"foo/bar/2",
		nil,
		strings.NewReader(""),
		StringMatches("This test asserts that you can set response files inside folders.\n"),
	)
}

func Test_E2E_Response_JsonResponse(t *testing.T) {
	RunTest(
		t,
		"config_with_file_responses/config.json",
		"POST",
		"foo/bar/3",
		nil,
		strings.NewReader(""),
		JsonMatches(map[string]interface{}{
			"response_text": "This is a JSON response.",
		}),
	)
}

func Test_E2E_Response_WithFile_WithAbsolutePath(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filePath := fmt.Sprintf("%s/data/response.txt", pwd)

	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route hello/world",
			"--method get",
			fmt.Sprintf("--response 'file:%s'", filePath),
		},
		"GET",
		"hello/world",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(200),
		StringMatches("Hello world!\n"),
	)
}

func Test_E2E_Response_WithDynamicFileName(t *testing.T) {
	RunTest(
		t,
		"config_with_file_responses/config.json",
		"GET",
		"books/i_robot/content",
		nil,
		strings.NewReader(""),
		StringMatches("This is the book 'I, Robot'.\n"),
	)

	RunTest(
		t,
		"config_with_file_responses/config.json",
		"GET",
		"books/nightfall/content",
		nil,
		strings.NewReader(""),
		StringMatches("This is the book 'Nightfall'.\n"),
	)
}

func Test_E2E_Response_ShellScript(t *testing.T) {
	RunTest(
		t,
		"config_with_script_responses/config.json",
		"GET",
		"foo/bar",
		nil,
		strings.NewReader(""),
		StringMatches("Hello world! This response was generated from a shell script."),
	)
}

func Test_E2E_Response_ShellScript_WithCmdParams(t *testing.T) {
	responseFormats := []string{
		"--response 'sh:data/config_with_script_responses/handler.sh'",
		"--response-sh data/config_with_script_responses/handler.sh",
		"--shell-script data/config_with_script_responses/handler.sh",
	}

	for _, responseFormat := range responseFormats {
		RunTestWithNoConfigAndWithArgs(
			t,
			[]string{
				"--route foo/bar",
				responseFormat,
			},
			"GET",
			"foo/bar",
			nil,
			strings.NewReader(""),
			StatusCodeMatches(200),
			StringMatches("Hello world! This response was generated from a shell script."),
		)
	}
}

func Test_E2E_Response_ShellScript_WithAbsolutePath(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filePath := fmt.Sprintf("%s/data/config_with_script_responses/handler.sh", pwd)

	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--response 'sh:%s'", filePath),
		},
		"GET",
		"foo/bar",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(200),
		StringMatches("Hello world! This response was generated from a shell script."),
	)
}

func Test_E2E_Response_ShellScript_RequestDetailsFromEnvVariables(t *testing.T) {
	RunTest(
		t,
		"config_with_script_responses/config.json",
		"GET",
		"foo/bar/2?some_key=some_value&another_key=another_value",
		map[string]string{
			"Some-Header-Key":    "Some-Header-Value",
			"Another-Header-Key": "Another-Header-Value",
		},
		strings.NewReader(""),
		StringMatches(`Server Host: localhost:{{TEST_E2E_PORT}}
Request Host: localhost:{{TEST_E2E_PORT}}
URL: http://localhost:{{TEST_E2E_PORT}}/foo/bar/2
Endpoint: foo/bar/2
Method: get
Querystring: some_key=some_value&another_key=another_value
Headers:
accept-encoding: gzip
another-header-key: Another-Header-Value
some-header-key: Some-Header-Value
user-agent: Go-http-client/1.1`),
	)
}

func Test_E2E_Response_ShellScript_RequestDetailsFromEnvVariables_WithPayload(t *testing.T) {
	RunTest(
		t,
		"config_with_script_responses/config.json",
		"POST",
		"foo/bar/3",
		nil,
		strings.NewReader("This is the request payload."),
		StringMatches(`Payload:
This is the request payload.`),
	)
}

func Test_E2E_Response_ShellScript_CustomHeadersAndStatusCode(t *testing.T) {
	RunTest(
		t,
		"config_with_script_responses/config.json",
		"POST",
		"foo/bar/4",
		nil,
		strings.NewReader("This is the request payload."),
		StringMatches(`Hello world!`),
		HeadersMatch(map[string]string{
			"Some-Header-Key":    "Some Header Value",
			"Another-Header-Key": "Another Header Value",
		}),
		StatusCodeMatches(201),
	)
}

func Test_E2E_Response_ShellScript_WithParameter(t *testing.T) {
	RunTest(
		t,
		"config_with_script_responses/config.json",
		"POST",
		"foo/bar/5",
		nil,
		strings.NewReader("This is the request payload."),
		StringMatches(`Parameter: foobar`),
	)
}

func Test_E2E_Response_ShellScript_CommandFailing(t *testing.T) {
	runningInGithubCi := EnvVarExists("CI")
	expectedFailureLine := "{{WD}}/data/config_with_script_responses/handler_with_command_that_fails.sh: line 3: please_fail: command not found"
	if runningInGithubCi {
		expectedFailureLine = "{{WD}}/data/config_with_script_responses/handler_with_command_that_fails.sh: 3: please_fail: not found"
	}

	RunTest(
		t,
		"config_with_script_responses/config.json",
		"POST",
		"foo/bar/6",
		nil,
		strings.NewReader("This is the request payload."),
		LineEquals(1, `Hello world!`),
		ApplicationOutputHasLines([]string{
			"Executing shell script located in {{WD}}/data/config_with_script_responses/handler_with_command_that_fails.sh",
			"Output from program execution:",
			"",
			expectedFailureLine,
		}),
	)
}

func Test_E2E_Response_ShellScript_ReadingEndpointParams(t *testing.T) {
	RunTest(
		t,
		"config_with_script_responses/config.json",
		"GET",
		"users/country/brazil/page/7",
		nil,
		strings.NewReader(""),
		LineEquals(1, `Country: brazil`),
		LineEquals(2, `Page: 7`),
	)
}

func Test_E2E_Response_Json_UsingVariables(t *testing.T) {
	RunTest(
		t,
		"config_responses_using_variables/config.json",
		"GET",
		"response_json_using_variables?param_one=value_one&param_two=value_two",
		nil,
		strings.NewReader(""),
		JsonMatches(map[string]interface{}{
			"MOCK_HOST":                          "localhost:{{TEST_E2E_PORT}}",
			"MOCK_REQUEST_HOST":                  "localhost:{{TEST_E2E_PORT}}",
			"MOCK_REQUEST_URL":                   "http://localhost:{{TEST_E2E_PORT}}/response_json_using_variables",
			"MOCK_REQUEST_ENDPOINT":              "response_json_using_variables",
			"MOCK_REQUEST_METHOD":                "get",
			"MOCK_REQUEST_QUERYSTRING":           "param_one=value_one&param_two=value_two",
			"MOCK_REQUEST_QUERYSTRING_PARAM_ONE": "value_one",
			"MOCK_REQUEST_QUERYSTRING_PARAM_TWO": "value_two",
		}),
	)
}

func Test_E2E_Response_PlainText_UsingVariables(t *testing.T) {
	responseStr := strings.Join([]string{
		"MOCK_HOST: ${MOCK_HOST}",
		"MOCK_REQUEST_HOST: ${MOCK_REQUEST_HOST}",
		"MOCK_REQUEST_URL: ${MOCK_REQUEST_URL}",
		"MOCK_REQUEST_ENDPOINT: ${MOCK_REQUEST_ENDPOINT}",
		"MOCK_REQUEST_METHOD: ${MOCK_REQUEST_METHOD}",
		"MOCK_REQUEST_QUERYSTRING: ${MOCK_REQUEST_QUERYSTRING}",
		"MOCK_REQUEST_QUERYSTRING_PARAM_ONE: ${MOCK_REQUEST_QUERYSTRING_PARAM_ONE}",
		"MOCK_REQUEST_QUERYSTRING_PARAM_TWO: ${MOCK_REQUEST_QUERYSTRING_PARAM_TWO}",
		"MOCK_REQUEST_HEADER_SOME_HEADER_ONE: ${MOCK_REQUEST_HEADER_SOME_HEADER_ONE}",
		"MOCK_REQUEST_HEADER_SOME_HEADER_TWO: ${MOCK_REQUEST_HEADER_SOME_HEADER_TWO}",
		"MOCK_ROUTE_PARAM_SOME_PARAM: ${MOCK_ROUTE_PARAM_SOME_PARAM}",
	}, "\n")

	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar/{some_param}",
			fmt.Sprintf("--response '%s'", responseStr),
		},
		"GET",
		"foo/bar/test?param_one=value_one&param_two=value_two",
		map[string]string{
			"Some-header-one": "header_value_one",
			"Some-header-two": "header_value_two",
		},
		strings.NewReader(""),
		LineEquals(1, `MOCK_HOST: localhost:{{TEST_E2E_PORT}}`),
		LineEquals(2, `MOCK_REQUEST_HOST: localhost:{{TEST_E2E_PORT}}`),
		LineEquals(3, `MOCK_REQUEST_URL: http://localhost:{{TEST_E2E_PORT}}/foo/bar/test`),
		LineEquals(4, `MOCK_REQUEST_ENDPOINT: foo/bar/test`),
		LineEquals(5, `MOCK_REQUEST_METHOD: get`),
		LineEquals(6, `MOCK_REQUEST_QUERYSTRING: param_one=value_one&param_two=value_two`),
		LineEquals(7, `MOCK_REQUEST_QUERYSTRING_PARAM_ONE: value_one`),
		LineEquals(8, `MOCK_REQUEST_QUERYSTRING_PARAM_TWO: value_two`),
		LineEquals(9, `MOCK_REQUEST_HEADER_SOME_HEADER_ONE: header_value_one`),
		LineEquals(10, `MOCK_REQUEST_HEADER_SOME_HEADER_TWO: header_value_two`),
		LineEquals(11, `MOCK_ROUTE_PARAM_SOME_PARAM: test`),
	)
}

func Test_E2E_Response_PlainText_ReadingRequestBody(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			"--response 'Request payload: ${MOCK_REQUEST_BODY}'",
		},
		"GET",
		"foo/bar",
		nil,
		strings.NewReader("THIS IS THE REQUEST PAYLOAD."),
		StringMatches("Request payload: THIS IS THE REQUEST PAYLOAD."),
	)
}

func Test_E2E_Response_Json_UsingVariables_WithFile(t *testing.T) {
	RunTest(
		t,
		"config_responses_using_variables/config.json",
		"GET",
		"response_json_using_variables/with_file?param_one=value_one&param_two=value_two",
		nil,
		strings.NewReader(""),
		JsonMatches(map[string]interface{}{
			"MOCK_HOST":                          "localhost:{{TEST_E2E_PORT}}",
			"MOCK_REQUEST_HOST":                  "localhost:{{TEST_E2E_PORT}}",
			"MOCK_REQUEST_URL":                   "http://localhost:{{TEST_E2E_PORT}}/response_json_using_variables/with_file",
			"MOCK_REQUEST_ENDPOINT":              "response_json_using_variables/with_file",
			"MOCK_REQUEST_METHOD":                "get",
			"MOCK_REQUEST_QUERYSTRING":           "param_one=value_one&param_two=value_two",
			"MOCK_REQUEST_QUERYSTRING_PARAM_ONE": "value_one",
			"MOCK_REQUEST_QUERYSTRING_PARAM_TWO": "value_two",
		}),
	)
}

func Test_E2E_Response_Json_ReadingRouteParams_WithFile(t *testing.T) {
	RunTest(
		t,
		"config_responses_using_variables/config.json",
		"GET",
		"response_json_reading_route_params/foo/bar",
		nil,
		strings.NewReader(""),
		JsonMatches(map[string]interface{}{
			"var_a": "foo",
			"var_b": "bar",
		}),
	)
}

func Test_E2E_WithNoMethodDefinedDefaultsToGet(t *testing.T) {
	RunTest(
		t,
		"config_with_file_responses/config.json",
		"GET",
		"with_no_method_defined",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(200),
		StringMatches("Hello world."),
	)

	RunTest(
		t,
		"config_with_file_responses/config.json",
		"POST",
		"with_no_method_defined",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(405),
	)

	RunTest(
		t,
		"config_with_file_responses/config.json",
		"PUT",
		"with_no_method_defined",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(405),
	)

	RunTest(
		t,
		"config_with_file_responses/config.json",
		"DELETE",
		"with_no_method_defined",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(405),
	)

	RunTest(
		t,
		"config_with_file_responses/config.json",
		"PATCH",
		"with_no_method_defined",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(405),
	)
}

func Test_E2E_Response_Exec(t *testing.T) {
	RunTest(
		t,
		"config_with_script_responses/config.json",
		"GET",
		"with/exec",
		nil,
		strings.NewReader(""),
		StringMatches("Hello world!"),
	)
}

func Test_E2E_Response_Exec_WithPipe(t *testing.T) {
	RunTest(
		t,
		"config_with_script_responses/config.json",
		"GET",
		"with/exec/with/pipe",
		nil,
		strings.NewReader(""),
		StringMatches("Hello world!"),
	)
}

func Test_E2E_Response_Exec_WithEnvVariable(t *testing.T) {
	RunTestWithEnv(
		t,
		"config_with_script_responses/config.json",
		"GET",
		"with/exec/with/env/var",
		nil,
		strings.NewReader(""),
		map[string]string{
			"FOO": "bar",
		},
		StringMatches("foo: bar"),
	)
}

func Test_E2E_Response_Exec_PrintingEnv(t *testing.T) {
	RunTestWithEnv(
		t,
		"config_with_script_responses/config.json",
		"GET",
		"with/exec/print/env/with/param/bar?foo=bar",
		map[string]string{
			"Header-One": "Header value one",
			"Header-Two": "Header value two",
		},
		strings.NewReader(""),
		map[string]string{
			"FOO": "bar",
		},
		LineEquals(1, `MOCK_HOST=localhost:{{TEST_E2E_PORT}}`),
		LineRegexMatches(2, `MOCK_REQUEST_BODY=.*`),
		LineEquals(3, `MOCK_REQUEST_ENDPOINT=with/exec/print/env/with/param/bar`),
		LineRegexMatches(4, `MOCK_REQUEST_HEADERS=.*`),
		LineEquals(5, `MOCK_REQUEST_HEADER_HEADER_ONE=Header value one`),
		LineEquals(6, `MOCK_REQUEST_HEADER_HEADER_TWO=Header value two`),
		LineEquals(7, `MOCK_REQUEST_HOST=localhost:{{TEST_E2E_PORT}}`),
		LineEquals(8, `MOCK_REQUEST_METHOD=get`),
		LineEquals(9, `MOCK_REQUEST_NTH=1`),
		LineEquals(10, `MOCK_REQUEST_QUERYSTRING=foo=bar`),
		LineEquals(11, `MOCK_REQUEST_QUERYSTRING_FOO=bar`),
		LineEquals(12, `MOCK_REQUEST_URL=http://localhost:{{TEST_E2E_PORT}}/with/exec/print/env/with/param/bar`),
		LineRegexMatches(13, `MOCK_RESPONSE_BODY=.*`),
		LineRegexMatches(14, `MOCK_RESPONSE_HEADERS=.*`),
		LineRegexMatches(15, `MOCK_RESPONSE_STATUS_CODE=.*`),
		LineEquals(16, `MOCK_ROUTE_PARAM_FOO=bar`),
	)
}

func Test_E2E_Response_Exec_PrintingRequestNth(t *testing.T) {
	RunTest(
		t,
		"config_with_script_responses/config.json",
		"GET",
		"print_request_nth",
		nil,
		strings.NewReader(""),
		LineEquals(1, `MOCK_REQUEST_NTH=1`),
	)

	RunTestWithMultipleRequests(
		t,
		"config_with_script_responses/config.json",
		[]TestRequest{
			*NewGetTestRequest("print_request_nth"),
			*NewGetTestRequest("print_request_nth"),
		},
		LineEquals(1, `MOCK_REQUEST_NTH=2`),
	)
}

func Test_E2E_Response_Exec_WithCmdParams(t *testing.T) {
	flagVariations := []string{
		`--exec 'printf "cexe hguorht detareneg saw txet siht" | rev > $MOCK_RESPONSE_BODY'`,
		`--response-exec 'printf "cexe hguorht detareneg saw txet siht" | rev > $MOCK_RESPONSE_BODY'`,
		`--response 'exec:printf "cexe hguorht detareneg saw txet siht" | rev > $MOCK_RESPONSE_BODY'`,
	}

	for _, flagVariation := range flagVariations {
		RunTestWithNoConfigAndWithArgs(
			t,
			[]string{
				"--route foo/bar",
				flagVariation,
			},
			"GET",
			"foo/bar",
			nil,
			strings.NewReader(""),
			StatusCodeMatches(200),
			StringMatches("this text was generated through exec"),
		)
	}
}
