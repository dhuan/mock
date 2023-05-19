package tests_e2e

import (
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_RequestWithWildcard(t *testing.T) {
	RunTest(
		t,
		"config_with_wildcards/config.json",
		"GET",
		"foo/bar/hello/world",
		nil,
		strings.NewReader(""),
		StringMatches("Test 1."),
	)
}

func Test_E2E_RequestWithPlaceholderVariable(t *testing.T) {
	RunTest(
		t,
		"config_with_wildcards/config.json",
		"GET",
		"user/123",
		nil,
		strings.NewReader(""),
		StringMatches("User ID: 123"),
	)
}
