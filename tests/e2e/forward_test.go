package tests_e2e

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_Forward(t *testing.T) {
	state := NewState()
	killMockBase, _, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--response 'Hello world.'",
		}, " "),
		nil,
		true,
	)
	defer killMockBase()

	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			fmt.Sprintf("--base 'http://localhost:%d'", state.Port),
			"--route foo/bar",
			"--exec '{{MOCK_EXECUTABLE}} forward; printf \" Modified!\" >> $MOCK_RESPONSE_BODY'",
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		StringMatches("Hello world. Modified!"),
	)
}
