package types

import (
	"net/http"
)

type State struct {
	RequestRecordDirectoryPath string
}

type EndpointConfig struct {
	Route   string `json:"route"`
	Method  string `json:"method"`
	Content string `json:"content"`
}

type RequestRecord struct {
	Route   string      `json:"route"`
	Headers http.Header `json:"headers"`
	Body    *[]byte     `json:"body"`
}

type MockFs interface {
	StoreRequestRecord(r *http.Request, endpointConfig *EndpointConfig) error
	GetRecordsMatchingRoute(route string) ([]*RequestRecord, error)
}
