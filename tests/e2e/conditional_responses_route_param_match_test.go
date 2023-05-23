package tests_e2e

import (
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_ConditionalResponses_RouteParamMatch_DefaultResponse(t *testing.T) {
	requestRoutesThatWillNotMatch := []string{
		"conditional_response/route_param_match/hello/world",
		"conditional_response/route_param_match/foo/wrong",
		"conditional_response/route_param_match/wrong/bar",
		"conditional_response/route_param_match/foo/barr",
	}

	for i := range requestRoutesThatWillNotMatch {
		RunTest(
			t,
			"config_with_conditional_response/config.json",
			"GET",
			requestRoutesThatWillNotMatch[i],
			nil,
			strings.NewReader(""),
			StringMatches("Default response"),
		)
	}
}

func Test_E2E_ConditionalResponses_RouteParamMatch_ConditionalResponseMach(t *testing.T) {
	requestRoutesThatWillMatch := []string{
		"conditional_response/route_param_match/foo/bar",
	}

	for i := range requestRoutesThatWillMatch {
		RunTest(
			t,
			"config_with_conditional_response/config.json",
			"GET",
			requestRoutesThatWillMatch[i],
			nil,
			strings.NewReader(""),
			StringMatches("Conditional response with Route Param Match resolved."),
		)
	}
}
