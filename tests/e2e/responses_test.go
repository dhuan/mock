package tests_e2e

import (
	"fmt"
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
		"",
		StringMatches("Hello world! This is response A.\n"),
	)
}

func Test_E2E_Response_ResponseInsideFolder(t *testing.T) {
	RunTest(
		t,
		"config_with_file_responses/config.json",
		"POST",
		"foo/bar/2",
		nil,
		"",
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
		"",
		JsonMatches(map[string]interface{}{
			"response_text": "This is a JSON response.",
		}),
	)
}

func Test_E2E_Response_WithDynamicFileName(t *testing.T) {
	RunTest(
		t,
		"config_with_file_responses/config.json",
		"GET",
		"books/i_robot/content",
		nil,
		"",
		StringMatches("This is the book 'I, Robot'.\n"),
	)

	RunTest(
		t,
		"config_with_file_responses/config.json",
		"GET",
		"books/nightfall/content",
		nil,
		"",
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
		"",
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
		"",
		StringMatches(fmt.Sprintf(`Server Host: localhost:4000
URL: http://localhost:%s/foo/bar/2
Endpoint: foo/bar/2
Method: GET
Querystring: some_key=some_value&another_key=another_value
Headers:
accept-encoding: gzip
another-header-key: Another-Header-Value
some-header-key: Some-Header-Value
user-agent: Go-http-client/1.1`, GetTestPort())),
	)
}

func Test_E2E_Response_ShellScript_RequestDetailsFromEnvVariables_WithPayload(t *testing.T) {
	RunTest(
		t,
		"config_with_script_responses/config.json",
		"POST",
		"foo/bar/3",
		nil,
		"This is the request payload.",
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
		"This is the request payload.",
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
		"This is the request payload.",
		StringMatches(`Parameter: foobar`),
	)
}

func Test_E2E_Response_ShellScript_CommandFailing(t *testing.T) {
	RunTest(
		t,
		"config_with_script_responses/config.json",
		"POST",
		"foo/bar/6",
		nil,
		"This is the request payload.",
		LineEquals(1, `Hello world!`),
		LineRegexMatches(2, `tests/e2e/data/config_with_script_responses/handler_with_command_that_fails.sh:.*3: please_fail:.* not found$`),
	)
}

func Test_E2E_Response_ShellScript_ReadingEndpointParams(t *testing.T) {
	RunTest(
		t,
		"config_with_script_responses/config.json",
		"GET",
		"users/country/brazil/page/7",
		nil,
		"",
		LineEquals(1, `Country: brazil`),
		LineEquals(2, `Page: 7`),
	)
}

func Test_E2E_Response_Json_UsingVariables(t *testing.T) {
	RunTest(
		t,
		"config_responses_using_variables/config.json",
		"GET",
		"response_json_using_variables",
		nil,
		"",
		JsonMatches(map[string]interface{}{
			"MOCK_HOST":                fmt.Sprintf("localhost:%s", GetTestPort()),
			"MOCK_REQUEST_URL":         "http://localhost:4000/response_json_using_variables",
			"MOCK_REQUEST_ENDPOINT":    "response_json_using_variables",
			"MOCK_REQUEST_METHOD":      "GET",
			"MOCK_REQUEST_QUERYSTRING": "",
		}),
	)
}

func Test_E2E_Response_Json_UsingVariables_WithFile(t *testing.T) {
	RunTest(
		t,
		"config_responses_using_variables/config.json",
		"GET",
		"response_json_using_variables/with_file",
		nil,
		"",
		JsonMatches(map[string]interface{}{
			"MOCK_HOST":                fmt.Sprintf("localhost:%s", GetTestPort()),
			"MOCK_REQUEST_URL":         "http://localhost:4000/response_json_using_variables/with_file",
			"MOCK_REQUEST_ENDPOINT":    "response_json_using_variables/with_file",
			"MOCK_REQUEST_METHOD":      "GET",
			"MOCK_REQUEST_QUERYSTRING": "",
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
		"",
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
		"",
		StatusCodeMatches(200),
		StringMatches("Hello world."),
	)

	RunTest(
		t,
		"config_with_file_responses/config.json",
		"POST",
		"with_no_method_defined",
		nil,
		"",
		StatusCodeMatches(405),
	)

	RunTest(
		t,
		"config_with_file_responses/config.json",
		"PUT",
		"with_no_method_defined",
		nil,
		"",
		StatusCodeMatches(405),
	)

	RunTest(
		t,
		"config_with_file_responses/config.json",
		"DELETE",
		"with_no_method_defined",
		nil,
		"",
		StatusCodeMatches(405),
	)

	RunTest(
		t,
		"config_with_file_responses/config.json",
		"PATCH",
		"with_no_method_defined",
		nil,
		"",
		StatusCodeMatches(405),
	)
}
