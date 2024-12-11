package tests_e2e

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_GetPayload_AllPayload(t *testing.T) {
	RunTest3(
		t,
		[]string{
			"--route foo/bar",
			"--method POST",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} get-payload | {{MOCK_EXECUTABLE}} write`,
			}, ";")),
		},
		Post("foo/bar", nil, strings.NewReader("Hello, world. This is the payload.")),
		StringMatches("Hello, world. This is the payload."),
	)
}

func Test_E2E_GetPayload_GetJsonField_OK(t *testing.T) {
	RunTest3(
		t,
		[]string{
			"--route foo/bar",
			"--method POST",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} get-payload foo | {{MOCK_EXECUTABLE}} write`,
			}, ";")),
		},
		Post("foo/bar", json_header, strings.NewReader(`{"foo": "bar"}`)),
		StringMatches("bar\n"),
	)
}

func Test_E2E_GetPayload_GetJsonField_FieldDoesNotExist(t *testing.T) {
	RunTest3(
		t,
		[]string{
			"--route foo/bar",
			"--method POST",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} get-payload foo | {{MOCK_EXECUTABLE}} write`,
				`{{MOCK_EXECUTABLE}} get-payload foo`,
				`printf $? | {{MOCK_EXECUTABLE}} write -a`,
			}, ";")),
		},
		Post("foo/bar", json_header, strings.NewReader(`{"hello": "world"}`)),
		StringMatches("1"),
	)
}

func Test_E2E_GetPayload_GetJsonField_WithEmptyPayload_Exit1(t *testing.T) {
	RunTest3(
		t,
		[]string{
			"--route foo/bar",
			"--method POST",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} get-payload foo | {{MOCK_EXECUTABLE}} write`,
				`{{MOCK_EXECUTABLE}} get-payload foo`,
				`printf $? | {{MOCK_EXECUTABLE}} write -a`,
			}, ";")),
		},
		Post("foo/bar", json_header, nil),
		StringMatches("1"),
	)
}
