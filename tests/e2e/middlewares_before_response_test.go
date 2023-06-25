package tests_e2e

import (
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_Middlewares_BeforeResponse_ModifyBody(t *testing.T) {
	RunTest(
		t,
		"config_with_middlewares/config.json",
		"GET",
		"middleware/before_response/modify_body",
		nil,
		strings.NewReader(""),
		StringMatches("Text: bar."),
	)
}

func Test_Middlewares_BeforeResponse_ModifyBody_WithCmdParams(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route middleware/before_response/modify_body",
			"--response 'Text: foo.'",
			"--middleware-before-response 'sh data/config_with_middlewares/middleware_replace_foo_with_bar.sh'",
		},
		"GET",
		"middleware/before_response/modify_body",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(200),
		StringMatches("Text: bar."),
	)
}

func Test_Middlewares_BeforeResponse_ModifyBody_WithAbsoluteScriptPath_WithCmdParams(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route middleware/before_response/modify_body",
			"--response 'Text: foo.'",
			fmt.Sprintf("--middleware-before-response 'sh %s/data/config_with_middlewares/middleware_replace_foo_with_bar.sh'", pwd),
		},
		"GET",
		"middleware/before_response/modify_body",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(200),
		StringMatches("Text: bar."),
	)
}

func Test_Middlewares_BeforeResponse_ModifyBody_WithFilteredRoute(t *testing.T) {
	RunTest(
		t,
		"config_with_middlewares/config.json",
		"GET",
		"middleware/before_response/modify_body/filtered_routes",
		nil,
		strings.NewReader(""),
		StringMatches("Hello world!Hello world!"),
	)
}
