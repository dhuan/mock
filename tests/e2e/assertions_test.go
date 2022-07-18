package tests_e2e

import (
	"testing"

	mocklib "github.com/dhuan/mock/pkg/mock"
	e2eutils "github.com/dhuan/mock/tests/e2e/utils"
	"github.com/stretchr/testify/assert"
)

func Test_E2E_Assertion_NoCalls(t *testing.T) {
	killMock := e2eutils.RunMockBg(e2eutils.NewState(), "serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}")
	defer killMock()

	mockConfig := mocklib.Init("localhost:4000")
	validationErrors := mocklib.Assert(mockConfig, &mocklib.AssertConfig{
		Route: "foo/bar",
		Assert: &mocklib.AssertOptions{
			Type:  mocklib.AssertType_MethodMatch,
			Value: "post",
		},
	})

	assert.Equal(
		t,
		[]mocklib.ValidationError{
			{Code: "no_call", Metadata: map[string]string{}},
		},
		validationErrors,
	)
}

func Test_E2E_Assertion_BasicAssertion_WithValidationErrors(t *testing.T) {
	killMock := e2eutils.RunMockBg(e2eutils.NewState(), "serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}")
	defer killMock()

	mockConfig := mocklib.Init("localhost:4000")

	e2eutils.Request(mockConfig, "POST", "foo/bar", `{"foo":"bar"}`, map[string]string{})

	validationErrors := mocklib.Assert(mockConfig, &mocklib.AssertConfig{
		Route: "foo/bar",
		Assert: &mocklib.AssertOptions{
			Type:  mocklib.AssertType_MethodMatch,
			Value: "put",
		},
	})

	assert.Equal(
		t,
		[]mocklib.ValidationError{
			{
				Code: "method_mismatch",
				Metadata: map[string]string{
					"method_expected":  "put",
					"method_requested": "post",
				},
			},
		},
		validationErrors,
	)
}

func Test_E2E_Assertion_BasicAssertion_WithoutValidationErrors(t *testing.T) {
	killMock := e2eutils.RunMockBg(e2eutils.NewState(), "serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}")
	defer killMock()

	mockConfig := mocklib.Init("localhost:4000")

	e2eutils.Request(mockConfig, "POST", "foo/bar", `{"foo":"bar"}`, map[string]string{})

	validationErrors := mocklib.Assert(mockConfig, &mocklib.AssertConfig{
		Route: "foo/bar",
		Assert: &mocklib.AssertOptions{
			Type:  mocklib.AssertType_MethodMatch,
			Value: "post",
		},
	})

	assert.Equal(
		t,
		[]mocklib.ValidationError{},
		validationErrors,
	)
}

func Test_E2E_Assertion_Chaining_WithValidationErrors(t *testing.T) {
	killMock := e2eutils.RunMockBg(e2eutils.NewState(), "serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}")
	defer killMock()

	mockConfig := mocklib.Init("localhost:4000")

	e2eutils.Request(mockConfig, "POST", "foo/bar", `{"foo":"bar"}`, map[string]string{})

	validationErrors := mocklib.Assert(mockConfig, &mocklib.AssertConfig{
		Route: "foo/bar",
		Assert: &mocklib.AssertOptions{
			Type:  mocklib.AssertType_MethodMatch,
			Value: "post",
			And: &mocklib.AssertOptions{
				Type: mocklib.AssertType_HeaderMatch,
				KeyValues: map[string]interface{}{
					"some_header_key": "some_header_value",
				},
			},
		},
	})

	assert.Equal(
		t,
		[]mocklib.ValidationError{
			{
				Code: "header_not_included",
				Metadata: map[string]string{
					"missing_header_key": "some_header_key",
				},
			},
		},
		validationErrors,
	)
}

func Test_E2E_Assertion_Chaining_WithoutValidationErrors(t *testing.T) {
	killMock := e2eutils.RunMockBg(e2eutils.NewState(), "serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}")
	defer killMock()

	mockConfig := mocklib.Init("localhost:4000")

	e2eutils.Request(mockConfig, "POST", "foo/bar", `{"foo":"bar"}`, map[string]string{
		"some_header_key": "some_header_value",
	})

	validationErrors := mocklib.Assert(mockConfig, &mocklib.AssertConfig{
		Route: "foo/bar",
		Assert: &mocklib.AssertOptions{
			Type:  mocklib.AssertType_MethodMatch,
			Value: "post",
			And: &mocklib.AssertOptions{
				Type: mocklib.AssertType_HeaderMatch,
				KeyValues: map[string]interface{}{
					"some_header_key": "some_header_value",
				},
			},
		},
	})

	assert.Equal(
		t,
		[]mocklib.ValidationError{},
		validationErrors,
	)
}
