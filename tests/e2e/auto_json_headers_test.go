package tests_e2e

import (
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_AutomaticallyAddJsonHeaders_IfResponseIsJsonObject(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route hello/world",
			"--method get",
			`--response '{"hello":"world"}'`,
		},
		"GET",
		"hello/world",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(200),
		HeadersMatch(map[string][]string{
			"Content-Type": {"application/json"},
		}),
	)
}

func Test_E2E_AutomaticallyAddJsonHeaders_IfResponseIsJsonArray(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route hello/world",
			"--method get",
			`--response '[{"hello":"world"}]'`,
		},
		"GET",
		"hello/world",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(200),
		HeadersMatch(map[string][]string{
			"Content-Type": {"application/json"},
		}),
	)
}
