package tests_e2e

import (
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
		"",
		HeadersMatch(map[string]string{
			"Some-Header-Key":    "Some header value",
			"Another-Header-Key": "Another header value",
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
		"",
		HeadersMatch(map[string]string{
			"Base-Header-One":    "A base header",
			"Base-Header-Two":    "Another base header",
			"Some-Header-Key":    "Some header value",
			"Another-Header-Key": "Another header value",
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
		"",
		HeadersMatch(map[string]string{
			"Base-Header-One":                     "A base header",
			"Base-Header-Two":                     "Another base header",
			"Header-For-Conditional-Response-One": "Some header value",
			"Header-For-Conditional-Response-Two": "Another header value",
		}),
		HeaderKeysNotIncluded([]string{
			"Some-Header-Key",
			"Another-Header-Key",
		}),
	)
}
