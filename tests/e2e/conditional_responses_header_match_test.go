package tests_e2e

import (
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_ConditionalResponses_HeaderMatch_DefaultResponse(t *testing.T) {
	headersThatWillNotMatch := []map[string]string{
		nil,
		{"some_header": "some_value"},
		{"some_header": "some_value", "another_header": "another_value"},
		{"foo": "baar"},
		{"fooo": "bar"},
	}

	for i := range headersThatWillNotMatch {
		RunTest(
			t,
			"config_with_conditional_response/config.json",
			"GET",
			"conditional_response/header_match",
			headersThatWillNotMatch[i],
			strings.NewReader(""),
			StringMatches("Default response"),
		)
	}
}

func Test_E2E_ConditionalResponses_HeaderMatch_MatchConditionalResponse(t *testing.T) {
	headersThatWillMatch := []map[string]string{
		{"foo": "bar"},
		{"hello": "world", "foo": "bar"},
	}

	for i := range headersThatWillMatch {
		RunTest(
			t,
			"config_with_conditional_response/config.json",
			"GET",
			"conditional_response/header_match",
			headersThatWillMatch[i],
			strings.NewReader(""),
			StringMatches("Conditional response with Header Match resolved."),
		)
	}
}
