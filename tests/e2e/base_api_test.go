package tests_e2e

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_BaseApi_NormalRequest(t *testing.T) {
	state := NewState()
	killMockBase, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--status-code 205",
			"--response 'Hello world! This is the base API.'",
		}, " "),
		nil,
	)
	defer killMockBase()

	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			fmt.Sprintf("--base 'localhost:%d'", state.Port),
			"--route hello/world",
			"--response 'Hello world!'",
		},
		"GET",
		"hello/world",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(200),
		StringMatches("Hello world!"),
	)
}

func Test_E2E_BaseApi_RequestForwardedToBaseApi(t *testing.T) {
	state := NewState()
	killMockBase, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--status-code 205",
			"--response 'Hello world! This is the base API.'",
		}, " "),
		nil,
	)
	defer killMockBase()

	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			fmt.Sprintf("--base 'localhost:%d'", state.Port),
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
