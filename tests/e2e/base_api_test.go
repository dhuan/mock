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

func Test_E2E_BaseApi_RequestForwardedToBaseApi_WithQuerystring(t *testing.T) {
	state := NewState()
	killMockBase, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--response 'Querystring: ${MOCK_REQUEST_QUERYSTRING}'",
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
		"foo/bar?param_one=value_one&param_two=value_two",
		nil,
		strings.NewReader(""),
		StringMatches("Querystring: param_one=value_one&param_two=value_two"),
	)
}

func Test_E2E_BaseApi_RequestForwardedToBaseApi_ResponseHeadersAreForwaded(t *testing.T) {
	state := NewState()
	killMockBase, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--header 'Header-One: value one'",
			"--header 'Header-Two: value two'",
			"--response 'Hello world!'",
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
		HeadersMatch(map[string]string{
			"Header-One": "value one",
			"Header-Two": "value two",
		}),
	)
}

func Test_E2E_BaseApi_RequestForwardedToBaseApi_RequestBodyIsForwarded(t *testing.T) {
	state := NewState()
	killMockBase, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--method POST",
			"--response 'Request body: ${MOCK_REQUEST_BODY}'",
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
		"POST",
		"foo/bar",
		nil,
		strings.NewReader("This is a request payload."),
		StringMatches("Request body: This is a request payload."),
	)
}
