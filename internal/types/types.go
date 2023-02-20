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

type RequestRecord struct {
	Route       string      `json:"route"`
	Querystring string      `json:"querystring"`
	Method      string      `json:"method"`
	Headers     http.Header `json:"headers"`
	Body        *[]byte     `json:"body"`
}

type MockFs interface {
	StoreRequestRecord(r *http.Request, requestBody []byte, endpointConfig *EndpointConfig) error
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
	Endpoint_content_type_json
	Endpoint_content_type_plaintext
	Endpoint_content_type_unknown
)
