package mock_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/dhuan/mock/internal/mock"
	"github.com/dhuan/mock/internal/types"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

var readFileMockReturn = []byte("")

var requestMock = &http.Request{URL: &url.URL{RawQuery: "hello=world"}}

type osMock struct {
	testifymock.Mock
}

func (this *osMock) ReadFile(name string) ([]byte, error) {
	args := this.Called(name)

	return args.Get(0).([]byte), nil
}

func Test_ResolveEndpointResponse_GettingResponse_Json(t *testing.T) {
	osMockInstance := osMock{}
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

	response, endpointContentType, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, &state, &endpointConfig)

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

func Test_ResolveEndpointResponse_GettingResponse_PlainText(t *testing.T) {
	osMockInstance := osMock{}
	state := types.State{
		RequestRecordDirectoryPath: "/path/to/somewhere",
		ConfigFolderPath:           "/path/to/somewhere",
	}
	endpointConfig := types.EndpointConfig{
		Route:   "foo/bar",
		Method:  "post",
		Content: []byte(`Hello world!`),
		Headers: map[string]string{},
	}

	response, endpointContentType, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, &state, &endpointConfig)

	assert.Equal(
		t,
		`Hello world!`,
		string(response),
	)

	assert.Equal(
		t,
		types.Endpoint_content_type_plaintext,
		endpointContentType,
	)
}

func Test_ResolveEndpointResponse_EndpointWithResponseByFile(t *testing.T) {
	osMockInstance := osMock{}
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

	osMockInstance.On("ReadFile", "/path/to/somewhere/./response_foobar").Return([]byte("Hello world!"), nil)

	response, endpointContentType, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, &state, &endpointConfig)

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

func Test_ResolveEndpointResponse_WithQueryStringCondition(t *testing.T) {
	osMockInstance := osMock{}
	state := types.State{
		RequestRecordDirectoryPath: "/path/to/somewhere",
		ConfigFolderPath:           "/path/to/somewhere",
	}
	endpointConfig := types.EndpointConfig{
		Route:   "foo/bar",
		Method:  "post",
		Headers: map[string]string{},
		Content: []byte(`file:./response_foobar`),
		ResponseIf: []types.ResponseIf{
			types.ResponseIf{
				Response: []byte(`{"result": "response_one"}`),
				QuerystringMatches: []types.Kv{
					types.Kv{Key: "foo", Value: "bar"},
				},
			},
			types.ResponseIf{
				Response: []byte(`{"result": "response_two"}`),
				QuerystringMatches: []types.Kv{
					types.Kv{Key: "hello", Value: "world"},
				},
			},
		},
	}

	response, endpointContentType, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, &state, &endpointConfig)

	assert.Equal(
		t,
		`{"result":"response_two"}`,
		string(response),
	)

	assert.Equal(
		t,
		types.Endpoint_content_type_json,
		endpointContentType,
	)
}

func Test_ResolveEndpointResponse_WithQueryStringCondition_FallbackResponse(t *testing.T) {
	requestMock = &http.Request{URL: &url.URL{RawQuery: "hello=WORLD"}}
	osMockInstance := osMock{}
	state := types.State{
		RequestRecordDirectoryPath: "/path/to/somewhere",
		ConfigFolderPath:           "/path/to/somewhere",
	}
	endpointConfig := types.EndpointConfig{
		Route:   "foo/bar",
		Method:  "post",
		Headers: map[string]string{},
		Content: []byte(`Fallback response!`),
		ResponseIf: []types.ResponseIf{
			types.ResponseIf{
				Response: []byte(`{"result": "response_one"}`),
				QuerystringMatches: []types.Kv{
					types.Kv{Key: "foo", Value: "bar"},
				},
			},
			types.ResponseIf{
				Response: []byte(`{"result": "response_two"}`),
				QuerystringMatches: []types.Kv{
					types.Kv{Key: "hello", Value: "world"},
				},
			},
		},
	}

	response, endpointContentType, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, &state, &endpointConfig)

	assert.Equal(
		t,
		`Fallback response!`,
		string(response),
	)

	assert.Equal(
		t,
		types.Endpoint_content_type_plaintext,
		endpointContentType,
	)
}

func Test_ResolveEndpointResponse_WithExactQueryStringCondition_Matching(t *testing.T) {
	requestMock = &http.Request{URL: &url.URL{RawQuery: "hello=world&foo=bar"}}
	osMockInstance := osMock{}
	state := types.State{
		RequestRecordDirectoryPath: "/path/to/somewhere",
		ConfigFolderPath:           "/path/to/somewhere",
	}
	endpointConfig := types.EndpointConfig{
		Route:   "foo/bar",
		Method:  "post",
		Headers: map[string]string{},
		Content: []byte(`Fallback response!`),
		ResponseIf: []types.ResponseIf{
			types.ResponseIf{
				Response: []byte(`response_one`),
				QuerystringMatchesExact: []types.Kv{
					types.Kv{Key: "hello", Value: "world"},
				},
			},
			types.ResponseIf{
				Response: []byte(`response_two`),
				QuerystringMatchesExact: []types.Kv{
					types.Kv{Key: "hello", Value: "world"},
					types.Kv{Key: "foo", Value: "bar"},
				},
			},
		},
	}

	response, endpointContentType, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, &state, &endpointConfig)

	assert.Equal(
		t,
		`response_two`,
		string(response),
	)

	assert.Equal(
		t,
		types.Endpoint_content_type_plaintext,
		endpointContentType,
	)
}
