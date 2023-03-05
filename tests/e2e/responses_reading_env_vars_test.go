package tests_e2e

import (
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
		"",
		map[string]string{
			"FOO": "BAR",
		},
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
		"",
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
		"",
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
		"",
		map[string]string{
			"FOO": "BAR",
		},
		JsonMatches(map[string]interface{}{
			"FOO": "BAR",
		}),
	)
}