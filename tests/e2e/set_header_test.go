package tests_e2e

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_SetHeader(t *testing.T) {
	for _, tc := range []struct {
		headerCommand       string
		expectedHeaderKey   string
		expectedHeaderValue string
	}{
		{"foo bar", "Foo", "bar"},
		{`some-header-key "some header value"`, "Some-Header-Key", "some header value"},
	} {
		RunTestWithNoConfigAndWithArgs(
			t,
			[]string{
				"--route foo/bar",
				fmt.Sprintf("--exec '%s'", strings.Join([]string{
					fmt.Sprintf(`{{MOCK_EXECUTABLE}} set-header %s`, tc.headerCommand),
				}, ";")),
			},
			"GET",
			"foo/bar",
			nil,
			nil,
			HeadersMatch(map[string][]string{
				tc.expectedHeaderKey: {tc.expectedHeaderValue},
			}),
		)
	}
}

func Test_E2E_SetHeader_Overwriting(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				fmt.Sprintf(`{{MOCK_EXECUTABLE}} set-header foo bar`),
				fmt.Sprintf(`{{MOCK_EXECUTABLE}} set-header hello world`),
				fmt.Sprintf(`{{MOCK_EXECUTABLE}} set-header foo MODIFIED`),
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		HeadersMatch(map[string][]string{
			"Hello": {"world"},
			"Foo":   {"MODIFIED"},
		}),
	)
}
