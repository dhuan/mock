package tests_e2e

import (
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_Middlewares_NotFound(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			`--middleware 'test "${MOCK_REQUEST_NOT_FOUND}" = "true" && {{MOCK_EXECUTABLE}} set-status 201 && (echo "NOT FOUND!" | mock write)'`,
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
