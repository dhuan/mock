package tests_e2e

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_GetHeader_All(t *testing.T) {
	getHeaderTest(
		t,
		[]string{
			`{{MOCK_EXECUTABLE}} get-header | {{MOCK_EXECUTABLE}} write`,
			`printf $? >> $MOCK_RESPONSE_BODY`,
		},
		[]string{
			"accept-encoding: gzip",
			"user-agent: Go-http-client/1.1",
			"0",
		},
	)
}

func Test_E2E_GetHeader_NoMatches(t *testing.T) {
	getHeaderTest(
		t,
		[]string{
			`{{MOCK_EXECUTABLE}} get-header foobar | {{MOCK_EXECUTABLE}} write`,
		},
		[]string{
			"",
		},
	)
}

func Test_E2E_GetHeader_NoMatches_ExitCode1(t *testing.T) {
	getHeaderTest(
		t,
		[]string{
			`{{MOCK_EXECUTABLE}} get-header foobar`,
			`printf $? >> $MOCK_RESPONSE_BODY`,
		},
		[]string{
			"1",
		},
	)
}

func getHeaderTest(t *testing.T, exec, expectOutput []string) {
	RunTestWithNoConfigAndWithArgs(
		t,
		append(
			[]string{
				"--route foo/bar",
				"--header 'some-header-key: some header value'",
				"--header 'another-header-key: another header value'",
			},
			fmt.Sprintf("--exec '%s'", strings.Join(exec, ";")),
		),
		"GET",
		"foo/bar",
		nil,
		nil,
		StringMatches(strings.Join(expectOutput, "\n")),
	)
}
