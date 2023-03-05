package tests_e2e

import (
	"testing"

	e2eutils "github.com/dhuan/mock/tests/e2e/utils"
	"github.com/stretchr/testify/assert"
)

func Test_E2E_Cors_HeadersAreSet(t *testing.T) {
	killMock, _, mockConfig := e2eutils.RunMockBg(
		e2eutils.NewState(),
		"serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}} --cors",
		nil,
	)
	defer killMock()

	response := e2eutils.Request(mockConfig, "POST", "foo/bar", `{"foo":"bar"}`, e2eutils.ContentTypeJsonHeaders)

	e2eutils.AssertMapHasValues(t, response.Headers, map[string]string{
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Headers":     "*",
		"Access-Control-Allow-Methods":     "POST, GET, OPTIONS, PUT, DELETE",
		"Access-Control-Allow-Origin":      "*",
	})
}

func Test_E2E_Cors_HeadersAreNotSet(t *testing.T) {
	killMock, _, mockConfig := e2eutils.RunMockBg(
		e2eutils.NewState(),
		"serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}",
		nil,
	)
	defer killMock()

	response := e2eutils.Request(mockConfig, "POST", "foo/bar", `{"foo":"bar"}`, e2eutils.ContentTypeJsonHeaders)

	headerKeys := e2eutils.GetKeys(response.Headers)

	assert.Equal(t, e2eutils.IndexOf(headerKeys, "Access-Control-Allow-Credentials"), -1)
	assert.Equal(t, e2eutils.IndexOf(headerKeys, "Access-Control-Allow-Headers"), -1)
	assert.Equal(t, e2eutils.IndexOf(headerKeys, "Access-Control-Allow-Methods"), -1)
	assert.Equal(t, e2eutils.IndexOf(headerKeys, "Access-Control-Allow-Origin"), -1)
}
