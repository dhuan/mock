package mock_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/dhuan/mock/internal/mock"
	"github.com/dhuan/mock/internal/types"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

var readFileMockReturn = []byte("")

var requestMock = &http.Request{URL: &url.URL{RawQuery: "hello=world"}}
var requestBody = []byte("")

var state types.State = types.State{
	RequestRecordDirectoryPath: "/path/to/somewhere",
	ConfigFolderPath:           "/path/to/somewhere",
}

type osMock struct {
	testifymock.Mock
}

func (this *osMock) ReadFile(name string) ([]byte, error) {
	args := this.Called(name)

	return args.Get(0).([]byte), nil
}

func Test_ResolveEndpointResponse_GettingResponse_Json(t *testing.T) {
	osMockInstance := osMock{}
	endpointConfig := types.EndpointConfig{
		Route:    "foo/bar",
		Method:   "post",
		Response: []byte(`{"foo":"bar"}`),
		Headers:  map[string]string{},
	}

	response, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, requestBody, &state, &endpointConfig)

	assert.Equal(
		t,
		`{"foo":"bar"}`,
		string(response.Body),
	)

	assert.Equal(
		t,
		types.Endpoint_content_type_json,
		response.EndpointContentType,
	)
}

func Test_ResolveEndpointResponse_GettingResponse_PlainText(t *testing.T) {
	osMockInstance := osMock{}
	endpointConfig := types.EndpointConfig{
		Route:    "foo/bar",
		Method:   "post",
		Response: []byte(`Hello world!`),
		Headers:  map[string]string{},
	}

	response, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, requestBody, &state, &endpointConfig)

	assert.Equal(
		t,
		`Hello world!`,
		string(response.Body),
	)

	assert.Equal(
		t,
		types.Endpoint_content_type_plaintext,
		response.EndpointContentType,
	)
}

func Test_ResolveEndpointResponse_EndpointWithResponseByFile(t *testing.T) {
	osMockInstance := osMock{}
	endpointConfig := types.EndpointConfig{
		Route:    "foo/bar",
		Method:   "post",
		Response: []byte(`file:./response_foobar`),
		Headers:  map[string]string{},
	}

	osMockInstance.On("ReadFile", "/path/to/somewhere/./response_foobar").Return([]byte("Hello world!"), nil)

	response, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, requestBody, &state, &endpointConfig)

	assert.Equal(
		t,
		"Hello world!",
		string(response.Body),
	)

	assert.Equal(
		t,
		types.Endpoint_content_type_file,
		response.EndpointContentType,
	)
}

func Test_ResolveEndpointResponse_DefaultResponseStatusCode(t *testing.T) {
	osMockInstance := osMock{}
	endpointConfig := types.EndpointConfig{
		Route:    "foo/bar",
		Method:   "post",
		Response: []byte(`Hello world!`),
		Headers:  map[string]string{},
	}

	response, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, requestBody, &state, &endpointConfig)

	assert.Equal(t, 200, response.StatusCode)
}

func Test_ResolveEndpointResponse_ResponseStatusCode(t *testing.T) {
	osMockInstance := osMock{}
	endpointConfig := types.EndpointConfig{
		Route:              "foo/bar",
		Method:             "post",
		Response:           []byte(`Hello world!`),
		ResponseStatusCode: 201,
		Headers:            map[string]string{},
	}

	response, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, requestBody, &state, &endpointConfig)

	assert.Equal(t, 201, response.StatusCode)
}

func Test_ResolveEndpointResponse_WithQueryStringCondition(t *testing.T) {
	osMockInstance := osMock{}
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

	response, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, requestBody, &state, &endpointConfig)

	assert.Equal(t, 203, response.StatusCode)

	assert.Equal(
		t,
		`{"result":"response_two"}`,
		string(response.Body),
	)

	assert.Equal(
		t,
		types.Endpoint_content_type_json,
		response.EndpointContentType,
	)
}

func Test_ResolveEndpointResponse_WithQueryStringCondition_FallbackResponse(t *testing.T) {
	requestMock = &http.Request{URL: &url.URL{RawQuery: "hello=WORLD"}}
	osMockInstance := osMock{}
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

	response, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, requestBody, &state, &endpointConfig)

	assert.Equal(
		t,
		`Fallback response!`,
		string(response.Body),
	)

	assert.Equal(
		t,
		types.Endpoint_content_type_plaintext,
		response.EndpointContentType,
	)
}

