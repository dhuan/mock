package tests_e2e

import (
	"testing"

	e2eutils "github.com/dhuan/mock/tests/e2e/utils"
	"github.com/stretchr/testify/assert"
)

func Test_E2E_WithInvalidConfig(t *testing.T) {
	out, _ := e2eutils.RunMock(e2eutils.NewState(), "serve -c {{TEST_DATA_PATH}}/config_invalid/config.json -p {{TEST_E2E_PORT}}")

	assert.Equal(
		t,
		`mock can't be started. The following errors were found in your configuration:

1: Endpoint #1 (invalid_method foo/bar):
The given method, "invalid_method" , is invalid. The available HTTP Methods you can use are POST, GET, PUT, PATCH, and DELETE.

2: Endpoint #2 (get another/endpoint?foo=bar):
Routes cannot have querystrings. Read about "response_if" in the documentation to learn how to set Conditional Responses based on querystrings.`,
		string(out),
	)
}
