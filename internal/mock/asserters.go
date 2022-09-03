package mock

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/dhuan/mock/internal/types"
	. "github.com/dhuan/mock/pkg/mock"
)

func assertHeaderMatch(requestRecord *types.RequestRecord, assert *AssertOptions) (*[]ValidationError, error) {
	validationErrors := make([]ValidationError, 0)
	keyValues := assert.KeyValues
	if keyValues == nil {
		keyValues = make(map[string]interface{})
	}

	if assert.Key != "" && assert.Value != "" {
		keyValues[assert.Key] = fmt.Sprint(assert.Value)
	}

	for i, _ := range keyValues {
		key := i
		value := keyValues[i]

		valueFromRequestRecord, ok := requestRecord.Headers[key]
		if !ok {
			validationErrors = append(validationErrors, ValidationError{
				Code: ValidationErrorCode_HeaderNotIncluded,
				Metadata: map[string]string{
					"missing_header_key": key,
				},
			})

			continue
		}

		if value != strings.Join(valueFromRequestRecord, "") {
			validationErrors = append(validationErrors, ValidationError{
				Code: ValidationErrorCode_HeaderValueMismatch,
				Metadata: map[string]string{
					"header_key":             key,
					"header_value_requested": strings.Join(valueFromRequestRecord, ""),
					"header_value_expected":  value.(string),
				},
			})
		}
	}

	return &validationErrors, nil
}

func assertMethodMatch(requestRecord *types.RequestRecord, assert *AssertOptions) (*[]ValidationError, error) {
	validationErrors := make([]ValidationError, 0)

	if requestRecord.Method != assert.Value {
		validationErrors = append(validationErrors, ValidationError{
			Code: ValidationErrorCode_MethodMismatch,
			Metadata: map[string]string{
				"method_requested": requestRecord.Method,
				"method_expected":  assert.Value,
			},
		})
	}

	return &validationErrors, nil
}

func assertFormMatch(requestRecord *types.RequestRecord, assert *AssertOptions) (*[]ValidationError, error) {
	validationErrors := make([]ValidationError, 0)
	requestBody := string(*requestRecord.Body)

	parsedForm, err := parseForm(requestBody)
	if err != nil {
		panic(err)
	}

	for i, _ := range assert.KeyValues {
		value, ok := parsedForm[i]
		if !ok {
			validationErrors = append(
				validationErrors,
				ValidationError{
					Code: ValidationErrorCode_FormKeyDoesNotExist,
					Metadata: map[string]string{
						"form_key": i,
					},
				},
			)

			continue
		}

		if value != assert.KeyValues[i] {
			validationErrors = append(
				validationErrors,
				ValidationError{
					Code: ValidationErrorCode_FormValueMismatch,
					Metadata: map[string]string{
						"form_key":             i,
						"form_value_requested": value,
						"form_value_expected":  assert.KeyValues[i].(string),
					},
				},
			)
		}
	}

	return &validationErrors, nil
}

func assertQuerystringMatch(requestRecord *types.RequestRecord, assert *AssertOptions) (*[]ValidationError, error) {
	validationErrors := make([]ValidationError, 0)

	if requestRecord.Querystring == "" {
		return &[]ValidationError{ValidationError{
			Code:     ValidationErrorCode_RequestHasNoQuerystring,
			Metadata: map[string]string{},
		}}, nil
	}

	parsedQuery, err := url.ParseQuery(requestRecord.Querystring)
	if err != nil {
		return &validationErrors, err
	}

	expectedKeyValuePairs := make(map[string]string)

	if assert.Key != "" {
		expectedKeyValuePairs[assert.Key] = assert.Value
	}

	for key, value := range assert.KeyValues {
		expectedKeyValuePairs[key] = value.(string)
	}

	for key, _ := range expectedKeyValuePairs {
		_, ok := parsedQuery[key]
		if !ok {
			validationErrors = append(
				validationErrors,
				ValidationError{
					Code: ValidationErrorCode_QuerystringKeyNotSet,
					Metadata: map[string]string{
						"querystring_key": key,
					},
				},
			)

			continue
		}

		if expectedKeyValuePairs[key] != parsedQuery[key][0] {
			validationErrors = append(
				validationErrors,
				ValidationError{
					Code: ValidationErrorCode_QuerystringMismatch,
					Metadata: map[string]string{
						"querystring_key":             key,
						"querystring_value_expected":  expectedKeyValuePairs[key],
						"querystring_value_requested": parsedQuery[key][0],
					},
				},
			)
		}
	}

	return &validationErrors, nil
}

func assertJsonBodyMatch(jsonValidate JsonValidate) func(requestRecord *types.RequestRecord, assert *AssertOptions) (*[]ValidationError, error) {
	return func(requestRecord *types.RequestRecord, assert *AssertOptions) (*[]ValidationError, error) {
		validationErrors := make([]ValidationError, 0)

		var jsonResult map[string]interface{}
		err := json.Unmarshal(*requestRecord.Body, &jsonResult)
		if err != nil {
			return &validationErrors, err
		}

		jsonValidationResult := jsonValidate(jsonResult, assert.Data)
		if !jsonValidationResult {
			assertJson, err := json.Marshal(assert.Data)
			if err != nil {
				panic(err)
			}

			requestRecordReformatted, err := reformatJson(requestRecord.Body)
			if err != nil {
				panic(err)
			}

			validationErrors = append(
				validationErrors,
				ValidationError{Code: ValidationErrorCode_BodyMismatch, Metadata: map[string]string{
					"body_requested": string(requestRecordReformatted),
					"body_expected":  string(assertJson),
				}},
			)
		}

		return &validationErrors, nil
	}
}

func reformatJson(jsonEncoded *[]byte) ([]byte, error) {
	var result map[string]interface{}
	err := json.Unmarshal(*jsonEncoded, &result)
	if err != nil {
		return []byte(""), err
	}

	newJsonEncoded, err := json.Marshal(result)
	if err != nil {
		return []byte(""), err
	}

	return newJsonEncoded, nil
}

func parseForm(requestBody string) (map[string]string, error) {
	formValues := make(map[string]string)

	values, err := url.ParseQuery(requestBody)
	if err != nil {
		return formValues, err
	}

	for i, _ := range values {
		formValues[i] = values[i][0]
	}

	return formValues, nil
}
