package mock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dhuan/mock/internal/types"
	. "github.com/dhuan/mock/pkg/mock"
)

type Kv struct {
	Key   string
	Value string
}

type JsonValidate func(jsonA map[string]interface{}, jsonB map[string]interface{}) bool

func ParseAssertRequest(req *http.Request) (*AssertConfig, error) {
	var assertConfig AssertConfig
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&assertConfig)

	return &assertConfig, err
}

func Validate(
	mockFs types.MockFs,
	jsonValidate JsonValidate,
	assertConfig *AssertConfig,
) (*[]ValidationError, error) {
	validationErrors := make([]ValidationError, 0)
	requestRecords, err := getRequestRecordMatchingRoute(mockFs, assertConfig.Route)
	if err != nil {
		return &validationErrors, err
	}
	if len(requestRecords) == 0 {
		validationErrors = append(
			validationErrors,
			ValidationError{
				Code:     Validation_error_code_no_call,
				Metadata: map[string]string{},
			},
		)

		return &validationErrors, nil
	}

	nth := assertConfig.Nth
	if nth == 0 {
		nth = 1
	}
	requestRecord := requestRecords[nth-1]

	return validate(requestRecord, assertConfig.Assert, jsonValidate)
}

func validate(requestRecord *types.RequestRecord, assert *AssertOptions, jsonValidate JsonValidate) (*[]ValidationError, error) {
	hasAnd := assert.And != nil
	hasOr := assert.Or != nil
	validationErrors := make([]ValidationError, 0)
	assertFunc := resolveAssertTypeFunc(assert.Type, jsonValidate)
	validationErrorsCurrent, err := assertFunc(requestRecord, assert)
	success := len(*validationErrorsCurrent) == 0
	if err != nil {
		return &validationErrors, err
	}

	if !success {
		validationErrors = append(validationErrors, *validationErrorsCurrent...)
	}

	if success && !hasAnd {
		return &validationErrors, nil
	}

	if success && hasAnd {
		furtherValidationErrors, err := validate(requestRecord, assert.And, jsonValidate)
		if err != nil {
			return &validationErrors, err
		}

		validationErrors = append(*furtherValidationErrors, *validationErrorsCurrent...)
	}

	if !success && hasOr {
		furtherValidationErrors, err := validate(requestRecord, assert.Or, jsonValidate)
		if err != nil {
			return &validationErrors, err
		}

		if len(*furtherValidationErrors) == 0 {
			return furtherValidationErrors, nil
		}

		validationErrors = append(*furtherValidationErrors, *validationErrorsCurrent...)
	}

	return &validationErrors, nil
}

func resolveAssertTypeFunc(
	assertType AssertType,
	jsonValidate JsonValidate,
) func(requestRecord *types.RequestRecord, assert *AssertOptions) (*[]ValidationError, error) {
	if assertType == AssertType_HeaderMatch {
		return assertHeaderMatch
	}

	if assertType == AssertType_MethodMatch {
		return assertMethodMatch
	}

	if assertType == AssertType_JsonBodyMatch {
		return assertJsonBodyMatch(jsonValidate)
	}

	if assertType == AssertType_FormMatch {
		return assertFormMatch
	}

	panic(fmt.Sprintf("Failed to resolve assert type: %d", assertType))
}

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
				Code: Validation_error_code_header_not_included,
				Metadata: map[string]string{
					"missing_header_key": key,
				},
			})

			continue
		}

		if value != strings.Join(valueFromRequestRecord, "") {
			validationErrors = append(validationErrors, ValidationError{
				Code: Validation_error_code_header_value_mismatch,
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
			Code: Validation_error_code_method_mismatch,
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
					Code: Validation_error_code_form_key_does_not_exist,
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
					Code: Validation_error_code_form_value_mismatch,
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
				ValidationError{Code: Validation_error_code_body_mismatch, Metadata: map[string]string{
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

func getRequestRecordMatchingRoute(mockFs types.MockFs, route string) ([]*types.RequestRecord, error) {
	requestRecords, err := mockFs.GetRecordsMatchingRoute(route)
	if err != nil {
		return requestRecords, err
	}

	if len(requestRecords) == 0 {
		return requestRecords, err
	}

	return requestRecords, nil
}
