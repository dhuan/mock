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

	validationErrors := e2eutils.MockAssert(&mocklib.AssertConfig{
		Route: "foo/bar",
		Assert: &mocklib.AssertOptions{
			Type:  mocklib.AssertType_MethodMatch,
			Value: "post",
		},
	})

	assert.Equal(
		t,
		[]mocklib.ValidationError{
			{Code: mocklib.ValidationErrorCode_NoCall, Metadata: map[string]string{}},
		},
		validationErrors,
	)
}

func Test_E2E_Assertion_BasicAssertion_WithValidationErrors(t *testing.T) {
	killMock := e2eutils.RunMockBg(e2eutils.NewState(), "serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}")
	defer killMock()

	mockConfig := mocklib.Init("localhost:4000")
	e2eutils.Request(mockConfig, "POST", "foo/bar", `{"foo":"bar"}`, map[string]string{})

	validationErrors := e2eutils.MockAssert(&mocklib.AssertConfig{
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
				Code: mocklib.ValidationErrorCode_MethodMismatch,
				Metadata: map[string]string{
					"method_expected":  "put",
					"method_requested": "post",
				},
			},
		},
		validationErrors,
	)
}

func Test_E2E_Assertion_WithNth(t *testing.T) {
	killMock := e2eutils.RunMockBg(e2eutils.NewState(), "serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}")
	defer killMock()

	mockConfig := mocklib.Init("localhost:4000")
	e2eutils.Request(mockConfig, "POST", "foo/bar", `{"foo":"not_bar"}`, e2eutils.ContentTypeJsonHeaders)
	e2eutils.Request(mockConfig, "POST", "foo/bar", `{"foo":"bar"}`, e2eutils.ContentTypeJsonHeaders)

	validationErrors := e2eutils.MockAssert(&mocklib.AssertConfig{
		Route: "foo/bar",
		Nth:   2,
		Assert: &mocklib.AssertOptions{
			Type: mocklib.AssertType_JsonBodyMatch,
			Data: map[string]interface{}{
				"foo": "bar",
			},
		},
	})

	assert.Equal(
		t,
		[]mocklib.ValidationError{},
		validationErrors,
	)
}

func Test_E2E_Assertion_WithNth_Failing(t *testing.T) {
	killMock := e2eutils.RunMockBg(e2eutils.NewState(), "serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}")
	defer killMock()

	mockConfig := mocklib.Init("localhost:4000")
	e2eutils.Request(mockConfig, "POST", "foo/bar", `{"foo":"not_bar"}`, e2eutils.ContentTypeJsonHeaders)
	e2eutils.Request(mockConfig, "POST", "foo/bar", `{"foo":"bar"}`, e2eutils.ContentTypeJsonHeaders)

	validationErrors := e2eutils.MockAssert(&mocklib.AssertConfig{
		Route: "foo/bar",
		Nth:   1,
		Assert: &mocklib.AssertOptions{
			Type: mocklib.AssertType_JsonBodyMatch,
			Data: map[string]interface{}{
				"foo": "bar",
			},
		},
	})

	assert.Equal(
		t,
		[]mocklib.ValidationError{
			mocklib.ValidationError{
				Code: mocklib.ValidationErrorCode_BodyMismatch,
				Metadata: map[string]string{
					"body_expected":  `{"foo":"bar"}`,
					"body_requested": `{"foo":"not_bar"}`,
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

	validationErrors := e2eutils.MockAssert(&mocklib.AssertConfig{
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

	validationErrors := e2eutils.MockAssert(&mocklib.AssertConfig{
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
				Code: mocklib.ValidationErrorCode_HeaderNotIncluded,
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

	validationErrors := e2eutils.MockAssert(&mocklib.AssertConfig{
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
