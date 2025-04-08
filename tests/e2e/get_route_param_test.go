package tests_e2e

import (
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_GetRouteParam(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			"--route say_hi/{name}/{location}",
			CmdExec(`printf "Hi. My name is $({{MOCK_EXECUTABLE}} get-route-param name). I live on $({{MOCK_EXECUTABLE}} get-route-param location)." > $MOCK_RESPONSE_BODY`),
		},
		Get("say_hi/john_doe/earth", nil),
		StringMatches("Hi. My name is john_doe. I live on earth."),
	)
}

func Test_E2E_GetRouteParam_WithUnexistingKeyExitsWith1(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			"--route say_hi/{name}",
			CmdExec(`{{MOCK_EXECUTABLE}} get-route-param this-param-does-not-exist > $MOCK_RESPONSE_BODY`),
		},
		Get("say_hi/john_doe", nil),
		StringMatches(""),
		ExitCodeHeaderMatches("1"),
	)
}

func Test_E2E_GetRouteParam_WithUnexistingKeyExitsWith1_WithRouteWithoutParams(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			"--route say_hi",
			CmdExec(`{{MOCK_EXECUTABLE}} get-route-param this-param-does-not-exist > $MOCK_RESPONSE_BODY`),
		},
		Get("say_hi", nil),
		StringMatches(""),
		ExitCodeHeaderMatches("1"),
	)
}
