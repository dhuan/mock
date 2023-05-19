package tests_e2e

import (
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_ConditionalResponses_QuerystringMatch_DefaultResponse(t *testing.T) {
	requestsThatWillNotMatch := []string{
		"conditional_response/querystring_match",
		"conditional_response/querystring_match?foo=some_value",
		"conditional_response/querystring_match?foo=baar",
		"conditional_response/querystring_match?some_key=some_value&another_key=another_value",
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

func Test_E2E_ConditionalResponses_QuerystringMatch_ConditionalResponseMatching(t *testing.T) {
	requestsThatWillNotMatch := []string{
		"conditional_response/querystring_match?foo=bar",
		"conditional_response/querystring_match?some_key=some_value&foo=bar",
	}

	for i := range requestsThatWillNotMatch {
		RunTest(
			t,
			"config_with_conditional_response/config.json",
			"GET",
			requestsThatWillNotMatch[i],
			nil,
			"",
			StringMatches("Conditional response with Querystring Match resolved."),
		)
	}
}
