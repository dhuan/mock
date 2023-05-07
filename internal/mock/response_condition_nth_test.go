package mock_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dhuan/mock/internal/mock"
	"github.com/dhuan/mock/internal/types"
	. "github.com/dhuan/mock/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func newRequestRecord(route, method string) types.RequestRecord {
	return types.RequestRecord{
		Route:  route,
		Method: method,
	}
}

var endpoint_config_with_nth_conditions = types.EndpointConfig{
	Route:    "foo/bar",
	Method:   "get",
	Response: []byte(`default response.`),
	ResponseIf: []types.ResponseIf{
		{
			Response: []byte(`this is the second response.`),
			Condition: &Condition{
				Type:  ConditionType_Nth,
				Value: "2",
			},
		},
		{
			Response: []byte(`this is the third response.`),
			Condition: &Condition{
				Type:  ConditionType_Nth,
				Value: "3",
			},
		},
	},
}

func Test_ResolveEndpointResponse_Condition_Nth_FirstRequest(t *testing.T) {
	osMockInstance := osMock{}
	execMockInstance := execMock{}

	response, _, _ := mock.ResolveEndpointResponse(
		osMockInstance.ReadFile,
		execMockInstance.Exec,
		newRequest("foo/bar", http.MethodGet),
		requestBody,
		&state,
		&endpoint_config_with_nth_conditions,
		map[string]string{},
		map[string]string{},
		[]types.RequestRecord{
			newRequestRecord("irrelevant_request", "get"),
			newRequestRecord("irrelevant_request", "post"),
		},
	)

	assert.Equal(t, []byte(`default response.`), response.Body)
}

func Test_ResolveEndpointResponse_Condition_Nth_SecondRequest(t *testing.T) {
	osMockInstance := osMock{}
	execMockInstance := execMock{}

	response, _, _ := mock.ResolveEndpointResponse(
		osMockInstance.ReadFile,
		execMockInstance.Exec,
		newRequest("foo/bar", http.MethodGet),
		requestBody,
		&state,
		&endpoint_config_with_nth_conditions,
		map[string]string{},
		map[string]string{},
		[]types.RequestRecord{
			newRequestRecord("irrelevant_request", "get"),
			newRequestRecord("foo/bar", "get"),
			newRequestRecord("irrelevant_request", "get"),
			newRequestRecord("foo/bar", "post"),
			newRequestRecord("foo/bar", "put"),
		},
	)

	assert.Equal(t, []byte(`this is the second response.`), response.Body)
}

func Test_ResolveEndpointResponse_Condition_Nth_ThirdRequest(t *testing.T) {
	osMockInstance := osMock{}
	execMockInstance := execMock{}

	response, _, _ := mock.ResolveEndpointResponse(
		osMockInstance.ReadFile,
		execMockInstance.Exec,
		newRequest("foo/bar", http.MethodGet),
		requestBody,
		&state,
		&endpoint_config_with_nth_conditions,
		map[string]string{},
		map[string]string{},
		[]types.RequestRecord{
			newRequestRecord("irrelevant_request", "get"),
			newRequestRecord("foo/bar", "get"),
			newRequestRecord("irrelevant_request", "get"),
			newRequestRecord("foo/bar", "get"),
		},
	)

	assert.Equal(t, []byte(`this is the third response.`), response.Body)
}

func Test_ResolveEndpointResponse_Condition_Nth_FourthRequest(t *testing.T) {
	osMockInstance := osMock{}
	execMockInstance := execMock{}

	response, _, _ := mock.ResolveEndpointResponse(
		osMockInstance.ReadFile,
		execMockInstance.Exec,
		newRequest("foo/bar", http.MethodGet),
		requestBody,
		&state,
		&endpoint_config_with_nth_conditions,
		map[string]string{},
		map[string]string{},
		[]types.RequestRecord{
			newRequestRecord("irrelevant_request", "get"),
			newRequestRecord("foo/bar", "get"),
			newRequestRecord("irrelevant_request", "get"),
			newRequestRecord("foo/bar", "get"),
			newRequestRecord("foo/bar", "get"),
			newRequestRecord("foo/bar", "get"),
			newRequestRecord("foo/bar", "get"),
		},
	)

	assert.Equal(t, []byte(`default response.`), response.Body)
}

func newRequest(route, method string) *http.Request {
	requestMock = httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	requestMock.RequestURI = "foo/bar"

	return requestMock
}
