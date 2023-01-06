package tests_e2e

import (
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_ConditionalResponses_ReceivingDefaultResponse(t *testing.T) {
	RunTest(
		t,
		"config_with_conditional_response/config.json",
		"POST",
		"foo/bar",
		nil,
		"",
		StringMatches("This is the default response."),
	)
}

func Test_E2E_ConditionalResponses_ReceivingConditionalResponse(t *testing.T) {
	RunTest(
		t,
		"config_with_conditional_response/config.json",
		"POST",
		"foo/bar?key1=value1&key2=value2",
		nil,
		"",
		StringMatches("Hello world!"),
	)
}

func Test_E2E_ConditionalResponses_ReceivingConditionalResponse_WithAndChaining(t *testing.T) {
	RunTest(
		t,
		"config_with_conditional_response/config.json",
		"POST",
		"foo/bar?key1=value1&key2=value2&key4=value4",
		nil,
		"",
		StringMatches("Hello world! (Condition with AND chaining)"),
	)
}

func Test_E2E_ConditionalResponses_ReceivingConditionalResponse_WithOrChaining(t *testing.T) {
	RunTest(
		t,
		"config_with_conditional_response/config.json",
		"POST",
		"foo/bar?key1=value1&key6=value6",
		nil,
		"",
		StringMatches("Hello world! (Condition with OR chaining)"),
	)
}
