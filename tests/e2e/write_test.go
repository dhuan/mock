package tests_e2e

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_Write(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`printf "Hello, world! Write was used." | mock write`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		StringMatches("Hello, world! Write was used."),
	)
}

func Test_E2E_Write_Append(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`printf "Hello, world!" | mock write`,
				`printf " Append was used." | mock write -a`,
				`printf " Again." | mock write --append`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		StringMatches("Hello, world! Append was used. Again."),
	)
}

func Test_E2E_Write_WritingMultipleTimesOverwrites(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`printf "Hello, world!" | mock write`,
				`printf "Write again." | mock write`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		StringMatches("Write again."),
	)
}

func Test_E2E_Write_WithJsonOption_Ok(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				// Space excess is intentional in order to assert
				// that --json formats the JSON for us.
				`printf "{\"foo\":          \"bar\"}" | mock write --json`,
				`mock set-header Exit-Status-Code "${?}"`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		StringMatches(`{"foo":"bar"}`),
		HeadersMatch(map[string][]string{
			"Content-Type":     {"application/json"},
			"Exit-Status-Code": {"0"},
		}),
	)
}

func Test_E2E_Write_WithJsonOption_InvalidJson(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`printf "{\"foo\":INVALID}" | mock write --json`,
				`mock set-header Exit-Status-Code "${?}"`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		StringMatches(``),
		HeadersMatch(map[string][]string{
			"Exit-Status-Code": {"1"},
		}),
		HeaderKeysNotIncluded([]string{"Content-Type"}),
	)
}

func Test_E2E_Write_WithJsonOption_CannotUseWithAppend(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`printf "{\"foo\":123}" | mock write --json --append`,
				`mock set-header Exit-Status-Code "${?}"`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		StringMatches(``),
		HeadersMatch(map[string][]string{
			"Exit-Status-Code": {"1"},
		}),
		HeaderKeysNotIncluded([]string{"Content-Type"}),
	)
}
