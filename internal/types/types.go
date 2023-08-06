package types

import (
	"net/http"
	"strings"

	"github.com/dhuan/mock/internal/utils"
	. "github.com/dhuan/mock/pkg/mock"
)

type State struct {
	ListenPort                 string
	RequestRecordDirectoryPath string
	ConfigFolderPath           string
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

type MiddlewareConfig struct {
	Exec       string         `json:"exec"`
	Type       MiddlewareType `json:"type"`
	RouteMatch string         `json:"route_match"`
	Condition  *Condition     `json:"condition"`
}

type RequestRecord struct {
	Route             string            `json:"route"`
	Querystring       string            `json:"querystring"`
	QuerystringParsed map[string]string `json:"querystring_parsed"`
	Method            string            `json:"method"`
	Host              string            `json:"host"`
	Https             bool              `json:"https"`
	Headers           http.Header       `json:"headers"`
	Body              *[]byte           `json:"body"`
	RouteParams       map[string]string `json:"route_params"`
}

type MockFs interface {
	StoreRequestRecord(requestRecord *RequestRecord, endpointConfig *EndpointConfig) error
	GetRecordsMatchingRoute(route string) ([]RequestRecord, error)
	RemoveAllRequestRecords() error
}

type EndpointConfigResponse []byte

func (this *EndpointConfigResponse) UnmarshalJSON(data []byte) (err error) {
	if utils.BeginsWith(string(data), `"file:`) {
		*this = []byte(utils.Unquote(string(data)))

		return nil
	}

	if utils.BeginsWith(string(data), `"sh:`) {
		*this = []byte(utils.Unquote(string(data)))

		return nil
	}

	if utils.BeginsWith(string(data), `"exec:`) {
		*this = []byte(utils.Unquote(string(data)))

		return nil
	}

	if utils.BeginsWith(string(data), `"fs:`) {
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
	Endpoint_content_type_shell
	Endpoint_content_type_exec
	Endpoint_content_type_fileserver
	Endpoint_content_type_json
	Endpoint_content_type_plaintext
	Endpoint_content_type_unknown
)

type MiddlewareType int

const (
	MiddlewareType_Unknown MiddlewareType = iota
	MiddlewareType_BeforeResponse
)

func (this *MiddlewareType) UnmarshalJSON(data []byte) (err error) {
	text := utils.Unquote(string(data))

	if text == "before_response" {
		*this = MiddlewareType_BeforeResponse

		return
	}

	*this = MiddlewareType_BeforeResponse

	return nil
}

var Middleware_type_code_encoding_map = map[MiddlewareType]string{
	MiddlewareType_Unknown:        "unknown",
	MiddlewareType_BeforeResponse: "before_response",
}

type ReadFileFunc = func(filePath string) ([]byte, error)
