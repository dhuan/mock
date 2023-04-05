package args2config_test

import (
	"testing"

	"github.com/dhuan/mock/internal/args2config"
	"github.com/dhuan/mock/internal/types"
	"github.com/stretchr/testify/assert"
)

func Test_WithEmptyArgs(t *testing.T) {
	assert.Equal(
		t,
		[]types.EndpointConfig{},
		args2config.Parse([]string{}),
	)
}

func Test_WithOneEndpoint(t *testing.T) {
	assert.Equal(
		t,
		[]types.EndpointConfig{
			{
				Route:              "foo/bar",
				Method:             "post",
				Response:           []byte("Hello world!"),
				ResponseStatusCode: 0,
				Headers:            nil,
				ResponseIf:         nil,
				HeadersBase:        nil,
			},
		},
		args2config.Parse([]string{
			"--route",
			"foo/bar",
			"--method",
			"post",
			"--response",
			"Hello world!",
		}),
	)
}

func Test_WithTwoEndpoints(t *testing.T) {
	assert.Equal(
		t,
		[]types.EndpointConfig{
			{
				Route:              "endpoint/one",
				Method:             "get",
				Response:           []byte("Endpoint one's response."),
				ResponseStatusCode: 0,
				Headers:            nil,
				ResponseIf:         nil,
				HeadersBase:        nil,
			},
			{
				Route:              "endpoint/two",
				Method:             "post",
				Response:           []byte("Endpoint two's response."),
				ResponseStatusCode: 0,
				Headers:            nil,
				ResponseIf:         nil,
				HeadersBase:        nil,
			},
		},
		args2config.Parse([]string{
			"--route",
			"endpoint/one",
			"--method",
			"get",
			"--response",
			"Endpoint one's response.",
			"--route",
			"endpoint/two",
			"--method",
			"post",
			"--response",
			"Endpoint two's response.",
		}),
	)
}

func Test_WithStatusCode(t *testing.T) {
	assert.Equal(
		t,
		[]types.EndpointConfig{
			{
				Route:              "endpoint/one",
				Method:             "",
				Response:           []byte("Endpoint one's response."),
				ResponseStatusCode: 0,
				Headers:            nil,
				ResponseIf:         nil,
				HeadersBase:        nil,
			},
			{
				Route:              "endpoint/two",
				Method:             "",
				Response:           []byte("Endpoint two's response."),
				ResponseStatusCode: 201,
				Headers:            nil,
				ResponseIf:         nil,
				HeadersBase:        nil,
			},
		},
		args2config.Parse([]string{
			"--route",
			"endpoint/one",
			"--response",
			"Endpoint one's response.",
			"--route",
			"endpoint/two",
			"--status-code",
			"201",
			"--response",
			"Endpoint two's response.",
		}),
	)
}
