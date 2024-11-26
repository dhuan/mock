package tests_e2e

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_GetQuery(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} get-query someKey | {{MOCK_EXECUTABLE}} write`,
			}, ";")),
		},
		"GET",
		"foo/bar?someKey=someValue",
		nil,
		nil,
		StringMatches("someValue"),
	)
}

func Test_E2E_GetQuery_GetAll(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} get-query | {{MOCK_EXECUTABLE}} write`,
			}, ";")),
		},
		"GET",
		"foo/bar?anotherKey=anotherValue&someKey=someValue",
		nil,
		nil,
		StringMatches("anotherKey=anotherValue&someKey=someValue"),
	)
}

func Test_E2E_GetQuery_ExitCode1WhenKeyDoesNotExist(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} get-query someKey ; printf $? | {{MOCK_EXECUTABLE}} write`,
			}, ";")),
		},
		"GET",
		"foo/bar?foo=bar",
		nil,
		nil,
		StringMatches("1"),
	)
}
