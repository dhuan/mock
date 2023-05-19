package tests_e2e

import (
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_ConditionalResponses_ReceivingDefaultResponse(t *testing.T) {
	RunTest(
		t,
		"config_with_conditional_response/config.json",
		"POST",
		"with_different_responses_based_on_querystring",
		nil,
		strings.NewReader(""),
		StringMatches("This is the default response."),
	)
}

func Test_E2E_ConditionalResponses_ReceivingConditionalResponse(t *testing.T) {
	RunTest(
		t,
		"config_with_conditional_response/config.json",
		"POST",
		"with_different_responses_based_on_querystring?key1=value1&key2=value2",
		nil,
		strings.NewReader(""),
		StringMatches("Hello world!"),
	)
}

func Test_E2E_ConditionalResponses_ReceivingConditionalResponse_WithAndChaining(t *testing.T) {
	RunTest(
		t,
		"config_with_conditional_response/config.json",
		"POST",
		"with_different_responses_based_on_querystring?key1=value1&key2=value2&key4=value4",
		nil,
		strings.NewReader(""),
		StringMatches("Hello world! (Condition with AND chaining)"),
	)
}

func Test_E2E_ConditionalResponses_ReceivingConditionalResponse_WithOrChaining(t *testing.T) {
	RunTest(
		t,
		"config_with_conditional_response/config.json",
		"POST",
		"with_different_responses_based_on_querystring?key1=value1&key6=value6",
		nil,
		strings.NewReader(""),
		StringMatches("Hello world! (Condition with OR chaining)"),
	)
}
