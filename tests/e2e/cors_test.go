package tests_e2e

import (
	"net/http"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
	e2eutils "github.com/dhuan/mock/tests/e2e/utils"
	"github.com/stretchr/testify/assert"
)

func Test_E2E_Cors_HeadersAreSet(t *testing.T) {
	killMock, serverOutput, mockConfig, _ := e2eutils.RunMockBg(
		e2eutils.NewState(),
		"serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}} --cors",
		nil,
		true,
		nil,
	)
	defer killMock()

	response := e2eutils.Request(mockConfig, "POST", "foo/bar", strings.NewReader(`{"foo":"bar"}`), e2eutils.ContentTypeJsonHeaders, serverOutput)

	e2eutils.AssertMapHasValues(t, response.Headers, http.Header{
		"Access-Control-Allow-Credentials": []string{"true"},
		"Access-Control-Allow-Headers":     []string{"*"},
		"Access-Control-Allow-Methods":     []string{"POST, GET, OPTIONS, PUT, DELETE"},
		"Access-Control-Allow-Origin":      []string{"*"},
	})
}

func Test_E2E_Cors_HeadersAreNotSet(t *testing.T) {
	killMock, serverOutput, mockConfig, _ := e2eutils.RunMockBg(
		e2eutils.NewState(),
		"serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}",
		nil,
		true,
		nil,
	)
	defer killMock()

	response := e2eutils.Request(mockConfig, "POST", "foo/bar", strings.NewReader(`{"foo":"bar"}`), e2eutils.ContentTypeJsonHeaders, serverOutput)

	headerKeys := e2eutils.GetKeys(response.Headers)

	assert.Equal(t, e2eutils.IndexOf(headerKeys, "Access-Control-Allow-Credentials"), -1)
	assert.Equal(t, e2eutils.IndexOf(headerKeys, "Access-Control-Allow-Headers"), -1)
	assert.Equal(t, e2eutils.IndexOf(headerKeys, "Access-Control-Allow-Methods"), -1)
	assert.Equal(t, e2eutils.IndexOf(headerKeys, "Access-Control-Allow-Origin"), -1)
}

func Test_E2E_Cors_WithUnexistingRoute(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			"--route test",
			"--response 'Hello, world.'",
			"--cors",
		},
		Get("this/route/does/not/exist", nil),
		StatusCodeMatches(405),
		HeadersMatch(http.Header{
			"Access-Control-Allow-Credentials": []string{"true"},
			"Access-Control-Allow-Headers":     []string{"*"},
			"Access-Control-Allow-Methods":     []string{"POST, GET, OPTIONS, PUT, DELETE"},
			"Access-Control-Allow-Origin":      []string{"*"},
		}),
	)
}
