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
		StringMatches("Hello world! This is response A.\n"),
	)
}
