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

func TestValidate_NoCalls(t *testing.T) {
	assertConfig := mock.AssertConfig{
		Route: "foobar",
		Headers: map[string][]string{
			"foo": []string{"bar"},
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

func TestValidate_HeaderNotIncluded(t *testing.T) {
	reset()
	addToMockedRequestRecords(
		"foobar",
		"get",
		[][]string{[]string{"some_header_key", "some_header_value"}},
		[]byte(``),
	)

	assertConfig := mock.AssertConfig{
		Route: "foobar",
		Headers: map[string][]string{
			"foo": []string{"bar"},
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

func TestValidate_HeaderMismatch(t *testing.T) {
	reset()
	addToMockedRequestRecords(
		"foobar",
		"get",
		[][]string{[]string{"foo", "not_bar"}},
		[]byte(``),
	)

	assertConfig := mock.AssertConfig{
		Route: "foobar",
		Headers: map[string][]string{
			"foo": []string{"bar"},
		},
	}

	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidate, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{
			mock.ValidationError{
				Code: mock.Validation_error_code_header_value_mismatch,
				Metadata: map[string]string{
					"header_key":             "foo",
					"header_value_requested": "not_bar",
					"header_value_expected":  "bar",
				},
			},
		},
		validationErrors,
	)
}

func TestValidate_BodyJson_ValueMatches(t *testing.T) {
	reset()
	addToMockedRequestRecords(
		"foobar",
		"post",
		[][]string{
			[]string{"content-type", "application/json"},
		},
		[]byte(`{"foo":"not_bar"}`),
	)

	assertConfig := mock.AssertConfig{
		Route: "foobar",
		BodyJson: map[string]interface{}{
			"foo": "bar",
		},
	}

	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidate, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{
			mock.ValidationError{
				Code: mock.Validation_error_code_body_mismatch,
				Metadata: map[string]string{
					"body_requested": `{"foo":"not_bar"}`,
					"body_expected":  `{"foo":"bar"}`,
				},
			},
		},
		validationErrors,
	)
}

func TestValidate_BodyJson_RequestWithBodyButNoBodyAssertion(t *testing.T) {
	reset()
	addToMockedRequestRecords(
		"foobar",
		"post",
		[][]string{
			[]string{"content-type", "application/json"},
		},
		[]byte(`{"foo":"bar"}`),
	)

	assertConfig := mock.AssertConfig{
		Route: "foobar",
	}

	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidate, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{},
		validationErrors,
	)
}

func TestValidate_BodyJson_RequestWithoutBodyButWithBodyAssertion(t *testing.T) {
	reset()
	addToMockedRequestRecords(
		"foobar",
		"post",
		[][]string{
			[]string{"content-type", "application/json"},
		},
		[]byte(""),
	)

	assertConfig := mock.AssertConfig{
		Route: "foobar",
		BodyJson: map[string]interface{}{
			"foo": "bar",
		},
	}

	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidate, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{
			mock.ValidationError{
				Code:     mock.Validation_error_code_request_has_no_body_content,
				Metadata: map[string]string{},
			},
		},
		validationErrors,
	)
}

func TestValidate_MethodMismatch(t *testing.T) {
	reset()
	addToMockedRequestRecords(
		"foobar",
		"post",
		[][]string{},
		[]byte(""),
	)

	assertConfig := mock.AssertConfig{
		Route:  "foobar",
		Method: "put",
	}

	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidate, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{
			mock.ValidationError{
				Code: mock.Validation_error_code_method_mismatch,
				Metadata: map[string]string{
					"method_requested": "post",
					"method_expected":  "put",
				},
			},
		},
		validationErrors,
	)
}

func TestValidate_MethodMatch(t *testing.T) {
	reset()
	addToMockedRequestRecords(
		"foobar",
		"post",
		[][]string{},
		[]byte(""),
	)

	assertConfig := mock.AssertConfig{
		Route:  "foobar",
		Method: "post",
	}

	validationErrors, _ := mock.Validate(mockMockFs{}, mockJsonValidate, &assertConfig)

	assert.Equal(
		t,
		&[]mock.ValidationError{},
		validationErrors,
	)
}
