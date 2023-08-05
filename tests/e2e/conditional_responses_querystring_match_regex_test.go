package tests_e2e

import (
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_ConditionalResponses_QuerystringMatchRegex_DefaultResponse(t *testing.T) {
	requestsThatWillNotMatch := []string{
		"conditional_response/querystring_match_regex",
		"conditional_response/querystring_match_regex?foo=123",
		"conditional_response/querystring_match_regex?foo=baar",
	}

	for i := range requestsThatWillNotMatch {
		RunTest(
			t,
			"config_with_conditional_response/config.json",
			"GET",
			requestsThatWillNotMatch[i],
			nil,
			strings.NewReader(""),
			StringMatches("Default response"),
		)
	}
}

func Test_E2E_ConditionalResponses_QuerystringMatchRegex_ConditionalResponseMatching(t *testing.T) {
	requestsThatWillNotMatch := []string{
		"conditional_response/querystring_match_regex?foo=bar",
		"conditional_response/querystring_match_regex?foo=abc",
	}

	for i := range requestsThatWillNotMatch {
		RunTest(
			t,
			"config_with_conditional_response/config.json",
			"GET",
			requestsThatWillNotMatch[i],
			nil,
			strings.NewReader(""),
			StringMatches("Conditional response with Querystring Match Regex resolved."),
		)
	}
}
