package mock

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
)

type AssertType int

const (
	AssertType_None AssertType = iota
	AssertType_HeaderMatch
	AssertType_MethodMatch
	AssertType_JsonBodyMatch
	AssertType_FormMatch
)

func (this *AssertType) UnmarshalJSON(data []byte) error {
	assertTypeText := utils.Unquote(string(data))

	if assertTypeText == "header_match" {
		*this = AssertType_HeaderMatch

		return nil
	}

	if assertTypeText == "method_match" {
		*this = AssertType_MethodMatch

		return nil
	}

	if assertTypeText == "json_body_match" {
		*this = AssertType_JsonBodyMatch

		return nil
	}

	if assertTypeText == "form_match" {
		*this = AssertType_FormMatch

		return nil
	}

	return errors.New(fmt.Sprintf("Failed to parse Assert Type: %s", assertTypeText))
}

var (
	Validation_error_code_header_value_mismatch       = "header_value_mismatch"
	Validation_error_code_no_call                     = "no_call"
	Validation_error_code_header_not_included         = "header_not_included"
	Validation_error_code_body_mismatch               = "body_mismatch"
	Validation_error_code_request_has_no_body_content = "request_has_no_body_content"
	Validation_error_code_method_mismatch             = "method_mismatch"
	Validation_error_code_form_key_does_not_exist     = "form_key_does_not_exist"
	Validation_error_code_form_value_mismatch         = "form_value_mismatch"
)

type AssertHeader map[string][]string

type AssertConfig struct {
	Route    string                 `json:"route"`
	Nth      int                    `json:"nth"`
	Method   string                 `json:"method"`
	Headers  AssertHeader           `json:"headers"`
	BodyJson map[string]interface{} `json:"body_json"`
	Assert   *Assert                `json:"assert"`
}

type Assert struct {
	Type      AssertType             `json:"type"`
	Data      map[string]interface{} `json:"data"`
	KeyValues map[string]interface{} `json:"key_values"`
	Key       string                 `json:"key"`
	Value     string                 `json:"value"`
	And       *Assert                `json:"and"`
	Or        *Assert                `json:"or"`
}

type Kv struct {
	Key   string
	Value string
}

type ValidationError struct {
	Code     string            `json:"code"`
	Metadata map[string]string `json:"metadata"`
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
		validationErrors = append(validationErrors, ValidationError{Validation_error_code_no_call, map[string]string{}})

		return &validationErrors, nil
	}

	nth := assertConfig.Nth
	if nth == 0 {
		nth = 1
	}
	requestRecord := requestRecords[nth-1]

	return validate(requestRecord, assertConfig.Assert, jsonValidate)
}

func validate(requestRecord *types.RequestRecord, assert *Assert, jsonValidate JsonValidate) (*[]ValidationError, error) {
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
) func(requestRecord *types.RequestRecord, assert *Assert) (*[]ValidationError, error) {
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

func assertHeaderMatch(requestRecord *types.RequestRecord, assert *Assert) (*[]ValidationError, error) {
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

func assertMethodMatch(requestRecord *types.RequestRecord, assert *Assert) (*[]ValidationError, error) {
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

func assertFormMatch(requestRecord *types.RequestRecord, assert *Assert) (*[]ValidationError, error) {
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

func assertJsonBodyMatch(jsonValidate JsonValidate) func(requestRecord *types.RequestRecord, assert *Assert) (*[]ValidationError, error) {
	return func(requestRecord *types.RequestRecord, assert *Assert) (*[]ValidationError, error) {
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

func handleBodyJsonAssertion(
	requestRecord *types.RequestRecord,
	jsonValidate JsonValidate,
	assertConfig *AssertConfig,
) (*[]ValidationError, error) {
	validationErrors := make([]ValidationError, 0)

	if string(*requestRecord.Body) == "" {
		validationErrors = append(
			validationErrors,
			ValidationError{Code: Validation_error_code_request_has_no_body_content, Metadata: map[string]string{}},
		)

		return &validationErrors, nil
	}

	var jsonA map[string]interface{}
	err := json.Unmarshal(*requestRecord.Body, &jsonA)
	if err != nil {
		return &validationErrors, err
	}

	bodyValidationResult := jsonValidate(jsonA, assertConfig.BodyJson)
	bodyRequest, err := json.Marshal(jsonA)
	if err != nil {
		return &validationErrors, err
	}
	bodyAssert, err := json.Marshal(assertConfig.BodyJson)
	if err != nil {
		return &validationErrors, err
	}

	if !bodyValidationResult {
		validationErrors = append(
			validationErrors,
			ValidationError{Validation_error_code_body_mismatch, map[string]string{
				"body_requested": string(bodyRequest),
				"body_expected":  string(bodyAssert),
			}})
	}

	return &validationErrors, nil
}

func validateHeadersMatch(requestRecord *types.RequestRecord, assertConfig *AssertConfig) *[]ValidationError {
	validationErrors := make([]ValidationError, 0)

	for headerKey, header := range assertConfig.Headers {
		headerB, ok := requestRecord.Headers[headerKey]
		if !ok {
			validationErrors = append(validationErrors, ValidationError{
				Code: Validation_error_code_header_not_included,
				Metadata: map[string]string{
					"missing_header_key": headerKey,
				},
			})

			continue
		}

		if !utils.ListsEqual[string](header, headerB) {
			validationErrors = append(validationErrors, ValidationError{
				Code: Validation_error_code_header_value_mismatch,
				Metadata: map[string]string{
					"header_key":             headerKey,
					"header_value_expected":  strings.Join(header, ","),
					"header_value_requested": strings.Join(headerB, ","),
				},
			})
		}
	}

	return &validationErrors
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
