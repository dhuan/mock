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

func Test_ResolveEndpointResponse_Condition_Nth_1(t *testing.T) {
	requestMock = httptest.NewRequest(
		http.MethodGet,
		"/",
		strings.NewReader(""),
	)
	osMockInstance := osMock{}
	execMockInstance := execMock{}
	endpointConfig := types.EndpointConfig{
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
		},
	}

	response, _, _ := mock.ResolveEndpointResponse(
		osMockInstance.ReadFile,
		execMockInstance.Exec,
		requestMock,
		requestBody,
		&state,
		&endpointConfig,
		map[string]string{},
		map[string]string{},
		[]types.RequestRecord{},
	)

	assert.Equal(
		t,
		[]byte(`default response.`),
		response.Body,
	)
}
