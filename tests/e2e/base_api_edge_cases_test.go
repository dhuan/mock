package tests_e2e

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_BaseApi_RequestForwardedToBaseApi_SupportDifferentFormatsOfUrl(t *testing.T) {
	baseApis := []string{
		"--base localhost:%d",
		"--base localhost:%d/",
		"--base http://localhost:%d",
		"--base http://localhost:%d/",
	}

	state := NewState()
	killMockBase, _, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--status-code 205",
			"--response 'Hello world! This is the base API.'",
		}, " "),
		nil,
		true,
	)
	defer killMockBase()

	for _, baseApi := range baseApis {
		RunTestWithNoConfigAndWithArgs(
			t,
			[]string{
				fmt.Sprintf(baseApi, state.Port),
				"--route hello/world",
				"--response 'Hello world!'",
			},
			"GET",
			"foo/bar",
			nil,
			strings.NewReader(""),
			StatusCodeMatches(205),
			StringMatches("Hello world! This is the base API."),
		)
	}
}

func Test_E2E_BaseApi_CannotBeStartedBecauseBaseApiIsInvalid(t *testing.T) {
	baseApis := []string{
		"....",
		"x",
		"--some-invalid-hostname--",
		"http//invalid",
		"http:invalid",
		"http:///invalid",
		"http://-----------.com",
	}

	for _, baseApi := range baseApis {
		RunTestWithNoConfigAndWithArgsFailing(
			t,
			[]string{
				fmt.Sprintf("--base %s", baseApi),
				"--route hello/world",
				"--response 'Hello world!'",
			},
			"GET",
			"foo/bar",
			nil,
			strings.NewReader(""),
			ApplicationOutputMatches([]string{
				fmt.Sprintf("Base API is not valid: %s", baseApi),
				"Set it as a valid domain name such as google.com",
			}),
		)
	}
}
