package mock_test

import (
	"testing"

	"github.com/dhuan/mock/internal/mock"
	"github.com/dhuan/mock/internal/types"
	"github.com/stretchr/testify/assert"
)

var readFileMockReturn = []byte("")

func readFileMock(name string) ([]byte, error) {
	return readFileMockReturn, nil
}

func Test_ResolveEndpointResponse_GettingResponse(t *testing.T) {
	state := types.State{
		RequestRecordDirectoryPath: "/path/to/somewhere",
		ConfigFolderPath:           "/path/to/somewhere",
	}
	endpointConfig := types.EndpointConfig{
		Route:   "foo/bar",
		Method:  "post",
		Content: []byte(`{"foo":"bar"}`),
		Headers: map[string]string{},
	}

	response, endpointContentType, _ := mock.ResolveEndpointResponse(readFileMock, &state, &endpointConfig)

	assert.Equal(
		t,
		`{"foo":"bar"}`,
		string(response),
	)

	assert.Equal(
		t,
		types.Endpoint_content_type_json,
		endpointContentType,
	)
}

func Test_ResolveEndpointResponse_EndpointWithResponseByFile(t *testing.T) {
	readFileMockReturn = []byte("Hello world!")
	state := types.State{
		RequestRecordDirectoryPath: "/path/to/somewhere",
		ConfigFolderPath:           "/path/to/somewhere",
	}
	endpointConfig := types.EndpointConfig{
		Route:   "foo/bar",
		Method:  "post",
		Content: []byte(`file:./response_foobar`),
		Headers: map[string]string{},
	}

	response, endpointContentType, _ := mock.ResolveEndpointResponse(readFileMock, &state, &endpointConfig)

	assert.Equal(
		t,
		"Hello world!",
		string(response),
	)

	assert.Equal(
		t,
		types.Endpoint_content_type_file,
		endpointContentType,
	)
}
