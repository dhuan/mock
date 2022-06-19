package types

import (
	"net/http"
	"strings"

	"github.com/dhuan/mock/internal/utils"
)

type State struct {
	RequestRecordDirectoryPath string
	ConfigFolderPath           string
}

type ResponseIf struct {
	Response                EndpointConfigResponse `json:"response"`
	QuerystringMatches      []Kv                   `json:"querystring_matches"`
	QuerystringMatchesExact []Kv                   `json:"querystring_matches_exact"`
}

type Kv struct {
	Key   string
	Value string
}

type EndpointConfig struct {
	Route      string                 `json:"route"`
	Method     string                 `json:"method"`
	Response   EndpointConfigResponse `json:"response"`
	Headers    map[string]string      `json:"response_headers"`
	ResponseIf []ResponseIf           `json:"response_if"`
}

type RequestRecord struct {
	Route   string      `json:"route"`
	Method  string      `json:"method"`
	Headers http.Header `json:"headers"`
	Body    *[]byte     `json:"body"`
}

type MockFs interface {
	StoreRequestRecord(r *http.Request, endpointConfig *EndpointConfig) error
	GetRecordsMatchingRoute(route string) ([]*RequestRecord, error)
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
