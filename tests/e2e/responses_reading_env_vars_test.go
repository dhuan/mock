package tests_e2e

import (
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_Response_ReadingEnvironmentVariable_TextResponse(t *testing.T) {
	RunTestWithEnv(
		t,
		"config_responses_using_variables/config.json",
		"GET",
		"reading/env/vars/text",
		nil,
		strings.NewReader(""),
		map[string]string{
			"FOO": "BAR",
		},
		StringMatches("The value of FOO is: BAR."),
	)
}

func Test_E2E_Response_ReadingEnvironmentVariable_TextResponse_WithCmdParams(t *testing.T) {
	RunTestWithArgsAndEnv(
		t,
		[]string{
			"--route reading/env/vars/text",
			"--response 'The value of FOO is: ${FOO}.'",
		},
		"GET",
		"reading/env/vars/text",
		nil,
		strings.NewReader(""),
		map[string]string{
			"FOO": "BAR",
		},
		StatusCodeMatches(200),
		StringMatches("The value of FOO is: BAR."),
	)
}

func Test_E2E_Response_ReadingEnvironmentVariable_TextFileResponse(t *testing.T) {
	RunTestWithEnv(
		t,
		"config_responses_using_variables/config.json",
		"GET",
		"reading/env/vars/text_file",
		nil,
		strings.NewReader(""),
		map[string]string{
			"FOO": "BAR",
		},
		StringMatches("The value of FOO is: BAR.\n"),
	)
}

func Test_E2E_Response_ReadingEnvironmentVariable_JsonResponse(t *testing.T) {
	RunTestWithEnv(
		t,
		"config_responses_using_variables/config.json",
		"GET",
		"reading/env/vars/json",
		nil,
		strings.NewReader(""),
		map[string]string{
			"FOO": "BAR",
		},
		JsonMatches(map[string]interface{}{
			"FOO": "BAR",
		}),
	)
}

func Test_E2E_Response_ReadingEnvironmentVariable_JsonFileResponse(t *testing.T) {
	RunTestWithEnv(
		t,
		"config_responses_using_variables/config.json",
		"GET",
		"reading/env/vars/json_file",
		nil,
		strings.NewReader(""),
		map[string]string{
			"FOO": "BAR",
		},
		JsonMatches(map[string]interface{}{
			"FOO": "BAR",
		}),
	)
}
