package tests_e2e

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_BaseApi_NormalRequest(t *testing.T) {
	state := NewState()
	killMockBase, _, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--status-code 205",
			"--response 'Hello world! This is the base API.'",
		}, " "),
		nil,
		true,
	)
	defer killMockBase()

	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			fmt.Sprintf("--base 'localhost:%d'", state.Port),
			"--route hello/world",
			"--response 'Hello world!'",
		},
		"GET",
		"hello/world",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(200),
		StringMatches("Hello world!"),
	)
}

func Test_E2E_BaseApi_RequestForwardedToBaseApi(t *testing.T) {
	state := NewState()
	killMockBase, _, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--status-code 205",
			"--response 'Hello world! This is the base API.'",
		}, " "),
		nil,
		true,
	)
	defer killMockBase()

	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			fmt.Sprintf("--base 'localhost:%d'", state.Port),
			"--route hello/world",
			"--response 'Hello world!'",
		},
		"GET",
		"foo/bar",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(205),
		StringMatches("Hello world! This is the base API."),
	)
}

func Test_E2E_BaseApi_RequestForwardedToBaseApi_WithJsonConfig(t *testing.T) {
	state := NewState()
	killMockBase, _, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--status-code 205",
			"--response 'Hello world! This is the base API.'",
		}, " "),
		nil,
		true,
	)
	defer killMockBase()

	RunTestWithJsonConfig(
		t,
		fmt.Sprintf(`{
    "base": "localhost:%d",
    "endpoints": [
        {
            "route": "hello/world",
            "response": "Hello world!"
        }
    ]
}`, state.Port),
		[]string{},
		"GET",
		"foo/bar",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(205),
		StringMatches("Hello world! This is the base API."),
	)
}

func Test_E2E_BaseApi_RequestForwardedToBaseApi_WithQuerystring(t *testing.T) {
	state := NewState()
	killMockBase, _, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--response 'Querystring: ${MOCK_REQUEST_QUERYSTRING}'",
		}, " "),
		nil,
		true,
	)
	defer killMockBase()

	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			fmt.Sprintf("--base 'localhost:%d'", state.Port),
			"--route hello/world",
			"--response 'Hello world!'",
		},
		"GET",
		"foo/bar?param_one=value_one&param_two=value_two",
		nil,
		strings.NewReader(""),
		StringMatches("Querystring: param_one=value_one&param_two=value_two"),
	)
}

func Test_E2E_BaseApi_RequestForwardedToBaseApi_ResponseHeadersAreForwaded(t *testing.T) {
	state := NewState()
	killMockBase, _, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--header 'Header-One: value one'",
			"--header 'Header-Two: value two'",
			"--response 'Hello world!'",
		}, " "),
		nil,
		true,
	)
	defer killMockBase()

	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			fmt.Sprintf("--base 'localhost:%d'", state.Port),
			"--route hello/world",
			"--response 'Hello world!'",
		},
		"GET",
		"foo/bar",
		nil,
		strings.NewReader(""),
		HeadersMatch(map[string][]string{
			"Header-One": []string{"value one"},
			"Header-Two": []string{"value two"},
		}),
	)
}

func Test_E2E_BaseApi_RequestForwardedToBaseApi_RequestBodyIsForwarded(t *testing.T) {
	state := NewState()
	killMockBase, _, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--method POST",
			"--response 'Request body: ${MOCK_REQUEST_BODY}'",
		}, " "),
		nil,
		true,
	)
	defer killMockBase()

	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			fmt.Sprintf("--base 'localhost:%d'", state.Port),
			"--route hello/world",
			"--response 'Hello world!'",
		},
		"POST",
		"foo/bar",
		nil,
		strings.NewReader("This is a request payload."),
		StringMatches("Request body: This is a request payload."),
	)
}

func Test_E2E_BaseApi_RequestForwardedToBaseApi_RequestHeadersAreForwarded(t *testing.T) {
	state := NewState()
	killMockBase, _, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--response 'Request header foo: ${MOCK_REQUEST_HEADER_FOO}'",
		}, " "),
		nil,
		true,
	)
	defer killMockBase()

	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			fmt.Sprintf("--base 'localhost:%d'", state.Port),
			"--route hello/world",
			"--response 'Hello world!'",
		},
		"GET",
		"foo/bar",
		map[string]string{
			"Foo": "bar",
		},
		strings.NewReader(""),
		StringMatches("Request header foo: bar"),
	)
}

func Test_E2E_BaseApi_Middleware_ModifyingResponse(t *testing.T) {
	state := NewState()
	killMockBase, _, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--header 'Some-Header-From-Base-Api: some value'",
			"--response 'Base Api Response.'",
		}, " "),
		nil,
		true,
	)
	defer killMockBase()

	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			fmt.Sprintf("--base 'localhost:%d'", state.Port),
			"--route hello/world",
			"--response 'Hello world!'",
			"--middleware 'printf \"Foo: bar\" >> $MOCK_RESPONSE_HEADERS'",
			"--middleware 'printf \" Modified!\" >> $MOCK_RESPONSE_BODY'",
		},
		"GET",
		"foo/bar",
		nil,
		strings.NewReader(""),
		HeadersMatch(map[string][]string{
			"Some-Header-From-Base-Api": []string{"some value"},
			"Foo":                       []string{"bar"},
		}),
		StringMatches("Base Api Response. Modified!"),
	)
}

