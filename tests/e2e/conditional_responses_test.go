package tests_e2e

import (
	"testing"

	mocklib "github.com/dhuan/mock/pkg/mock"
	e2eutils "github.com/dhuan/mock/tests/e2e/utils"
	"github.com/stretchr/testify/assert"
)

func Test_E2E_ConditionalResponses_ReceivingDefaultResponse(t *testing.T) {
	killMock := e2eutils.RunMockBg(e2eutils.NewState(), "serve -c {{TEST_DATA_PATH}}/config_with_conditional_response/config.json -p {{TEST_E2E_PORT}}")
	defer killMock()

	mockConfig := mocklib.Init("localhost:4000")

	responseBody := e2eutils.Request(mockConfig, "POST", "foo/bar", "", map[string]string{})

	assert.Equal(
		t,
		"This is the default response.",
		string(responseBody),
	)
}

func Test_E2E_ConditionalResponses_ReceivingConditionalResponse(t *testing.T) {
	killMock := e2eutils.RunMockBg(e2eutils.NewState(), "serve -c {{TEST_DATA_PATH}}/config_with_conditional_response/config.json -p {{TEST_E2E_PORT}}")
	defer killMock()

	mockConfig := mocklib.Init("localhost:4000")

	responseBody := e2eutils.Request(mockConfig, "POST", "foo/bar?key1=value1&key2=value2", "", map[string]string{})

	assert.Equal(
		t,
		"Hello world!",
		string(responseBody),
	)
}

func Test_E2E_ConditionalResponses_ReceivingConditionalResponse_WithAndChaining(t *testing.T) {
	killMock := e2eutils.RunMockBg(e2eutils.NewState(), "serve -c {{TEST_DATA_PATH}}/config_with_conditional_response/config.json -p {{TEST_E2E_PORT}}")
	defer killMock()

	mockConfig := mocklib.Init("localhost:4000")

	responseBody := e2eutils.Request(mockConfig, "POST", "foo/bar?key1=value1&key2=value2&key4=value4", "", map[string]string{})

	assert.Equal(
		t,
		"Hello world! (Condition with AND chaining)",
		string(responseBody),
	)
}

func Test_E2E_ConditionalResponses_ReceivingConditionalResponse_WithOrChaining(t *testing.T) {
	killMock := e2eutils.RunMockBg(e2eutils.NewState(), "serve -c {{TEST_DATA_PATH}}/config_with_conditional_response/config.json -p {{TEST_E2E_PORT}}")
	defer killMock()

	mockConfig := mocklib.Init("localhost:4000")

	responseBody := e2eutils.Request(mockConfig, "POST", "foo/bar?key1=value1&key6=value6", "", map[string]string{})

	assert.Equal(
		t,
		"Hello world! (Condition with OR chaining)",
		string(responseBody),
	)
}
