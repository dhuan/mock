package tests_e2e

import (
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_Middlewares_ModifyBody(t *testing.T) {
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

func Test_Middlewares_ModifyBody_WithCmdParams(t *testing.T) {
	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route middleware/before_response/modify_body",
			"--response 'Text: foo.'",
			"--middleware 'sh data/config_with_middlewares/middleware_replace_foo_with_bar.sh'",
		},
		"GET",
		"middleware/before_response/modify_body",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(200),
		StringMatches("Text: bar."),
	)
}

func Test_Middlewares_ModifyBody_WithAbsoluteScriptPath_WithCmdParams(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			"--route middleware/before_response/modify_body",
			"--response 'Text: foo.'",
			fmt.Sprintf("--middleware 'sh %s/data/config_with_middlewares/middleware_replace_foo_with_bar.sh'", pwd),
		},
		"GET",
		"middleware/before_response/modify_body",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(200),
		StringMatches("Text: bar."),
	)
}

func Test_Middlewares_ModifyBody_WithFilteredRoute(t *testing.T) {
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

func Test_Middlewares_ModifyHeaders(t *testing.T) {
	RunTest(
		t,
		"config_with_middlewares/config.json",
		"GET",
		"middleware/before_response/modify_headers",
		nil,
		strings.NewReader(""),
		HeadersMatch(map[string]string{
			"Foo":        "bar",
			"Header-One": "Value for header one",
			"Header-Two": "Value for header two",
		}),
	)
}

func Test_Middlewares_RemoveHeaders(t *testing.T) {
	RunTest(
		t,
		"config_with_middlewares/config.json",
		"GET",
		"middleware/before_response/remove_headers",
		nil,
		strings.NewReader(""),
		HeaderKeysNotIncluded([]string{
			"Header-One",
		}),
		HeadersMatch(map[string]string{
			"Header-Two": "Value for header two",
		}),
	)
}

func Test_Middlewares_ModifyStatusCode(t *testing.T) {
	RunTest(
		t,
		"config_with_middlewares/config.json",
		"GET",
		"middleware/before_response/modify_status_code",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(202),
	)
}

func Test_Middlewares_PrintRouteParams(t *testing.T) {
	RunTest(
		t,
		"config_with_middlewares/config.json",
		"GET",
		"middleware/before_response/print_route_params/some_param/another_param",
		nil,
		strings.NewReader(""),
		LineEquals(1, "ROUTE_PARAM_ONE: some_param"),
		LineEquals(2, "ROUTE_PARAM_TWO: another_param"),
	)
}

func Test_Middlewares_PrintEnvironmentVariables(t *testing.T) {
	RunTest(
		t,
		"config_with_middlewares/config.json",
		"GET",
		"middleware/before_response/print_env_vars/some_param/another_param?foo=bar",
		nil,
		strings.NewReader(""),
		LineEquals(1, fmt.Sprintf(`MOCK_HOST=localhost:%s`, GetTestPort())),
		LineRegexMatches(2, `MOCK_REQUEST_BODY=.*`),
		LineEquals(3, `MOCK_REQUEST_ENDPOINT=middleware/before_response/print_env_vars/some_param/another_param`),
		LineRegexMatches(4, `MOCK_REQUEST_HEADERS=.*`),
		LineEquals(5, fmt.Sprintf(`MOCK_REQUEST_HOST=localhost:%s`, GetTestPort())),
		LineEquals(6, `MOCK_REQUEST_METHOD=get`),
		LineEquals(7, `MOCK_REQUEST_QUERYSTRING=foo=bar`),
		LineEquals(8, fmt.Sprintf(`MOCK_REQUEST_URL=http://localhost:%s/middleware/before_response/print_env_vars/some_param/another_param`, GetTestPort())),
	)
}
