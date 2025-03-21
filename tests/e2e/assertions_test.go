package tests_e2e

import (
	"strings"
	"testing"

	mocklib "github.com/dhuan/mock/pkg/mock"
	e2eutils "github.com/dhuan/mock/tests/e2e/utils"
	"github.com/stretchr/testify/assert"
)

func Test_E2E_Assertion_NoCalls(t *testing.T) {
	state := e2eutils.NewState()
	killMock, serverOutput, _, _ := e2eutils.RunMockBg(
		state,
		"serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}",
		nil,
		true,
		nil,
	)

	defer killMock()

	validationErrors := e2eutils.MockAssert(&mocklib.AssertOptions{
		Route: "foo/bar",
		Condition: &mocklib.Condition{
			Type:  mocklib.ConditionType_MethodMatch,
			Value: "post",
		},
	}, serverOutput, state)

	assert.Equal(
		t,
		[]mocklib.ValidationError{
			{Code: mocklib.ValidationErrorCode_NoCall, Metadata: map[string]string{}},
		},
		validationErrors,
	)
}

func Test_E2E_Assertion_BasicAssertion_WithValidationErrors(t *testing.T) {
	state := e2eutils.NewState()
	killMock, serverOutput, mockConfig, _ := e2eutils.RunMockBg(
		state,
		"serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}",
		nil,
		true,
		nil,
	)
	defer killMock()

	e2eutils.Request(mockConfig, "POST", "foo/bar", strings.NewReader(`{"foo":"bar"}`), map[string]string{}, serverOutput)

	validationErrors := e2eutils.MockAssert(&mocklib.AssertOptions{
		Route: "foo/bar",
		Condition: &mocklib.Condition{
			Type:  mocklib.ConditionType_MethodMatch,
			Value: "put",
		},
	}, serverOutput, state)

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
	state := e2eutils.NewState()
	killMock, serverOutput, mockConfig, _ := e2eutils.RunMockBg(
		state,
		"serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}",
		nil,
		true,
		nil,
	)
	defer killMock()

	e2eutils.Request(mockConfig, "POST", "foo/bar", strings.NewReader(`{"foo":"not_bar"}`), e2eutils.ContentTypeJsonHeaders, serverOutput)
	e2eutils.Request(mockConfig, "POST", "foo/bar", strings.NewReader(`{"foo":"bar"}`), e2eutils.ContentTypeJsonHeaders, serverOutput)

	validationErrors := e2eutils.MockAssert(&mocklib.AssertOptions{
		Route: "foo/bar",
		Nth:   2,
		Condition: &mocklib.Condition{
			Type: mocklib.ConditionType_JsonBodyMatch,
			KeyValues: map[string]interface{}{
				"foo": "bar",
			},
		},
	}, serverOutput, state)

	assert.Equal(
		t,
		[]mocklib.ValidationError{},
		validationErrors,
	)
}

