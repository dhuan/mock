package mock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dhuan/mock/internal/utils"
)

var validation_error_code_encoding_map = map[ValidationErrorCode]string{
	ValidationErrorCode_HeaderValueMismatch:              "header_value_mismatch",
	ValidationErrorCode_NoCall:                           "no_call",
	ValidationErrorCode_HeaderNotIncluded:                "header_not_included",
	ValidationErrorCode_BodyMismatch:                     "body_mismatch",
	ValidationErrorCode_RequestHasNoBodyContent:          "request_has_no_body_content",
	ValidationErrorCode_MethodMismatch:                   "method_mismatch",
	ValidationErrorCode_FormKeyDoesNotExist:              "form_key_does_not_exist",
	ValidationErrorCode_FormValueMismatch:                "form_value_mismatch",
	ValidationErrorCode_NthOutOfRange:                    "nth_out_of_range",
	ValidationErrorCode_RequestHasNoQuerystring:          "request_has_no_querystring",
	ValidationErrorCode_QuerystringMismatch:              "querystring_mismatch",
	ValidationErrorCode_QuerystringKeyNotSet:             "querystring_key_not_set",
	ValidationErrorCode_RequestHasNoBody:                 "request_has_no_body",
	ValidationErrorCode_NthMismatch:                      "nth_mismatch",
	ValidationErrorCode_RouteParamValueMismatch:          "route_param_value_mismatch",
	ValidationErrorCode_RouteParamDoesNotExistInEndpoint: "route_param_does_not_exist_in_endpoint",
}

type AssertHeader map[string][]string

type AssertOptions struct {
	Route     string     `json:"route"`
	Nth       int        `json:"nth"`
	Condition *Condition `json:"condition"`
}

type ValidationError struct {
	Code     ValidationErrorCode `json:"code"`
	Metadata map[string]string   `json:"metadata"`
}

type ValidationErrorCode int

const (
	ValidationErrorCode_Unknown ValidationErrorCode = iota
	ValidationErrorCode_NoCall
	ValidationErrorCode_MethodMismatch
	ValidationErrorCode_HeaderNotIncluded
	ValidationErrorCode_HeaderValueMismatch
	ValidationErrorCode_BodyMismatch
	ValidationErrorCode_RequestHasNoBodyContent
	ValidationErrorCode_FormKeyDoesNotExist
	ValidationErrorCode_FormValueMismatch
	ValidationErrorCode_NthOutOfRange
	ValidationErrorCode_RequestHasNoQuerystring
	ValidationErrorCode_QuerystringMismatch
	ValidationErrorCode_QuerystringKeyNotSet
	ValidationErrorCode_RequestHasNoBody
	ValidationErrorCode_NthMismatch
	ValidationErrorCode_RouteParamValueMismatch
	ValidationErrorCode_RouteParamDoesNotExistInEndpoint
)

func (vec *ValidationErrorCode) MarshalJSON() ([]byte, error) {
	encodingMapPrepared := utils.MapMapValueOnly(
		validation_error_code_encoding_map,
		utils.WrapIn(`"`),
	)

	return utils.MarshalJsonHelper(
		encodingMapPrepared,
		"Failed to parse Validation Error Code: %d",
		vec,
	)
}

func (vec *ValidationErrorCode) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalJsonHelper(
		vec,
		validation_error_code_encoding_map,
		data,
		"Failed to parse Validation Error Code: %s",
	)
}

type MockConfig struct {
	Url string
}

func Init(url string) *MockConfig {
	return &MockConfig{
		Url: url,
	}
}

type AssertResponse struct {
	ValidationErrors []ValidationError `json:"validation_errors"`
}

func Assert(config *MockConfig, assertOptions *AssertOptions) ([]ValidationError, error) {
	bodyJson, err := json.Marshal(assertOptions)
	if err != nil {
		return make([]ValidationError, 0, 0), err
	}

	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://%s/__mock__/assert", config.Url),
		bytes.NewBuffer([]byte(bodyJson)),
	)
	if err != nil {
		return make([]ValidationError, 0, 0), err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return make([]ValidationError, 0, 0), err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return make([]ValidationError, 0, 0), err
	}

	var responseParsed AssertResponse
	err = json.Unmarshal(responseBody, &responseParsed)
	if err != nil {
		return make([]ValidationError, 0, 0), err
	}

	return responseParsed.ValidationErrors, nil
}

func ToReadableError(validationErrors []ValidationError) string {
	result := make([]string, 0)

	for i := range validationErrors {
		validationErrorEncoded, ok := validation_error_code_encoding_map[validationErrors[i].Code]

		if !ok {
			return "Failed to parse!"
		}

		metadataKeys := utils.GetSortedKeys(validationErrors[i].Metadata)
		metadataEncoded := make([]string, len(metadataKeys))
		for i2, metadataKey := range metadataKeys {
			metadataEncoded[i2] = fmt.Sprintf("%s: %s", metadataKey, validationErrors[i].Metadata[metadataKey])
		}

		errorMessage := fmt.Sprintf(
			"Error: %s",
			validationErrorEncoded,
		)

		if len(metadataKeys) > 0 {
			errorMessage = fmt.Sprintf(
				"Error: %s\n%s",
				validationErrorEncoded,
				strings.Join(metadataEncoded, "\n"),
			)
		}

		result = append(result, errorMessage)
	}

	return strings.Join(result, "\n\n")
}
