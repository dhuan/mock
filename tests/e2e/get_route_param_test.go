package tests_e2e

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_GetRouteParam(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route say_hi/{name}/{location}",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`printf "Hi. My name is $({{MOCK_EXECUTABLE}} get-route-param name). I live on $({{MOCK_EXECUTABLE}} get-route-param location)." | {{MOCK_EXECUTABLE}} write`,
			}, ";")),
		},
		"GET",
		"say_hi/john_doe/earth",
		nil,
		nil,
		StringMatches("Hi. My name is john_doe. I live on earth."),
	)
}

func Test_E2E_GetRouteParam_WithUnexistingKeyExitsWith1(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route say_hi/{name}",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} get-route-param this-param-does-not-exist`,
				`printf "Exit code is %d". "${?}" | mock write`,
			}, ";")),
		},
		"GET",
		"say_hi/john_doe",
		nil,
		nil,
		StringMatches("Exit code is 1."),
	)
}
