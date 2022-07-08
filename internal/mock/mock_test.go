package mock_test

import (
	"net/http"
	"testing"

	"github.com/dhuan/mock/internal/mock"
	"github.com/dhuan/mock/internal/types"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

var mock_request_records []*types.RequestRecord

type mockJsonValidate struct {
	testifymock.Mock
}

var mockJsonValidateInstance = mockJsonValidate{}

func (this *mockJsonValidate) JsonValidate(jsonA map[string]interface{}, jsonB map[string]interface{}) bool {
	args := this.Called(jsonA, jsonB)

	return args.Get(0).(bool)
}

func reset() {
	mock_request_records = make([]*types.RequestRecord, 0)
	mockJsonValidateInstance = mockJsonValidate{}
}

func addToMockedRequestRecords(route string, method string, headers [][]string, body []byte) {
	httpHeaders := make(http.Header)

	for _, headerValues := range headers {
		headerKey := headerValues[0]
		httpHeaders[headerKey] = headerValues[1:]
	}

	mock_request_records = append(
		mock_request_records,
		&types.RequestRecord{
			Route:   route,
			Method:  method,
			Headers: httpHeaders,
			Body:    &body,
		},
	)
}

type mockMockFs struct {
	State *types.State
}

func (this mockMockFs) StoreRequestRecord(r *http.Request, requestBody []byte, endpointConfig *types.EndpointConfig) error {
	return nil
}

func (this mockMockFs) GetRecordsMatchingRoute(route string) ([]*types.RequestRecord, error) {
	return mock_request_records, nil
}

func Test_Validate_NoCalls(t *testing.T) {
	assertConfig := mock.AssertConfig{
		Route: "foobar",
		Assert: &mock.Assert{
			Type: mock.AssertType_HeaderMatch,
			KeyValues: map[string]interface{}{
				"foo": "bar",
			},
		},
	}

	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidateInstance.JsonValidate, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{
			mock.ValidationError{
				Code:     mock.Validation_error_code_no_call,
				Metadata: map[string]string{},
			},
		},
		validationErrors,
	)
}

func Test_Validate_HeaderNotIncluded(t *testing.T) {
	reset()
	addToMockedRequestRecords(
		"foobar",
		"get",
		[][]string{[]string{"some_header_key", "some_header_value"}},
		[]byte(``),
	)

	assertConfig := mock.AssertConfig{
		Route: "foobar",
		Assert: &mock.Assert{
			Type: mock.AssertType_HeaderMatch,
			KeyValues: map[string]interface{}{
				"foo": "bar",
			},
		},
	}

	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidateInstance.JsonValidate, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{
			mock.ValidationError{
				Code: mock.Validation_error_code_header_not_included,
				Metadata: map[string]string{
					"missing_header_key": "foo",
				},
			},
		},
		validationErrors,
	)
}

func Test_Validate_HeaderNotIncludedMany(t *testing.T) {
	reset()
	addToMockedRequestRecords(
		"foobar",
		"get",
		[][]string{[]string{"some_header_key", "some_header_value"}},
		[]byte(``),
	)

	assertConfig := mock.AssertConfig{
		Route: "foobar",
		Assert: &mock.Assert{
			Type: mock.AssertType_HeaderMatch,
			KeyValues: map[string]interface{}{
				"foo":  "bar",
				"foo2": "bar2",
			},
		},
	}

	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidateInstance.JsonValidate, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{
			mock.ValidationError{
				Code: mock.Validation_error_code_header_not_included,
				Metadata: map[string]string{
					"missing_header_key": "foo",
				},
			},
			mock.ValidationError{
				Code: mock.Validation_error_code_header_not_included,
				Metadata: map[string]string{
					"missing_header_key": "foo2",
				},
			},
		},
		validationErrors,
	)
}

func Test_Validate_HeaderMismatch(t *testing.T) {
	reset()
	addToMockedRequestRecords(
		"foobar",
		"get",
		[][]string{[]string{"some_header_key", "some_header_value"}},
		[]byte(``),
	)

	assertConfig := mock.AssertConfig{
		Route: "foobar",
		Assert: &mock.Assert{
			Type: mock.AssertType_HeaderMatch,
			KeyValues: map[string]interface{}{
				"some_header_key": "a_different_header_value",
			},
		},
	}

	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidateInstance.JsonValidate, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{
			mock.ValidationError{
				Code: mock.Validation_error_code_header_value_mismatch,
				Metadata: map[string]string{
					"header_key":             "some_header_key",
					"header_value_requested": "some_header_value",
					"header_value_expected":  "a_different_header_value",
				},
			},
		},
		validationErrors,
	)
}

func Test_Validate_WithAndChainingAssertingMethodAndHeader_Fail(t *testing.T) {
	reset()
	addToMockedRequestRecords(
		"foobar",
		"get",
		[][]string{[]string{"some_header_key", "some_header_value"}},
		[]byte(``),
	)

	assertConfig := mock.AssertConfig{
		Route: "foobar",
		Assert: &mock.Assert{
			Type: mock.AssertType_HeaderMatch,
			KeyValues: map[string]interface{}{
				"some_header_key": "some_header_value",
			},
			And: &mock.Assert{
				Type:  mock.AssertType_MethodMatch,
				Value: "post",
			},
		},
	}

	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidateInstance.JsonValidate, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{
			mock.ValidationError{
				Code: mock.Validation_error_code_method_mismatch,
				Metadata: map[string]string{
					"method_requested": "get",
					"method_expected":  "post",
				},
			},
		},
		validationErrors,
	)
}

