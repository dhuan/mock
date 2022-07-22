package types

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dhuan/mock/internal/utils"
)

type ConditionType int

const (
	ConditionType_None ConditionType = iota
	ConditionType_QuerystringMatch
	ConditionType_QuerystringExactMatch
	ConditionType_FormMatch
)

func (this *ConditionType) UnmarshalJSON(data []byte) (err error) {
	conditionTypeText := utils.Unquote(string(data))

	if conditionTypeText == "querystring_match" {
		*this = ConditionType_QuerystringMatch

		return nil
	}

	if conditionTypeText == "querystring_exact_match" {
		*this = ConditionType_QuerystringExactMatch

		return nil
	}

	if conditionTypeText == "form_match" {
		*this = ConditionType_FormMatch

		return nil
	}

	return errors.New(fmt.Sprintf("Failed to parse Condition Type: %s", conditionTypeText))
}

type State struct {
	RequestRecordDirectoryPath string
	ConfigFolderPath           string
}

type Condition struct {
	Type      ConditionType          `json:"type"`
	Key       string                 `json:"key"`
	Value     string                 `json:"value"`
	KeyValues map[string]interface{} `json:"key_values"`
	And       *Condition             `json:"and"`
	Or        *Condition             `json:"or"`
}

type ResponseIf struct {
	Response           EndpointConfigResponse `json:"response"`
	ResponseStatusCode int                    `json:"response_status_code"`
	Condition          *Condition             `json:"condition"`
	Headers            map[string]string      `json:"response_headers"`
}

type EndpointConfig struct {
	Route              string                 `json:"route"`
	Method             string                 `json:"method"`
	Response           EndpointConfigResponse `json:"response"`
	ResponseStatusCode int                    `json:"response_status_code"`
	ResponseIf         []ResponseIf           `json:"response_if"`
	Headers            map[string]string      `json:"response_headers"`
	HeadersBase        map[string]string      `json:"response_headers_base"`
}

type RequestRecord struct {
	Route       string      `json:"route"`
	Querystring string      `json:"querystring"`
	Method      string      `json:"method"`
	Headers     http.Header `json:"headers"`
	Body        *[]byte     `json:"body"`
}

type MockFs interface {
	StoreRequestRecord(r *http.Request, requestBody []byte, endpointConfig *EndpointConfig) error
	GetRecordsMatchingRoute(route string) ([]*RequestRecord, error)
	RemoveAllRequestRecords() error
}

type EndpointConfigResponse []byte

func (this *EndpointConfigResponse) UnmarshalJSON(data []byte) (err error) {
	if utils.BeginsWith(string(data), `"file:`) {
		*this = []byte(utils.Unquote(string(data)))

		return nil
	}

	if strings.Index(string(data), "{") == 0 {
		*this = data

		return nil
	}

	*this = data

	return nil
}

type Endpoint_content_type int

const (
	Endpoint_content_type_file Endpoint_content_type = iota
	Endpoint_content_type_json
	Endpoint_content_type_plaintext
	Endpoint_content_type_unknown
)
