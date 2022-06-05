package mock

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
)

var (
	Validation_error_code_header_value_mismatch       = "header_value_mismatch"
	Validation_error_code_no_call                     = "no_call"
	Validation_error_code_header_not_included         = "header_not_included"
	Validation_error_code_body_mismatch               = "body_mismatch"
	Validation_error_code_request_has_no_body_content = "request_has_no_body_content"
)

type AssertHeader map[string][]string

type AssertConfig struct {
	Route    string                 `json:"route"`
	Headers  AssertHeader           `json:"headers"`
	BodyJson map[string]interface{} `json:"body_json"`
}

type ValidationError struct {
	Code     string   `json:"code"`
	Metadata []string `json:"metadata"`
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
) (bool, *[]ValidationError, error) {
	validationErrors := make([]ValidationError, 0)

	requestRecords, err := getRequestRecordMatchingRoute(mockFs, assertConfig.Route)
	if err != nil {
		return false, &validationErrors, err
	}
	if len(requestRecords) == 0 {
		validationErrors = append(validationErrors, ValidationError{Validation_error_code_no_call, []string{}})

		return false, &validationErrors, nil
	}

	requestRecord := requestRecords[0]

	if len(assertConfig.Headers) > 0 {
		headersMatchValidationErrors := validateHeadersMatch(requestRecord, assertConfig)

		if len(*headersMatchValidationErrors) > 0 {
			validationErrors = append(validationErrors, *headersMatchValidationErrors...)
		}
	}

	hasBodyJsonAssertion := len(assertConfig.BodyJson) > 0
	if hasBodyJsonAssertion {
		bodyJsonAssertionValidationErrors, err := handleBodyJsonAssertion(requestRecord, jsonValidate, assertConfig)
		if err != nil {
			return false, &validationErrors, err
		}

		if len(*bodyJsonAssertionValidationErrors) > 0 {
			validationErrors = append(validationErrors, *bodyJsonAssertionValidationErrors...)
		}
	}

	return len(validationErrors) == 0, &validationErrors, nil
}

func handleBodyJsonAssertion(
	requestRecord *types.RequestRecord,
	jsonValidate JsonValidate,
	assertConfig *AssertConfig,
) (*[]ValidationError, error) {
	validationErrors := make([]ValidationError, 0)

	if string(*requestRecord.Body) == "" {
		validationErrors = append(validationErrors, ValidationError{Code: Validation_error_code_request_has_no_body_content, Metadata: []string{}})

		return &validationErrors, nil
	}

	var jsonA map[string]interface{}
	err := json.Unmarshal(*requestRecord.Body, &jsonA)
	if err != nil {
		return &validationErrors, err
	}

	bodyValidationResult := jsonValidate(jsonA, assertConfig.BodyJson)

	if !bodyValidationResult {
		validationErrors = append(validationErrors, ValidationError{Validation_error_code_body_mismatch, []string{}})
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
				Metadata: []string{
					headerKey,
				},
			})

			continue
		}

		if !utils.ListsEqual[string](header, headerB) {
			validationErrors = append(validationErrors, ValidationError{
				Code: Validation_error_code_header_value_mismatch,
				Metadata: []string{
					headerKey,
					strings.Join(header, ","),
					strings.Join(headerB, ","),
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
