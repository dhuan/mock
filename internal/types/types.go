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

type EndpointConfig struct {
	Route   string                `json:"route"`
	Method  string                `json:"method"`
	Content EndpointConfigContent `json:"content"`
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

type EndpointConfigContent []byte

func (this *EndpointConfigContent) UnmarshalJSON(data []byte) (err error) {
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
