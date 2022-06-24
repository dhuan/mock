package mock_test

import (
	"net/http"
	"testing"

	"github.com/dhuan/mock/internal/mock"
	"github.com/dhuan/mock/internal/types"
	"github.com/stretchr/testify/assert"
)

var mock_request_records []*types.RequestRecord

func mockJsonValidate(jsonA map[string]interface{}, jsonB map[string]interface{}) bool {
	return false
}

func reset() {
	mock_request_records = make([]*types.RequestRecord, 0)
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

func (this mockMockFs) StoreRequestRecord(r *http.Request, endpointConfig *types.EndpointConfig) error {
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
			KeyValues: []mock.Kv{
				mock.Kv{Key: "foo", Value: "bar"},
			},
		},
	}

	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidate, &assertConfig)

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
			KeyValues: []mock.Kv{
				mock.Kv{
					Key:   "foo",
					Value: "bar",
				},
			},
		},
	}

	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidate, &assertConfig)

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
			KeyValues: []mock.Kv{
				mock.Kv{
					Key:   "foo",
					Value: "bar",
				},
				mock.Kv{
					Key:   "foo2",
					Value: "bar2",
				},
			},
		},
	}

	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidate, &assertConfig)

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
			KeyValues: []mock.Kv{
				mock.Kv{
					Key:   "some_header_key",
					Value: "some_header_value",
				},
			},
			And: &mock.Assert{
				Type:  mock.AssertType_MethodMatch,
				Value: "post",
			},
		},
	}

	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidate, &assertConfig)

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
			KeyValues: []mock.Kv{
				mock.Kv{
					Key:   "some_header_key",
					Value: "some_header_value",
				},
			},
			And: &mock.Assert{
				Type:  mock.AssertType_MethodMatch,
				Value: "get",
			},
		},
	}

	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidate, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{},
		validationErrors,
	)
}
