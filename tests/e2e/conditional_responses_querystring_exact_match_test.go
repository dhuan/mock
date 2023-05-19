package tests_e2e

import (
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_ConditionalResponses_QuerystringExactMatch_DefaultResponse(t *testing.T) {
	requestsThatWillNotMatch := []string{
		"conditional_response/querystring_exact_match",
		"conditional_response/querystring_exact_match?foo=wrong_value",
		"conditional_response/querystring_exact_match?foo=bar&some_key=some_value",
	}

	for i := range requestsThatWillNotMatch {
		RunTest(
			t,
			"config_with_conditional_response/config.json",
			"GET",
			requestsThatWillNotMatch[i],
			nil,
			"",
			StringMatches("Default response"),
		)
	}
}

func Test_E2E_ConditionalResponses_QuerystringExactMatch_ConditionalResponseMatch(t *testing.T) {
	requestsThatWillNotMatch := []string{
		"conditional_response/querystring_exact_match?foo=bar",
	}

	for i := range requestsThatWillNotMatch {
		RunTest(
			t,
			"config_with_conditional_response/config.json",
			"GET",
			requestsThatWillNotMatch[i],
			nil,
			"",
			StringMatches("Conditional response with Querystring Exact Match resolved."),
		)
	}
}
