package tests_e2e

import (
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_OptionsMethod_WithoutCorsFlag_405(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			"--route foo/bar",
			"--response 'Hello, world.'",
		},
		Options("foo/bar", nil),
		StatusCodeMatches(405),
	)
}

func Test_E2E_OptionsMethod_WithCorsFlag_CorsHeaders(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			"--cors",
			"--route foo/bar",
			"--response 'Hello, world.'",
		},
		Options("foo/bar", nil),
		HeadersMatch(map[string][]string{
			"Access-Control-Allow-Origin":      {"*"},
			"Access-Control-Allow-Credentials": {"true"},
			"Access-Control-Allow-Headers":     {"*"},
			"Access-Control-Allow-Methods":     {"POST, GET, OPTIONS, PUT, DELETE"},
		}),
	)
}

func Test_E2E_OptionsMethod_WithCorsFlag_WithMiddleware(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			"--cors",
			"--middleware '/home/work/work/mock/bin/mock set-header foo bar'",
			"--route foo/bar",
			"--response 'Hello, world.'",
		},
		Options("foo/bar", nil),
		HeadersMatch(map[string][]string{
			"Access-Control-Allow-Origin":      {"*"},
			"Access-Control-Allow-Credentials": {"true"},
			"Access-Control-Allow-Headers":     {"*"},
			"Access-Control-Allow-Methods":     {"POST, GET, OPTIONS, PUT, DELETE"},
			"Foo":                              {"bar"},
		}),
	)
}
