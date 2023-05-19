package tests_e2e

import (
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_ConditionalResponses_JsonBodyMatch_DefaultResponse(t *testing.T) {
	requestBodiesThatWillNotMatch := []string{
		``,
		`{"fooo":"bar"}`,
		`{"foo":"baar"}`,
		`{"foo":"bar","hello":"world"}`,
	}

	for i := range requestBodiesThatWillNotMatch {
		RunTest(
			t,
			"config_with_conditional_response/config.json",
			"GET",
			"conditional_response/json_body_match",
			nil,
			requestBodiesThatWillNotMatch[i],
			StringMatches("Default response"),
		)
	}
}

func Test_E2E_ConditionalResponses_JsonBodyMatch_ConditionalResponseMatch(t *testing.T) {
	requestBodiesThatWillMatch := []string{
		`{"foo":"bar"}`,
	}

	for i := range requestBodiesThatWillMatch {
		RunTest(
			t,
			"config_with_conditional_response/config.json",
			"GET",
			"conditional_response/json_body_match",
			nil,
			requestBodiesThatWillMatch[i],
			StringMatches("Conditional response with Json Body Match resolved."),
		)
	}
}

func Test_E2E_ConditionalResponses_JsonBodyMatch_ConditionalResponseMatch_WithMultipleFields(t *testing.T) {
	requestBodiesThatWillMatch := []string{
		`{"some_key":"some_value","another_key":"another_value"}`,
	}

	for i := range requestBodiesThatWillMatch {
		RunTest(
			t,
			"config_with_conditional_response/config.json",
			"GET",
			"conditional_response/json_body_match",
			nil,
			requestBodiesThatWillMatch[i],
			StringMatches("Conditional response with Json Body Match resolved - with multiple fields."),
		)
	}
}