func Test_E2E_Assertion_WithNth_Failing(t *testing.T) {
	state := e2eutils.NewState()

	killMock, serverOutput, mockConfig, _ := e2eutils.RunMockBg(
		state,
		"serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}",
		nil,
		true,
		nil,
	)
	defer killMock()

	e2eutils.Request(mockConfig, "POST", "foo/bar", strings.NewReader(`{"foo":"not_bar"}`), e2eutils.ContentTypeJsonHeaders, serverOutput)
	e2eutils.Request(mockConfig, "POST", "foo/bar", strings.NewReader(`{"foo":"bar"}`), e2eutils.ContentTypeJsonHeaders, serverOutput)

	validationErrors := e2eutils.MockAssert(&mocklib.AssertOptions{
		Route: "foo/bar",
		Nth:   1,
		Condition: &mocklib.Condition{
			Type: mocklib.ConditionType_JsonBodyMatch,
			KeyValues: map[string]interface{}{
				"foo": "bar",
			},
		},
	}, serverOutput, state)

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

func Test_E2E_Assertion_WithNth_OutOfRange(t *testing.T) {
	state := e2eutils.NewState()
	killMock, serverOutput, mockConfig, _ := e2eutils.RunMockBg(
		state,
		"serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}",
		nil,
		true,
		nil,
	)
	defer killMock()

	e2eutils.Request(mockConfig, "POST", "foo/bar", strings.NewReader(`{"foo":"bar"}`), e2eutils.ContentTypeJsonHeaders, serverOutput)

	validationErrors := e2eutils.MockAssert(&mocklib.AssertOptions{
		Route: "foo/bar",
		Nth:   2,
		Condition: &mocklib.Condition{
			Type: mocklib.ConditionType_JsonBodyMatch,
			KeyValues: map[string]interface{}{
				"foo": "bar",
			},
		},
	}, serverOutput, state)

	assert.Equal(
		t,
		[]mocklib.ValidationError{
			mocklib.ValidationError{
				Code:     mocklib.ValidationErrorCode_NthOutOfRange,
				Metadata: map[string]string{},
			},
		},
		validationErrors,
	)
}

func Test_E2E_Assertion_BasicAssertion_WithoutValidationErrors(t *testing.T) {
	state := e2eutils.NewState()
	killMock, serverOutput, mockConfig, _ := e2eutils.RunMockBg(
		state,
		"serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}",
		nil,
		true,
		nil,
	)
	defer killMock()

	e2eutils.Request(mockConfig, "POST", "foo/bar", strings.NewReader(`{"foo":"bar"}`), map[string]string{}, serverOutput)

	validationErrors := e2eutils.MockAssert(&mocklib.AssertOptions{
		Route: "foo/bar",
		Condition: &mocklib.Condition{
			Type:  mocklib.ConditionType_MethodMatch,
			Value: "post",
		},
	}, serverOutput, state)

	assert.Equal(
		t,
		[]mocklib.ValidationError{},
		validationErrors,
	)
}

func Test_E2E_Assertion_Chaining_WithValidationErrors(t *testing.T) {
	state := e2eutils.NewState()
	killMock, serverOutput, mockConfig, _ := e2eutils.RunMockBg(
		state,
		"serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}",
		nil,
		true,
		nil,
	)
	defer killMock()

	e2eutils.Request(mockConfig, "POST", "foo/bar", strings.NewReader(`{"foo":"bar"}`), map[string]string{}, serverOutput)

	validationErrors := e2eutils.MockAssert(&mocklib.AssertOptions{
		Route: "foo/bar",
		Condition: &mocklib.Condition{
			Type:  mocklib.ConditionType_MethodMatch,
			Value: "post",
			And: &mocklib.Condition{
				Type: mocklib.ConditionType_HeaderMatch,
				KeyValues: map[string]interface{}{
					"some_header_key": "some_header_value",
				},
			},
		},
	}, serverOutput, state)

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
	state := e2eutils.NewState()
	killMock, serverOutput, mockConfig, _ := e2eutils.RunMockBg(
		state,
		"serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}",
		nil,
		true,
		nil,
	)
	defer killMock()

	e2eutils.Request(mockConfig, "POST", "foo/bar", strings.NewReader(`{"foo":"bar"}`), map[string]string{
		"some_header_key": "some_header_value",
	}, serverOutput)

	validationErrors := e2eutils.MockAssert(&mocklib.AssertOptions{
		Route: "foo/bar",
		Condition: &mocklib.Condition{
			Type:  mocklib.ConditionType_MethodMatch,
			Value: "post",
			And: &mocklib.Condition{
				Type: mocklib.ConditionType_HeaderMatch,
				KeyValues: map[string]interface{}{
					"some_header_key": "some_header_value",
				},
			},
		},
	}, serverOutput, state)

	assert.Equal(
		t,
		[]mocklib.ValidationError{},
		validationErrors,
	)
}

func Test_E2E_Assertion_MethodMatchingIsCaseInsensitive(t *testing.T) {
	state := e2eutils.NewState()
	killMock, serverOutput, mockConfig, _ := e2eutils.RunMockBg(
		state,
		"serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}",
		nil,
		true,
		nil,
	)
	defer killMock()

	assertMethods := []string{"post", "POST"}

	for _, assertMethod := range assertMethods {
		e2eutils.RequestApiReset(mockConfig)
		e2eutils.Request(mockConfig, "POST", "foo/bar", strings.NewReader(`{"foo":"bar"}`), map[string]string{}, serverOutput)

		validationErrors := e2eutils.MockAssert(&mocklib.AssertOptions{
			Route: "foo/bar",
			Condition: &mocklib.Condition{
				Type:  mocklib.ConditionType_MethodMatch,
				Value: mocklib.ConditionValue(assertMethod),
			},
		}, serverOutput, state)

		assert.Equal(t, 0, len(validationErrors))
	}
}
