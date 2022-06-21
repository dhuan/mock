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
		Route:    "foo/bar",
		Method:   "post",
		Response: []byte(`{"foo":"bar"}`),
		Headers:  map[string]string{},
	}

	response, endpointContentType, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, &state, &endpointConfig)

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
		Route:    "foo/bar",
		Method:   "post",
		Response: []byte(`Hello world!`),
		Headers:  map[string]string{},
	}

	response, endpointContentType, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, &state, &endpointConfig)

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
		Route:    "foo/bar",
		Method:   "post",
		Response: []byte(`file:./response_foobar`),
		Headers:  map[string]string{},
	}

	osMockInstance.On("ReadFile", "/path/to/somewhere/./response_foobar").Return([]byte("Hello world!"), nil)

	response, endpointContentType, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, &state, &endpointConfig)

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

func Test_ResolveEndpointResponse_DefaultResponseStatusCode(t *testing.T) {
	osMockInstance := osMock{}
	state := types.State{
		RequestRecordDirectoryPath: "/path/to/somewhere",
		ConfigFolderPath:           "/path/to/somewhere",
	}
	endpointConfig := types.EndpointConfig{
		Route:    "foo/bar",
		Method:   "post",
		Response: []byte(`Hello world!`),
		Headers:  map[string]string{},
	}

	_, _, responseStatusCode, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, &state, &endpointConfig)

	assert.Equal(t, 200, responseStatusCode)
}

func Test_ResolveEndpointResponse_ResponseStatusCode(t *testing.T) {
	osMockInstance := osMock{}
	state := types.State{
		RequestRecordDirectoryPath: "/path/to/somewhere",
		ConfigFolderPath:           "/path/to/somewhere",
	}
	endpointConfig := types.EndpointConfig{
		Route:              "foo/bar",
		Method:             "post",
		Response:           []byte(`Hello world!`),
		ResponseStatusCode: 201,
		Headers:            map[string]string{},
	}

	_, _, responseStatusCode, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, &state, &endpointConfig)

	assert.Equal(t, 201, responseStatusCode)
}

func Test_ResolveEndpointResponse_WithQueryStringCondition(t *testing.T) {
	osMockInstance := osMock{}
	state := types.State{
		RequestRecordDirectoryPath: "/path/to/somewhere",
		ConfigFolderPath:           "/path/to/somewhere",
	}
	endpointConfig := types.EndpointConfig{
		Route:    "foo/bar",
		Method:   "post",
		Headers:  map[string]string{},
		Response: []byte(`file:./response_foobar`),
		ResponseIf: []types.ResponseIf{
			types.ResponseIf{
				Response:           []byte(`{"result": "response_one"}`),
				ResponseStatusCode: 202,
				Condition: &types.Condition{
					Type:  types.ConditionType_QuerystringMatch,
					Key:   "foo",
					Value: "bar",
				},
			},
			types.ResponseIf{
				Response:           []byte(`{"result": "response_two"}`),
				ResponseStatusCode: 203,
				Condition: &types.Condition{
					Type:  types.ConditionType_QuerystringMatch,
					Key:   "hello",
					Value: "world",
				},
			},
		},
	}

	response, endpointContentType, responseStatusCode, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, &state, &endpointConfig)

	assert.Equal(t, 203, responseStatusCode)

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
		Route:    "foo/bar",
		Method:   "post",
		Headers:  map[string]string{},
		Response: []byte(`Fallback response!`),
		ResponseIf: []types.ResponseIf{
			types.ResponseIf{
				Response: []byte(`{"result": "response_one"}`),
				Condition: &types.Condition{
					Type:  types.ConditionType_QuerystringMatch,
					Key:   "foo",
					Value: "bar",
				},
			},
			types.ResponseIf{
				Response: []byte(`{"result": "response_two"}`),
				Condition: &types.Condition{
					Type:  types.ConditionType_QuerystringMatch,
					Key:   "hello",
					Value: "world",
				},
			},
		},
	}

	response, endpointContentType, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, &state, &endpointConfig)

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

func Test_ResolveEndpointResponse_WithAndChaining(t *testing.T) {
	requestMock = &http.Request{URL: &url.URL{RawQuery: "hello=world&foo=bar"}}
	osMockInstance := osMock{}
	state := types.State{
		RequestRecordDirectoryPath: "/path/to/somewhere",
		ConfigFolderPath:           "/path/to/somewhere",
	}
	endpointConfig := types.EndpointConfig{
		Route:    "foo/bar",
		Method:   "post",
		Headers:  map[string]string{},
		Response: []byte(`Fallback response!`),
		ResponseIf: []types.ResponseIf{
			types.ResponseIf{
				Response: []byte(`response_two`),
				Condition: &types.Condition{
					Type:  types.ConditionType_QuerystringMatch,
					Key:   "hello",
					Value: "world",
					And: &types.Condition{
						Type:  types.ConditionType_QuerystringMatch,
						Key:   "foo",
						Value: "BAR",
					},
				},
			},
			types.ResponseIf{
				Response: []byte(`response_two`),
				Condition: &types.Condition{
					Type:  types.ConditionType_QuerystringMatch,
					Key:   "hello",
					Value: "world",
					And: &types.Condition{
						Type:  types.ConditionType_QuerystringMatch,
						Key:   "foo",
						Value: "bar",
					},
				},
			},
		},
	}

	response, endpointContentType, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, &state, &endpointConfig)

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
