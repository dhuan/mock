package tests_e2e

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

var middleware_script string = strings.Join([]string{
	`test "${MOCK_REQUEST_NOT_FOUND}" = "true"`,
	`&& {{MOCK_EXECUTABLE}} set-status 201`,
	`&& (echo "NOT FOUND!" | {{MOCK_EXECUTABLE}} write)`,
}, " ")

func Test_Middlewares_NotFound_ModifyResponse(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			fmt.Sprintf(`--middleware '%s'`, middleware_script),
			"--route foo/bar",
			`--response "Hello, world."`,
		},
		Get("foo/bar", nil),
		StringMatches("Hello, world."),
		StatusCodeMatches(200),
		Get("no_route", nil),
		StringMatches("NOT FOUND!"),
		StatusCodeMatches(201),
	)
}

func Test_Middlewares_NotFound_ModifyResponse_WithRouteFiltering(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			fmt.Sprintf(`--middleware '%s'`, middleware_script),
			"--route-match foo/bar/3",
			"--route foo/bar",
			`--response "Hello, world."`,
		},
		Get("foo/bar", nil),
		StringMatches("Hello, world."),
		StatusCodeMatches(200),
		Get("foo/bar/2", nil),
		StringMatches(""),
		StatusCodeMatches(405),
		Get("foo/bar/3", nil),
		StringMatches("NOT FOUND!"),
		StatusCodeMatches(201),
	)
}
