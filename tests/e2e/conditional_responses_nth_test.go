package tests_e2e

import (
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_ConditionalResponses_Nth_FirstRequest(t *testing.T) {
	nthEndpoints := []string{
		"conditional_response/nth",
		"conditional_response/nth/with_numbers",
		"conditional_response/nth/with_param/some_value",
	}

	for _, nthEndpoint := range nthEndpoints {
		RunTest(
			t,
			"config_with_conditional_response/config.json",
			"GET",
			nthEndpoint,
			nil,
			strings.NewReader(""),
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
				*NewGetTestRequest(nthEndpoint),
			},
			StringMatches("Default response"),
		)
	}
}

func Test_E2E_ConditionalResponses_Nth_SecondRequest(t *testing.T) {
	nthEndpoints := []string{
		"conditional_response/nth",
		"conditional_response/nth/with_numbers",
		"conditional_response/nth/with_param/some_value",
	}

	for _, nthEndpoint := range nthEndpoints {
		RunTestWithMultipleRequests(
			t,
			"config_with_conditional_response/config.json",
			[]TestRequest{
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest(nthEndpoint),
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest(nthEndpoint),
			},
			StringMatches("Second response"),
		)
	}
}

func Test_E2E_ConditionalResponses_Nth_ThirdRequest(t *testing.T) {
	nthEndpoints := []string{
		"conditional_response/nth",
		"conditional_response/nth/with_numbers",
		"conditional_response/nth/with_param/some_value",
	}

	for _, nthEndpoint := range nthEndpoints {
		RunTestWithMultipleRequests(
			t,
			"config_with_conditional_response/config.json",
			[]TestRequest{
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest(nthEndpoint),
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest(nthEndpoint),
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest(nthEndpoint),
			},
			StringMatches("Third response"),
		)
	}
}

func Test_E2E_ConditionalResponses_Nth_FourthRequestFallsbackToDefault(t *testing.T) {
	nthEndpoints := []string{
		"conditional_response/nth",
		"conditional_response/nth/with_numbers",
		"conditional_response/nth/with_param/some_value",
	}

	for _, nthEndpoint := range nthEndpoints {
		RunTestWithMultipleRequests(
			t,
			"config_with_conditional_response/config.json",
			[]TestRequest{
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest(nthEndpoint),
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest(nthEndpoint),
				*NewGetTestRequest("dummy_endpoint"),
				*NewGetTestRequest(nthEndpoint),
				*NewGetTestRequest(nthEndpoint),
			},
			StringMatches("Default response"),
		)
	}
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
