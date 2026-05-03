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
				`printf "Hello, world." | mock write`,
				`mock replace world WORLD`,
				`mock replace . !`,
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
				`printf "Hello, world." | mock write`,
				`mock replace --regex "w[a-z]{1,}" people`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		StringMatches("Hello, people."),
	)
}

func Test_E2E_Replace_Error_OnlyTwoArgsAreAllowed(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`printf "Hello, world." | mock write`,
				`mock replace one two three`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		StringMatches(`Hello, world.`),
		ApplicationOutputHasLines([]string{
			"Output from program execution:",
			"",
			`"replace" allows only 2 paramaters.`,
		}),
	)
}