func Test_Validate_WithAndChainingAssertingMethodAndHeader(t *testing.T) {
	reset()
	addToMockedRequestRecords(
		"foobar",
		"get",
		[][]string{[]string{"some_header_key", "some_header_value"}},
		[]byte(``),
	)

	assertConfig := mock.AssertConfig{
		Route: "foobar",
		Assert: &mock.Assert{
			Type: mock.AssertType_HeaderMatch,
			KeyValues: map[string]interface{}{
				"some_header_key": "some_header_value",
			},
			And: &mock.Assert{
				Type:  mock.AssertType_MethodMatch,
				Value: "get",
			},
		},
	}

	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidateInstance.JsonValidate, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{},
		validationErrors,
	)
}

func Test_Validate_JsonBodyAssertion_Match(t *testing.T) {
	reset()
	addToMockedRequestRecords(
		"foobar",
		"get",
		[][]string{},
		[]byte(`{"foo":"bar", "some_key": "some_value"}`),
	)

	assertConfig := mock.AssertConfig{
		Route: "foobar",
		Assert: &mock.Assert{
			Type: mock.AssertType_JsonBodyMatch,
			Data: map[string]interface{}{
				"foo":      "bar",
				"some_key": "some_value",
			},
		},
	}

	mockJsonValidateInstance.On(
		"JsonValidate",
		map[string]interface{}{"foo": "bar", "some_key": "some_value"},
		map[string]interface{}{"foo": "bar", "some_key": "some_value"},
	).Return(true)
	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidateInstance.JsonValidate, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{},
		validationErrors,
	)
}

func Test_Validate_JsonBodyAssertion_Mismatch(t *testing.T) {
	reset()
	addToMockedRequestRecords(
		"foobar",
		"get",
		[][]string{},
		[]byte(`{"foo":"bar","some_key":"some_value"}`),
	)

	assertConfig := mock.AssertConfig{
		Route: "foobar",
		Assert: &mock.Assert{
			Type: mock.AssertType_JsonBodyMatch,
			Data: map[string]interface{}{
				"foo":         "bar",
				"some_key":    "some_value",
				"another_key": "another_value",
			},
		},
	}

	mockJsonValidateInstance.On(
		"JsonValidate",
		map[string]interface{}{"foo": "bar", "some_key": "some_value"},
		map[string]interface{}{"foo": "bar", "some_key": "some_value", "another_key": "another_value"},
	).Return(false)
	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidateInstance.JsonValidate, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{
			mock.ValidationError{
				Code: mock.Validation_error_code_body_mismatch,
				Metadata: map[string]string{
					"body_requested": `{"foo":"bar","some_key":"some_value"}`,
					"body_expected":  `{"another_key":"another_value","foo":"bar","some_key":"some_value"}`,
				},
			},
		},
		validationErrors,
	)
}

func Test_Validate_Nth(t *testing.T) {
	reset()
	addToMockedRequestRecords(
		"foobar",
		"get",
		[][]string{},
		[]byte(``),
	)
	addToMockedRequestRecords(
		"foobar",
		"post",
		[][]string{},
		[]byte(`{"foo":"bar","some_key":"some_value"}`),
	)

	assertConfig := mock.AssertConfig{
		Route: "foobar",
		Assert: &mock.Assert{
			Type:  mock.AssertType_MethodMatch,
			Value: "get",
		},
	}
	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidateInstance.JsonValidate, &assertConfig)

	assert.Equal(t, 0, len(*validationErrors))

	assertConfig = mock.AssertConfig{
		Route: "foobar",
		Nth:   2,
		Assert: &mock.Assert{
			Type:  mock.AssertType_MethodMatch,
			Value: "get",
		},
	}
	validationErrors, _ = mock.Validate(mockMockFs{}, mockJsonValidateInstance.JsonValidate, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{
			mock.ValidationError{
				Code: mock.Validation_error_code_method_mismatch,
				Metadata: map[string]string{
					"method_requested": "post",
					"method_expected":  "get",
				},
			},
		},
		validationErrors,
	)

	assertConfig = mock.AssertConfig{
		Route: "foobar",
		Nth:   2,
		Assert: &mock.Assert{
			Type:  mock.AssertType_MethodMatch,
			Value: "post",
		},
	}
	validationErrors, _ = mock.Validate(mockMockFs{}, mockJsonValidateInstance.JsonValidate, &assertConfig)

	assert.Equal(t, 0, len(*validationErrors))

	assertConfig = mock.AssertConfig{
		Route: "foobar",
		Nth:   1,
		Assert: &mock.Assert{
			Type:  mock.AssertType_MethodMatch,
			Value: "get",
		},
	}
	validationErrors, _ = mock.Validate(mockMockFs{}, mockJsonValidateInstance.JsonValidate, &assertConfig)

	assert.Equal(t, 0, len(*validationErrors))
}