func Test_ResolveEndpointResponse_WithQueryStringCondition_WithMultipleValues(t *testing.T) {
	requestMock = &http.Request{URL: &url.URL{RawQuery: "hello=world&foo=bar"}}
	osMockInstance := osMock{}
	endpointConfig := types.EndpointConfig{
		Route:    "foo/bar",
		Method:   "post",
		Headers:  map[string]string{},
		Response: []byte(`default_response`),
		ResponseIf: []types.ResponseIf{
			types.ResponseIf{
				Response:           []byte(`response_one`),
				ResponseStatusCode: 202,
				Condition: &types.Condition{
					Type: types.ConditionType_QuerystringMatch,
					KeyValues: map[string]interface{}{
						"foo":   "not_bar",
						"hello": "world",
					},
				},
			},
			types.ResponseIf{
				Response:           []byte(`response_two`),
				ResponseStatusCode: 203,
				Condition: &types.Condition{
					Type: types.ConditionType_QuerystringMatch,
					KeyValues: map[string]interface{}{
						"foo":   "bar",
						"hello": "world",
					},
				},
			},
		},
	}

	response, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, requestBody, &state, &endpointConfig)

	assert.Equal(t, 203, response.StatusCode)

	assert.Equal(
		t,
		`response_two`,
		string(response.Body),
	)

	assert.Equal(
		t,
		types.Endpoint_content_type_plaintext,
		response.EndpointContentType,
	)
}

func Test_ResolveEndpointResponse_WithQueryStringExactCondition_FallingBackToDefault(t *testing.T) {
	requestMock = &http.Request{URL: &url.URL{RawQuery: "hello=world&foo=bar"}}
	osMockInstance := osMock{}
	endpointConfig := types.EndpointConfig{
		Route:    "foo/bar",
		Method:   "post",
		Headers:  map[string]string{},
		Response: []byte(`default_response`),
		ResponseIf: []types.ResponseIf{
			types.ResponseIf{
				Response:           []byte(`response_one`),
				ResponseStatusCode: 202,
				Condition: &types.Condition{
					Type: types.ConditionType_QuerystringMatch,
					KeyValues: map[string]interface{}{
						"some_key": "some_value",
					},
				},
			},
			types.ResponseIf{
				Response:           []byte(`response_two`),
				ResponseStatusCode: 203,
				Condition: &types.Condition{
					Type: types.ConditionType_QuerystringExactMatch,
					KeyValues: map[string]interface{}{
						"hello": "world",
					},
				},
			},
		},
	}

	response, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, requestBody, &state, &endpointConfig)

	assert.Equal(
		t,
		`default_response`,
		string(response.Body),
	)

	assert.Equal(
		t,
		types.Endpoint_content_type_plaintext,
		response.EndpointContentType,
	)
}

func Test_ResolveEndpointResponse_WithQueryStringExactCondition_ResolvingToConditionalResponse(t *testing.T) {
	requestMock = &http.Request{URL: &url.URL{RawQuery: "hello=world"}}
	osMockInstance := osMock{}
	endpointConfig := types.EndpointConfig{
		Route:    "foo/bar",
		Method:   "post",
		Headers:  map[string]string{},
		Response: []byte(`default_response`),
		ResponseIf: []types.ResponseIf{
			types.ResponseIf{
				Response:           []byte(`response_one`),
				ResponseStatusCode: 202,
				Condition: &types.Condition{
					Type: types.ConditionType_QuerystringMatch,
					KeyValues: map[string]interface{}{
						"some_key": "some_value",
					},
				},
			},
			types.ResponseIf{
				Response:           []byte(`response_two`),
				ResponseStatusCode: 203,
				Condition: &types.Condition{
					Type: types.ConditionType_QuerystringExactMatch,
					KeyValues: map[string]interface{}{
						"hello": "world",
					},
				},
			},
		},
	}

	response, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, requestBody, &state, &endpointConfig)

	assert.Equal(
		t,
		`response_two`,
		string(response.Body),
	)

	assert.Equal(
		t,
		types.Endpoint_content_type_plaintext,
		response.EndpointContentType,
	)
}

func Test_ResolveEndpointResponse_WithAndChaining(t *testing.T) {
	requestMock = &http.Request{URL: &url.URL{RawQuery: "hello=world&foo=bar"}}
	osMockInstance := osMock{}
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

	response, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, requestBody, &state, &endpointConfig)

	assert.Equal(
		t,
		`response_two`,
		string(response.Body),
	)

	assert.Equal(
		t,
		types.Endpoint_content_type_plaintext,
		response.EndpointContentType,
	)
}

func Test_ResolveEndpointResponse_WithOrChaining(t *testing.T) {
	requestMock = &http.Request{URL: &url.URL{RawQuery: "hello=world"}}
	osMockInstance := osMock{}
	endpointConfig := types.EndpointConfig{
		Route:    "foo/bar",
		Method:   "post",
		Headers:  map[string]string{},
		Response: []byte(`Fallback response!`),
		ResponseIf: []types.ResponseIf{
			types.ResponseIf{
				Response: []byte(`or chaining!`),
				Condition: &types.Condition{
					Type:  types.ConditionType_QuerystringMatch,
					Key:   "foo",
					Value: "bar",
					Or: &types.Condition{
						Type:  types.ConditionType_QuerystringMatch,
						Key:   "hello",
						Value: "world",
					},
				},
			},
		},
	}

	response, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, requestBody, &state, &endpointConfig)

	assert.Equal(
		t,
		`or chaining!`,
		string(response.Body),
	)

	assert.Equal(
		t,
		types.Endpoint_content_type_plaintext,
		response.EndpointContentType,
	)
}

