package tests_e2e

import (
	"io"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_ConditionalResponses_FormMatch_DefaultResponse(t *testing.T) {
	requestBodiesThatWillNotMatch := []io.Reader{
		strings.NewReader(""),
		BuildFormPayload(map[string]string{
			"some_key": "some_invalid_value",
		}),
	}

	for i := range requestBodiesThatWillNotMatch {
		RunTest(
			t,
			"config_with_conditional_response/config.json",
			"GET",
			"conditional_response/form_match",
			nil,
			requestBodiesThatWillNotMatch[i],
			StringMatches("Default response"),
		)
	}
}

func Test_E2E_ConditionalResponses_FormMatch_ConditionalResponseMatch(t *testing.T) {
	requestBodiesThatWillMatch := []io.Reader{
		BuildFormPayload(map[string]string{
			"some_key": "some_value",
		}),
		BuildFormPayload(map[string]string{
			"some_key":    "some_value",
			"another_key": "another_value",
		}),
	}

	for i := range requestBodiesThatWillMatch {
		RunTest(
			t,
			"config_with_conditional_response/config.json",
			"GET",
			"conditional_response/form_match",
			nil,
			requestBodiesThatWillMatch[i],
			StringMatches("Conditional response with Form Match resolved."),
		)
	}
}
