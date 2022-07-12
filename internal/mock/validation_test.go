package mock_test

import (
	"testing"

	"github.com/dhuan/mock/internal/mock"
	"github.com/dhuan/mock/internal/types"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateEndpointConfigs_Duplicates(t *testing.T) {
	endpointConfigs := []types.EndpointConfig{
		types.EndpointConfig{
			Route:    "foo/bar",
			Method:   "get",
			Response: []byte(`{"foo":"bar"}`),
		},
		types.EndpointConfig{
			Route:    "hello/world",
			Method:   "get",
			Response: []byte(`{"foo":"bar"}`),
		},
		types.EndpointConfig{
			Route:    "foo/bar",
			Method:   "get",
			Response: []byte(`{"hello":"world"}`),
		},
	}

	validationErrors, _ := mock.ValidateEndpointConfigs(endpointConfigs)

	assert.Equal(
		t,
		[]mock.EndpointConfigError{
			mock.EndpointConfigError{
				Code:          mock.EndpointConfigErrorCode_EndpointDuplicate,
				EndpointIndex: 0,
				Metadata: map[string]string{
					"duplicate_index": "2",
				},
			},
		},
		validationErrors,
	)
}

func Test_ValidateEndpointConfigs_InvalidMethod(t *testing.T) {
	endpointConfigs := []types.EndpointConfig{
		types.EndpointConfig{
			Route:    "foo/bar",
			Method:   "get",
			Response: []byte(`{"foo":"bar"}`),
		},
		types.EndpointConfig{
			Route:    "hello/world",
			Method:   "foobar",
			Response: []byte(`{"foo":"bar"}`),
		},
	}

	validationErrors, _ := mock.ValidateEndpointConfigs(endpointConfigs)

	assert.Equal(
		t,
		[]mock.EndpointConfigError{
			mock.EndpointConfigError{
				Code:          mock.EndpointConfigErrorCode_InvalidMethod,
				EndpointIndex: 1,
				Metadata: map[string]string{
					"method": "foobar",
				},
			},
		},
		validationErrors,
	)
}

func Test_ValidateEndpointConfigs_WithQuerystring(t *testing.T) {
	endpointConfigs := []types.EndpointConfig{
		types.EndpointConfig{
			Route:    "foo/bar",
			Method:   "get",
			Response: []byte(`{"foo":"bar"}`),
		},
		types.EndpointConfig{
			Route:    "hello/world?foo=bar",
			Method:   "get",
			Response: []byte(`{"foo":"bar"}`),
		},
	}

	validationErrors, _ := mock.ValidateEndpointConfigs(endpointConfigs)

	assert.Equal(
		t,
		[]mock.EndpointConfigError{
			mock.EndpointConfigError{
				Code:          mock.EndpointConfigErrorCode_RouteWithQuerystring,
				EndpointIndex: 1,
				Metadata:      map[string]string{},
			},
		},
		validationErrors,
	)
}
