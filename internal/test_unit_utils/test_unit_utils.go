package test_unit_utils

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/dhuan/mock/internal/mock"
	"github.com/dhuan/mock/internal/types"
	. "github.com/dhuan/mock/pkg/mock"
	"github.com/stretchr/testify/assert"
)

type TestOperationFunc func(t *testing.T, state *unitTestState)

type unitTestState struct {
	validationErrors []ValidationError
	requestRecords   []types.RequestRecord
}

func RunUnitTest(t *testing.T, funcs ...TestOperationFunc) {
	state := &unitTestState{
		make([]ValidationError, 0),
		make([]types.RequestRecord, 0),
	}

	for i := range funcs {
		funcs[i](t, state)
	}
}

func AddGetRequestRecord(route string) TestOperationFunc {
	return func(t *testing.T, state *unitTestState) {
		AddToMockedRequestRecords(
			state,
			route,
			"get",
			[][]string{},
			[]byte(``),
		)
	}
}

func AddGetRequestRecordWithHeaders(route string, headers [][]string) TestOperationFunc {
	return func(t *testing.T, state *unitTestState) {
		AddToMockedRequestRecords(
			state,
			route,
			"get",
			headers,
			[]byte(``),
		)
	}
}

func AddPostRequestRecordWithPayload(route string, payload string) TestOperationFunc {
	return func(t *testing.T, state *unitTestState) {
		AddToMockedRequestRecords(
			state,
			route,
			"post",
			[][]string{},
			[]byte(payload),
		)
	}
}

func AddToMockedRequestRecords(state *unitTestState, fullRoute, method string, headers [][]string, body []byte) {
	httpHeaders := make(http.Header)

	for _, headerValues := range headers {
		headerKey := headerValues[0]
		httpHeaders[headerKey] = headerValues[1:]
	}

	route, querystring := parseRoute(fullRoute)

	state.requestRecords = append(
		state.requestRecords,
		types.RequestRecord{
			Route:       route,
			Method:      method,
			Headers:     httpHeaders,
			Body:        &body,
			Querystring: querystring,
		},
	)
}

func Validate(route string, assertOptions *Condition) TestOperationFunc {
	return func(t *testing.T, state *unitTestState) {
		assertConfig := AssertConfig{
			Route:  route,
			Assert: assertOptions,
		}

		newValidationErrors, err := mock.Validate(mockMockFs{state}, &assertConfig)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}

		state.validationErrors = append(state.validationErrors, *newValidationErrors...)
	}
}

func ValidateNth(nth int, route string, assertOptions *Condition) TestOperationFunc {
	return func(t *testing.T, state *unitTestState) {
		assertConfig := AssertConfig{
			Nth:    nth,
			Route:  route,
			Assert: assertOptions,
		}

		newValidationErrors, err := mock.Validate(mockMockFs{state}, &assertConfig)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}

		state.validationErrors = append(state.validationErrors, *newValidationErrors...)
	}
}

func RemoveValidationErrors(t *testing.T, state *unitTestState) {
	state.validationErrors = make([]ValidationError, 0)
}

func AssertOptionsWithKeyValue(assertType ConditionType, key, value string) *Condition {
	return &Condition{
		Type:  assertType,
		Key:   key,
		Value: ConditionValue(value),
	}
}

func AssertOptionsWithKeyValues(assertType ConditionType, keyValues map[string]interface{}) *Condition {
	return &Condition{
		Type:      assertType,
		KeyValues: keyValues,
	}
}

func AssertOptionsWithData(assertType ConditionType, data map[string]interface{}) *Condition {
	return &Condition{
		Type:      assertType,
		KeyValues: data,
	}
}

func AssertOptionsWithValue(assertType ConditionType, value string) *Condition {
	return &Condition{
		Type:  assertType,
		Value: ConditionValue(value),
	}
}

func ExpectOneValidationError(errorCode ValidationErrorCode, errorMetadata map[string]string) TestOperationFunc {
	return func(t *testing.T, state *unitTestState) {
		assert.Equal(
			t,
			[]ValidationError{
				{
					Code:     errorCode,
					Metadata: errorMetadata,
				},
			},
			state.validationErrors,
		)
	}
}

func ExpectValidationErrorNth(index int, errorCode ValidationErrorCode, errorMetadata map[string]string) TestOperationFunc {
	return func(t *testing.T, state *unitTestState) {
		assert.Equal(
			t,
			ValidationError{
				Code:     errorCode,
				Metadata: errorMetadata,
			},
			state.validationErrors[index],
		)
	}
}

func ExpectZeroValidationErrors(t *testing.T, state *unitTestState) {
	assert.Equal(t, 0, len(state.validationErrors))
}

func ExpectValidationErrorsCount(count int) TestOperationFunc {
	return func(t *testing.T, state *unitTestState) {
		assert.Equal(
			t,
			count,
			len(state.validationErrors),
		)
	}
}

func Reset(state *unitTestState) {
	*state = unitTestState{}
	state.requestRecords = make([]types.RequestRecord, 0)
	state.validationErrors = make([]ValidationError, 0)
}

func parseRoute(fullRoute string) (string, string) {
	split := strings.Split(fullRoute, "?")
	if len(split) < 2 {
		return fullRoute, ""
	}

	return split[0], split[1]
}

type mockMockFs struct {
	State *unitTestState
}

func (this mockMockFs) StoreRequestRecord(r *http.Request, requestBody []byte, endpointConfig *types.EndpointConfig) error {
	return nil
}

func (this mockMockFs) GetRecordsMatchingRoute(route string) ([]types.RequestRecord, error) {
	return this.State.requestRecords, nil
}

func (this mockMockFs) RemoveAllRequestRecords() error {
	return nil
}
