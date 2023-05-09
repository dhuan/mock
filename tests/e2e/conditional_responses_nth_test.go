package tests_e2e

import (
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_ConditionalResponses_Nth_FirstRequest(t *testing.T) {
	RunTest(
		t,
		"config_with_conditional_response/config.json",
		"GET",
		"conditional_response/nth",
		nil,
		"",
		StringMatches("Default response"),
	)

	RunTestWithMultipleRequests(
		t,
		"config_with_conditional_response/config.json",
		[]TestRequest{
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewPostTestRequest("dummy_endpoint"),
			*NewGetTestRequest("conditional_response/nth"),
		},
		StringMatches("Default response"),
	)
}

func Test_E2E_ConditionalResponses_Nth_SecondRequest(t *testing.T) {
	RunTestWithMultipleRequests(
		t,
		"config_with_conditional_response/config.json",
		[]TestRequest{
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("conditional_response/nth"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("conditional_response/nth"),
		},
		StringMatches("Second response"),
	)
}

func Test_E2E_ConditionalResponses_Nth_ThirdRequest(t *testing.T) {
	RunTestWithMultipleRequests(
		t,
		"config_with_conditional_response/config.json",
		[]TestRequest{
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("conditional_response/nth"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("conditional_response/nth"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("conditional_response/nth"),
		},
		StringMatches("Third response"),
	)
}

func Test_E2E_ConditionalResponses_Nth_FourthRequestFallsbackToDefault(t *testing.T) {
	RunTestWithMultipleRequests(
		t,
		"config_with_conditional_response/config.json",
		[]TestRequest{
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("conditional_response/nth"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("conditional_response/nth"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("conditional_response/nth"),
			*NewGetTestRequest("conditional_response/nth"),
		},
		StringMatches("Default response"),
	)
}

func Test_E2E_ConditionalResponses_Nth_WithPlus(t *testing.T) {
	RunTestWithMultipleRequests(
		t,
		"config_with_conditional_response/config.json",
		[]TestRequest{
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("conditional_response/nth/with_plus"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("conditional_response/nth/with_plus"),
			*NewGetTestRequest("dummy_endpoint"),
			*NewGetTestRequest("conditional_response/nth/with_plus"),
			*NewGetTestRequest("conditional_response/nth/with_plus"),
		},
		StringMatches("Second response"),
	)
}
