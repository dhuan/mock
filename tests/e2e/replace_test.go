package tests_e2e

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_Replace(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`printf "Hello, world." | {{MOCK_EXECUTABLE}} write`,
				`{{MOCK_EXECUTABLE}} replace world WORLD`,
				`{{MOCK_EXECUTABLE}} replace . !`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		StringMatches("Hello, WORLD!"),
	)
}

func Test_E2E_Replace_WithRegex(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`printf "Hello, world." | {{MOCK_EXECUTABLE}} write`,
				`{{MOCK_EXECUTABLE}} replace --regex "w[a-z]{1,}" people`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		StringMatches("Hello, people."),
	)
}
