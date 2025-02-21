package tests_e2e

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_SetStatus(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} set-status 210`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		StatusCodeMatches(210),
	)
}

func Test_E2E_SetStatus_Error_InvalidStatusCodeParameter(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} set-status invalid`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		ApplicationOutputHasLines([]string{
			"Output from program execution:",
			"",
			`Invalid status code!`,
		}),
	)
}
