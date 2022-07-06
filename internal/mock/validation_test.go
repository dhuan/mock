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
			Method:   "GET",
			Response: []byte(`{"foo":"bar"}`),
		},
		types.EndpointConfig{
			Route:    "hello/world",
			Method:   "GET",
			Response: []byte(`{"foo":"bar"}`),
		},
		types.EndpointConfig{
			Route:    "foo/bar",
			Method:   "GET",
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
