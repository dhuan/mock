package tests_e2e

import (
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_ResponseWithHeaders(t *testing.T) {
	RunTest(
		t,
		"config_with_headers/config.json",
		"GET",
		"with/headers",
		nil,
		strings.NewReader(""),
		HeadersMatch(map[string][]string{
			"Some-Header-Key":    []string{"Some header value"},
			"Another-Header-Key": []string{"Another header value"},
		}),
	)
}

func Test_E2E_ResponseWithHeaders_AndBaseHeaders(t *testing.T) {
	RunTest(
		t,
		"config_with_headers/config.json",
		"GET",
		"with/headers/and/base/headers",
		nil,
		strings.NewReader(""),
		HeadersMatch(map[string][]string{
			"Base-Header-One":    []string{"A base header"},
			"Base-Header-Two":    []string{"Another base header"},
			"Some-Header-Key":    []string{"Some header value"},
			"Another-Header-Key": []string{"Another header value"},
		}),
	)
}

func Test_E2E_ResponseWithHeaders_WithConditionalResponse(t *testing.T) {
	RunTest(
		t,
		"config_with_headers/config.json",
		"GET",
		"with/conditional/responses/and/base/headers?some_key=some_value",
		nil,
		strings.NewReader(""),
		HeadersMatch(map[string][]string{
			"Base-Header-One":                     []string{"A base header"},
			"Base-Header-Two":                     []string{"Another base header"},
			"Header-For-Conditional-Response-One": []string{"Some header value"},
			"Header-For-Conditional-Response-Two": []string{"Another header value"},
		}),
		HeaderKeysNotIncluded([]string{
			"Some-Header-Key",
			"Another-Header-Key",
		}),
	)
}