func Test_E2E_BaseApi_Middleware_ModifyingResponse_UsingMOCK_BASE_API_RESPONSE(t *testing.T) {
	state := NewState()
	killMockBase, _, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--header 'Some-Header-From-Base-Api: some value'",
			"--response 'Base Api Response.'",
		}, " "),
		nil,
		true,
	)
	defer killMockBase()

	middlewareScript := []string{
		"if [ \"$MOCK_BASE_API_RESPONSE\" = true ];",
		"then printf \" This response was proxied from a Base API.\" >> $MOCK_RESPONSE_BODY;",
		"else printf \" This response was NOT proxied from a Base API.\" >> $MOCK_RESPONSE_BODY;",
		"fi",
	}

	commandOptions := []string{
		fmt.Sprintf("--base 'localhost:%d'", state.Port),
		"--route hello/world",
		"--response 'Hello world!'",
		fmt.Sprintf("--middleware '%s'", strings.Join(middlewareScript, " ")),
	}

	RunTestWithNoConfigAndWithArgs(
		t,
		commandOptions,
		"GET",
		"foo/bar",
		nil,
		strings.NewReader(""),
		StringMatches("Base Api Response. This response was proxied from a Base API."),
	)

	RunTestWithNoConfigAndWithArgs(
		t,
		commandOptions,
		"GET",
		"hello/world",
		nil,
		strings.NewReader(""),
		StringMatches("Hello world! This response was NOT proxied from a Base API."),
	)
}

func Test_E2E_BaseApi_RequestForwardedToBaseApi_CmdFlagOverwritesConfig(t *testing.T) {
	state := NewState()
	killMockBase, _, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--status-code 205",
			"--response 'Hello world! This is the base API 1.'",
		}, " "),
		nil,
		true,
	)
	defer killMockBase()

	state2 := NewState()
	killMockBase2, _, _, _ := RunMockBg(
		state2,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--status-code 205",
			"--response 'Hello world! This is the base API 2.'",
		}, " "),
		nil,
		true,
	)
	defer killMockBase2()

	type testCase struct {
		portBaseFlag     int
		portBaseConfig   int
		expectedResponse string
	}

	testCases := []testCase{
		{portBaseFlag: state2.Port, portBaseConfig: state.Port, expectedResponse: "Hello world! This is the base API 2."},
		{portBaseFlag: state.Port, portBaseConfig: state2.Port, expectedResponse: "Hello world! This is the base API 1."},
	}

	for _, testCase := range testCases {
		RunTestWithJsonConfig(
			t,
			fmt.Sprintf(`{
        "base": "localhost:%d",
        "endpoints": [
            {
                "route": "hello/world",
                "response": "Hello world!"
            }
        ]
    }`, testCase.portBaseConfig),
			[]string{fmt.Sprintf("--base localhost:%d", testCase.portBaseFlag)},
			"GET",
			"foo/bar",
			nil,
			strings.NewReader(""),
			StatusCodeMatches(205),
			StringMatches(testCase.expectedResponse),
		)
	}
}

func Test_E2E_BaseApi_WithoutAnyEndpoints(t *testing.T) {
	state := NewState()
	killMockBase, _, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--response 'Base Api Response.'",
		}, " "),
		nil,
		true,
	)
	defer killMockBase()

	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			fmt.Sprintf("--base 'localhost:%d'", state.Port),
		},
		"GET",
		"foo/bar",
		nil,
		strings.NewReader(""),
		StringMatches("Base Api Response."),
	)
}

func Test_E2E_BaseApi_WithoutAnyEndpoints_WithJsonConfig(t *testing.T) {
	state := NewState()
	killMockBase, _, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--response 'Base Api Response.'",
		}, " "),
		nil,
		true,
	)
	defer killMockBase()

	RunTestWithJsonConfig(
		t,
		fmt.Sprintf(`{
    "base": "localhost:%d"
}`, state.Port),
		[]string{},
		"GET",
		"foo/bar",
		nil,
		strings.NewReader(""),
		StringMatches("Base Api Response."),
	)
}

func Test_E2E_BaseApi_CorsFlagOverwritesCorsHeadersFromBaseApi(t *testing.T) {
	state := NewState()
	killMockBase, _, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--response 'Hello world! This is the base API.'",
			"--header 'Access-Control-Allow-Origin: some_origin'",
			"--header 'Access-Control-Allow-Credentials: false'",
			"--header 'Access-Control-Allow-Headers: some-header'",
			"--header 'Access-Control-Allow-Methods: GET'",
		}, " "),
		nil,
		true,
	)
	defer killMockBase()

	RunTestWithNoConfigAndWithArgs(
		t,
		[]string{
			fmt.Sprintf("--base 'localhost:%d'", state.Port),
			"--cors",
		},
		"GET",
		"foo/bar",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(200),
		HeadersMatch(map[string][]string{
			"Access-Control-Allow-Origin":      []string{"*"},
			"Access-Control-Allow-Credentials": []string{"true"},
			"Access-Control-Allow-Headers":     []string{"*"},
			"Access-Control-Allow-Methods":     []string{"POST, GET, OPTIONS, PUT, DELETE"},
		}),
	)
}
