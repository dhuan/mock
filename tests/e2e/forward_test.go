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
			"--status-code 206",
			"--header 'Header-One: value one'",
		}, " "),
		nil,
		true,
		nil,
	)
	defer killMockBase()

	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			fmt.Sprintf("--base 'http://localhost:%d'", state.Port),
			"--route foo/bar",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} forward`,
				`printf " Modified!" >> $MOCK_RESPONSE_BODY`,
				`STATUS_CODE=$(cat $MOCK_RESPONSE_STATUS_CODE)`,
				`echo $((STATUS_CODE+1)) > $MOCK_RESPONSE_STATUS_CODE`,
				`printf "\nHeader-Two: value two" >> $MOCK_RESPONSE_HEADERS`,
			}, ";")),
		},
		"GET",
		"foo/bar",
		nil,
		nil,
		StringMatches("Hello world. Modified!"),
		StatusCodeMatches(207),
		HeadersMatch(map[string][]string{
			"Header-One": {"value one"},
			"Header-Two": {"value two"},
		}),
	)
}
