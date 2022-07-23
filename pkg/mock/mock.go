package mock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dhuan/mock/internal/utils"
)

var assert_type_encoding_map = map[AssertType]string{
	AssertType_HeaderMatch:   "header_match",
	AssertType_MethodMatch:   "method_match",
	AssertType_JsonBodyMatch: "json_body_match",
	AssertType_FormMatch:     "form_match",
}

type AssertType int

const (
	AssertType_None AssertType = iota
	AssertType_HeaderMatch
	AssertType_MethodMatch
	AssertType_JsonBodyMatch
	AssertType_FormMatch
)

func (this *AssertType) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalJsonHelper[AssertType](this, assert_type_encoding_map, data, "Failed to parse Assert Type: %s")
}

func (this *AssertType) MarshalJSON() ([]byte, error) {
	encodingMapPrepared := utils.MapMapValueOnly[AssertType, string, string](assert_type_encoding_map, utils.WrapIn(`"`))

	return utils.MarshalJsonHelper[AssertType](
		encodingMapPrepared,
		"Failed to parse Assert Type: %d",
		this,
	)
}

var validation_error_code_encoding_map = map[ValidationErrorCode]string{
	ValidationErrorCode_HeaderValueMismatch:     "header_value_mismatch",
	ValidationErrorCode_NoCall:                  "no_call",
	ValidationErrorCode_HeaderNotIncluded:       "header_not_included",
	ValidationErrorCode_BodyMismatch:            "body_mismatch",
	ValidationErrorCode_RequestHasNoBodyContent: "request_has_no_body_content",
	ValidationErrorCode_MethodMismatch:          "method_mismatch",
	ValidationErrorCode_FormKeyDoesNotExist:     "form_key_does_not_exist",
	ValidationErrorCode_FormValueMismatch:       "form_value_mismatch",
}

type AssertHeader map[string][]string

type AssertConfig struct {
	Route  string         `json:"route"`
	Nth    int            `json:"nth"`
	Method string         `json:"method"`
	Assert *AssertOptions `json:"assert"`
}

type AssertOptions struct {
	Type      AssertType             `json:"type"`
	Data      map[string]interface{} `json:"data"`
	KeyValues map[string]interface{} `json:"key_values"`
	Key       string                 `json:"key"`
	Value     string                 `json:"value"`
	And       *AssertOptions         `json:"and"`
	Or        *AssertOptions         `json:"or"`
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
)

func (this *ValidationErrorCode) MarshalJSON() ([]byte, error) {
	encodingMapPrepared := utils.MapMapValueOnly[ValidationErrorCode, string, string](
		validation_error_code_encoding_map,
		utils.WrapIn(`"`),
	)

	return utils.MarshalJsonHelper[ValidationErrorCode](
		encodingMapPrepared,
		"Failed to parse Validation Error Code: %d",
		this,
	)
}

func (this *ValidationErrorCode) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalJsonHelper[ValidationErrorCode](
		this,
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

func Assert(config *MockConfig, assertConfig *AssertConfig) []ValidationError {
	bodyJson, err := json.Marshal(assertConfig)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://%s/__mock__/assert", config.Url),
		bytes.NewBuffer([]byte(bodyJson)),
	)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var responseParsed AssertResponse
	err = json.Unmarshal(responseBody, &responseParsed)
	if err != nil {
		panic(err)
	}

	return responseParsed.ValidationErrors
}