func Test_ResolveEndpointResponse_Headers_Match(t *testing.T) {
	osMockInstance := osMock{}
	endpointConfig := types.EndpointConfig{
		Route:    "foo/bar",
		Method:   "post",
		Response: []byte(`{"foo":"bar"}`),
		Headers: map[string]string{
			"Some-header-key": "Some header value",
		},
	}

	response, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, requestBody, &state, &endpointConfig)

	assert.Equal(
		t,
		map[string]string{
			"Some-header-key": "Some header value",
		},
		response.Headers,
	)
}

func Test_ResolveEndpointResponse_Headers_WithBase_Match(t *testing.T) {
	osMockInstance := osMock{}
	endpointConfig := types.EndpointConfig{
		Route:    "foo/bar",
		Method:   "post",
		Response: []byte(`{"foo":"bar"}`),
		HeadersBase: map[string]string{
			"Some-base-header-key": "Some base header value",
		},
		Headers: map[string]string{
			"Some-header-key": "Some header value",
		},
	}

	response, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, requestBody, &state, &endpointConfig)

	assert.Equal(
		t,
		map[string]string{
			"Some-base-header-key": "Some base header value",
			"Some-header-key":      "Some header value",
		},
		response.Headers,
	)
}

func Test_ResolveEndpointResponse_Headers_WithBase_WithConditionalResponse_Match(t *testing.T) {
	requestMock = &http.Request{URL: &url.URL{RawQuery: "foo=bar"}}
	osMockInstance := osMock{}
	endpointConfig := types.EndpointConfig{
		Route:    "foo/bar",
		Method:   "post",
		Response: []byte(`{"foo":"bar"}`),
		HeadersBase: map[string]string{
			"Some-base-header-key": "Some base header value",
		},
		Headers: map[string]string{
			"Some-header-key": "Some header value",
		},
		ResponseIf: []types.ResponseIf{
			types.ResponseIf{
				Response: []byte(`or chaining!`),
				Headers: map[string]string{
					"Another-header-key": "Another header value",
				},
				Condition: &types.Condition{
					Type:  types.ConditionType_QuerystringMatch,
					Key:   "foo",
					Value: "bar",
				},
			},
		},
	}

	response, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, requestBody, &state, &endpointConfig)

	assert.Equal(
		t,
		map[string]string{
			"Some-base-header-key": "Some base header value",
			"Another-header-key":   "Another header value",
		},
		response.Headers,
	)
}

func Test_ResolveEndpointResponse_Headers_WithBase_WithConditionalResponse_ConditionNotMatching(t *testing.T) {
	requestMock = &http.Request{URL: &url.URL{RawQuery: "foo=not_bar"}}
	osMockInstance := osMock{}
	endpointConfig := types.EndpointConfig{
		Route:    "foo/bar",
		Method:   "post",
		Response: []byte(`{"foo":"bar"}`),
		HeadersBase: map[string]string{
			"Some-base-header-key": "Some base header value",
		},
		Headers: map[string]string{
			"Some-header-key": "Some header value",
		},
		ResponseIf: []types.ResponseIf{
			types.ResponseIf{
				Response: []byte(`or chaining!`),
				Headers: map[string]string{
					"Another-header-key": "Another header value",
				},
				Condition: &types.Condition{
					Type:  types.ConditionType_QuerystringMatch,
					Key:   "foo",
					Value: "bar",
				},
			},
		},
	}

	response, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, requestBody, &state, &endpointConfig)

	assert.Equal(
		t,
		map[string]string{
			"Some-base-header-key": "Some base header value",
			"Some-header-key":      "Some header value",
		},
		response.Headers,
	)
}

func Test_ResolveEndpointResponse_FormMatch_Match(t *testing.T) {
	requestMock = httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(""),
	)
	requestMock.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	requestBody = []byte("foo=bar&foo2=bar2")
	osMockInstance := osMock{}
	endpointConfig := types.EndpointConfig{
		Route:    "foo/bar",
		Method:   "post",
		Response: []byte(`default response.`),
		ResponseIf: []types.ResponseIf{
			types.ResponseIf{
				Response: []byte(`this response shall not be returned.`),
				Condition: &types.Condition{
					Type:  types.ConditionType_QuerystringMatch,
					Key:   "foo",
					Value: "bar",
				},
			},
			types.ResponseIf{
				Response: []byte(`good!`),
				Condition: &types.Condition{
					Type: types.ConditionType_FormMatch,
					KeyValues: map[string]interface{}{
						"foo":  "bar",
						"foo2": "bar2",
					},
				},
			},
		},
	}

	response, _, _ := mock.ResolveEndpointResponse(osMockInstance.ReadFile, requestMock, requestBody, &state, &endpointConfig)

	assert.Equal(
		t,
		[]byte(`good!`),
		response.Body,
	)
}
