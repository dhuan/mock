package tests_e2e

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_WipeHeaders_WipeOne(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`printf "Header-One: value one\n" >> $MOCK_RESPONSE_HEADERS`,
				`printf "Header-Two: value two\n" >> $MOCK_RESPONSE_HEADERS`,
				`printf "Header-Three: value three\n" >> $MOCK_RESPONSE_HEADERS`,
				`{{MOCK_EXECUTABLE}} wipe-headers header-two`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		HeaderKeysNotIncluded([]string{"Header-Two"}),
		HeadersMatch(map[string][]string{
			"Header-One":   {"value one"},
			"Header-Three": {"value three"},
		}),
	)
}

func Test_E2E_WipeHeaders_WipeMultiple(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`printf "Header-One: value one\n" >> $MOCK_RESPONSE_HEADERS`,
				`printf "Header-Two: value two\n" >> $MOCK_RESPONSE_HEADERS`,
				`printf "Header-Three: value three\n" >> $MOCK_RESPONSE_HEADERS`,
				`{{MOCK_EXECUTABLE}} wipe-headers header-one header-two`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		HeaderKeysNotIncluded([]string{"Header-One", "Header-Two"}),
		HeadersMatch(map[string][]string{
			"Header-Three": {"value three"},
		}),
	)
}

func Test_E2E_WipeHeaders_WipeNone(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`printf "Header-One: value one\n" >> $MOCK_RESPONSE_HEADERS`,
				`printf "Header-Two: value two\n" >> $MOCK_RESPONSE_HEADERS`,
				`printf "Header-Three: value three\n" >> $MOCK_RESPONSE_HEADERS`,
				`{{MOCK_EXECUTABLE}} wipe-headers header`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		HeadersMatch(map[string][]string{
			"Header-One":   {"value one"},
			"Header-Two":   {"value two"},
			"Header-Three": {"value three"},
		}),
	)
}

func Test_E2E_WipeHeaders_Regex(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`printf "Header-One: value one\n" >> $MOCK_RESPONSE_HEADERS`,
				`printf "Header-Two: value two\n" >> $MOCK_RESPONSE_HEADERS`,
				`printf "Header-Three: value three\n" >> $MOCK_RESPONSE_HEADERS`,
				`{{MOCK_EXECUTABLE}} wipe-headers --regex t.ree`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		HeaderKeysNotIncluded([]string{"Header-Three"}),
		HeadersMatch(map[string][]string{
			"Header-One": {"value one"},
			"Header-Two": {"value two"},
		}),
	)
}

func Test_E2E_WipeHeaders_Regex_Many(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`printf "Header-One: value one\n" >> $MOCK_RESPONSE_HEADERS`,
				`printf "Header-Two: value two\n" >> $MOCK_RESPONSE_HEADERS`,
				`printf "Header-Three: value three\n" >> $MOCK_RESPONSE_HEADERS`,
				`{{MOCK_EXECUTABLE}} wipe-headers --regex foo bar "^header-.*$" hello world`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		HeaderKeysNotIncluded([]string{"Header-One", "Header-Two", "Header-Three"}),
	)
}

func Test_E2E_WipeHeaders_Regex_NoMatch(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`printf "Header-One: value one\n" >> $MOCK_RESPONSE_HEADERS`,
				`printf "Header-Two: value two\n" >> $MOCK_RESPONSE_HEADERS`,
				`printf "Header-Three: value three\n" >> $MOCK_RESPONSE_HEADERS`,
				`{{MOCK_EXECUTABLE}} wipe-headers --regex foo bar "^header$" hello world`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		HeadersMatch(map[string][]string{
			"Header-One":   {"value one"},
			"Header-Two":   {"value two"},
			"Header-Three": {"value three"},
		}),
	)
}
