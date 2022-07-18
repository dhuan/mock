package mock

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

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

func (this *AssertType) MarshalJSON() ([]byte, error) {
	if *this == AssertType_HeaderMatch {
		return []byte(`"header_match"`), nil
	}

	if *this == AssertType_MethodMatch {
		return []byte(`"method_match"`), nil
	}

	if *this == AssertType_JsonBodyMatch {
		return []byte(`"json_body_match"`), nil
	}

	if *this == AssertType_FormMatch {
		return []byte(`"form_match"`), nil
	}

	return []byte(""), errors.New(fmt.Sprintf("Failed to parse Assert Type: %d", *this))
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
	Code     string            `json:"code"`
	Metadata map[string]string `json:"metadata"`
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
