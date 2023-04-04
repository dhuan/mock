package tests_e2e

import (
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_CommandLineDefinedEndpoints_WithOneEndpoint(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route hello/world",
			"--method get",
			"--response 'Hello world!'",
		},
		"GET",
		"hello/world",
		nil,
		"",
		StatusCodeMatches(200),
		StringMatches("Hello world!"),
	)
}

func Test_E2E_CommandLineDefinedEndpoints_WithMultipleEndpoints(t *testing.T) {
	commandArgs := []string{
		"--route endpoint/one",
		"--method get",
		"--response 'First endpoint.'",
		"--route endpoint/two",
		"--method post",
		"--response 'Second endpoint.'",
	}

	RunTestWithNoConfigAndWithArgs(
		t,
		commandArgs,
		"GET",
		"endpoint/one",
		nil,
		"",
		StatusCodeMatches(200),
		StringMatches("First endpoint."),
	)

	RunTestWithNoConfigAndWithArgs(
		t,
		commandArgs,
		"POST",
		"endpoint/two",
		nil,
		"",
		StatusCodeMatches(200),
		StringMatches("Second endpoint."),
	)

	RunTestWithNoConfigAndWithArgs(
		t,
		commandArgs,
		"GET",
		"endpoint/two",
		nil,
		"",
		StatusCodeMatches(405),
	)
}

func Test_E2E_CommandLineDefinedEndpoints_WithoutMethodDefaultsToGet(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route hello/world",
			"--response 'Hello world!'",
		},
		"GET",
		"hello/world",
		nil,
		"",
		StatusCodeMatches(200),
		StringMatches("Hello world!"),
	)
}

func Test_E2E_CommandLineDefinedEndpoints_WithConfigAndArgs(t *testing.T) {
	RunTestWithArgs(
		t,
		"config_basic/config.json",
		[]string{
			"--route hello/world",
			"--response 'Hello world!'",
		},
		"GET",
		"hello/world",
		nil,
		"",
		StatusCodeMatches(200),
		StringMatches("Hello world!"),
	)
}
