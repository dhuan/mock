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
		StringMatches("Hello world! This is response A.\n"),
	)
}

func Test_E2E_Response_ResponseInsideFolder(t *testing.T) {
	RunTest(
		t,
		"config_with_file_responses/config.json",
		"POST",
		"foo/bar/2",
		StringMatches("This test asserts that you can set response files inside folders.\n"),
	)
}

func Test_E2E_Response_JsonResponse(t *testing.T) {
	RunTest(
		t,
		"config_with_file_responses/config.json",
		"POST",
		"foo/bar/3",
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
		StringMatches("Hello world! This response was generated from a shell script.\n"),
	)
}
