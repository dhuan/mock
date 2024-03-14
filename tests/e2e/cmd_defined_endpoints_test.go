package tests_e2e

import (
	"strings"
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
		strings.NewReader(""),
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
		strings.NewReader(""),
		StatusCodeMatches(200),
		StringMatches("First endpoint."),
	)

	RunTestWithNoConfigAndWithArgs(
		t,
		commandArgs,
		"POST",
		"endpoint/two",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(200),
		StringMatches("Second endpoint."),
	)

	RunTestWithNoConfigAndWithArgs(
		t,
		commandArgs,
		"GET",
		"endpoint/two",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(405),
	)
}

func Test_E2E_CommandLineDefinedEndpoints_WithStatusCode(t *testing.T) {
	commandArgs := []string{
		"--route endpoint/one",
		"--status-code 201",
		"--response 'First endpoint.'",
		"--route endpoint/two",
		"--status-code 202",
		"--response 'Second endpoint.'",
	}

	RunTestWithNoConfigAndWithArgs(
		t,
		commandArgs,
		"GET",
		"endpoint/one",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(201),
		StringMatches("First endpoint."),
	)

	RunTestWithNoConfigAndWithArgs(
		t,
		commandArgs,
		"GET",
		"endpoint/two",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(202),
		StringMatches("Second endpoint."),
	)
}

func Test_E2E_CommandLineDefinedEndpoints_WithHeaders(t *testing.T) {
	commandArgs := []string{
		"--route endpoint/one",
		"--response 'First endpoint.'",
		"--route endpoint/two",
		"--header 'Header-One: 1st header'",
		"--header 'Header-Two: 2nd header'",
		"--response 'Second endpoint.'",
		"--route endpoint/three",
		"--response 'Third endpoint.'",
		"--header 'Header-Three: 3rd header'",
	}

	RunTestWithNoConfigAndWithArgs(
		t,
		commandArgs,
		"GET",
		"endpoint/one",
		nil,
		strings.NewReader(""),
		StringMatches("First endpoint."),
		HeaderKeysNotIncluded([]string{
			"Header-One",
			"Header-Two",
			"Header-Three",
		}),
	)

	RunTestWithNoConfigAndWithArgs(
		t,
		commandArgs,
		"GET",
		"endpoint/two",
		nil,
		strings.NewReader(""),
		StringMatches("Second endpoint."),
		HeadersMatch(map[string][]string{
			"Header-One": []string{"1st header"},
			"Header-Two": []string{"2nd header"},
		}),
	)

	RunTestWithNoConfigAndWithArgs(
		t,
		commandArgs,
		"GET",
		"endpoint/three",
		nil,
		strings.NewReader(""),
		StringMatches("Third endpoint."),
		HeadersMatch(map[string][]string{
			"Header-Three": []string{"3rd header"},
		}),
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
		strings.NewReader(""),
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
		strings.NewReader(""),
		StatusCodeMatches(200),
		StringMatches("Hello world!"),
	)
}
