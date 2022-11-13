package tests_e2e

import (
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
		StringMatches(`URL: http://localhost:{{TEST_E2E_PORT}}/foo/bar/2
Endpoint: foo/bar/2
Method: GET
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
