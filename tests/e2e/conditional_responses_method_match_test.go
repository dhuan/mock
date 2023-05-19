package tests_e2e

import (
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_ConditionalResponses_MethodMatch_DefaultResponse(t *testing.T) {
	RunTest(
		t,
		"config_with_conditional_response/config.json",
		"GET",
		"conditional_response/method_match",
		nil,
		strings.NewReader(""),
		StringMatches("Default response"),
	)
}

// Unfortunately the following the test is not possible yet.
/*
func Test_E2E_ConditionalResponses_MethodMatch_ConditionalResponseMatch(t *testing.T) {
	RunTest(
		t,
		"config_with_conditional_response/config.json",
		"POST",
		"conditional_response/method_match",
        nil,
		"",
		StringMatches("Conditional response with Method Match resolved."),
	)
}
*/
