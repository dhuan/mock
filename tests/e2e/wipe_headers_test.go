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
