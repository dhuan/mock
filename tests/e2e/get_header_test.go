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
			"another-header-key: another header value",
			"some-header-key: some header value",
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

func Test_E2E_GetHeader_Match(t *testing.T) {
	getHeaderTest(
		t,
		[]string{
			`{{MOCK_EXECUTABLE}} get-header some-header-key | {{MOCK_EXECUTABLE}} write`,
			`printf $? >> $MOCK_RESPONSE_BODY`,
		},
		[]string{
			"some-header-key: some header value",
			"0",
		},
	)
}

func Test_E2E_GetHeader_PrintValueOnly(t *testing.T) {
	getHeaderTest(
		t,
		[]string{
			`{{MOCK_EXECUTABLE}} get-header -v some-header-key | {{MOCK_EXECUTABLE}} write`,
		},
		[]string{
			"some header value",
			"",
		},
	)
}

func getHeaderTest(t *testing.T, exec, expectOutput []string) {
	RunTestWithNoConfigAndWithArgs(
		t,
		append(
			[]string{
				"--route foo/bar",
				"--response 'Hello, world!'",
			},
			fmt.Sprintf("--exec '%s'", strings.Join(exec, ";")),
		),
		"GET",
		"foo/bar",
		map[string]string{
			"some-header-key":    "some header value",
			"another-header-key": "another header value",
		},
		nil,
		StringMatches(strings.Join(expectOutput, "\n")),
	)
}
