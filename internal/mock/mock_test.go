package mock_test

import (
	"net/http"
	"testing"

	"github.com/dhuan/mock/internal/mock"
	"github.com/dhuan/mock/internal/types"
	"github.com/stretchr/testify/assert"
)

var mock_request_records []*types.RequestRecord

func reset() {
	mock_request_records = make([]*types.RequestRecord, 0)
}

func addToMockedRequestRecords(route string, headers [][]string) {
	httpHeaders := make(http.Header)

	for _, headerValues := range headers {
		headerKey := headerValues[0]
		httpHeaders[headerKey] = headerValues[1:]
	}

	mock_request_records = append(
		mock_request_records,
		&types.RequestRecord{
			Route:   route,
			Headers: httpHeaders,
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

func TestValidate_NoCalls(t *testing.T) {
	assertConfig := mock.AssertConfig{
		Route: "foobar",
		Headers: map[string][]string{
			"foo": []string{"bar"},
		},
	}

	_, validationErrors, _ := mock.Validate(mockMockFs{}, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{
			mock.ValidationError{Code: mock.Validation_error_code_no_call, Metadata: []string{}},
		},
		validationErrors,
	)
}

func TestValidate_HeaderNotIncluded(t *testing.T) {
	reset()
	addToMockedRequestRecords(
		"foobar",
		[][]string{[]string{"some_header_key", "some_header_value"}},
	)

	assertConfig := mock.AssertConfig{
		Route: "foobar",
		Headers: map[string][]string{
			"foo": []string{"bar"},
		},
	}

	_, validationErrors, _ := mock.Validate(mockMockFs{}, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{
			mock.ValidationError{Code: mock.Validation_error_code_header_not_included, Metadata: []string{"foo"}},
		},
		validationErrors,
	)
}

func TestValidate_HeaderMismatch(t *testing.T) {
	reset()
	addToMockedRequestRecords(
		"foobar",
		[][]string{[]string{"foo", "not_bar"}},
	)

	assertConfig := mock.AssertConfig{
		Route: "foobar",
		Headers: map[string][]string{
			"foo": []string{"bar"},
		},
	}

	_, validationErrors, _ := mock.Validate(mockMockFs{}, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{
			mock.ValidationError{Code: mock.Validation_error_code_header_value_mismatch, Metadata: []string{"foo", "bar", "not_bar"}},
		},
		validationErrors,
	)
}
