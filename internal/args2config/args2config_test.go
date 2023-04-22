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

func Test_WithResponseFile(t *testing.T) {
	assert.Equal(
		t,
		[]types.EndpointConfig{
			{
				Route:              "hello/world",
				Method:             "",
				Response:           []byte("file:path/to/some/file.txt"),
				ResponseStatusCode: 0,
				Headers:            nil,
				ResponseIf:         nil,
				HeadersBase:        nil,
			},
		},
		args2config.Parse([]string{
			"--route",
			"hello/world",
			"--response-file",
			"path/to/some/file.txt",
		}),
	)
}

func Test_WithResponseFileServer(t *testing.T) {
	responseVariations := [][]string{
		{"--response-file-server", "path/to/my/files"},
		{"--file-server", "path/to/my/files"},
	}

	for _, responseVariation := range responseVariations {
		assert.Equal(
			t,
			[]types.EndpointConfig{
				{
					Route:              "public/*",
					Method:             "",
					Response:           []byte("fs:path/to/my/files"),
					ResponseStatusCode: 0,
					Headers:            nil,
					ResponseIf:         nil,
					HeadersBase:        nil,
				},
			},
			args2config.Parse([]string{
				"--route",
				"public/*",
				responseVariation[0],
				responseVariation[1],
			}),
		)
	}
}

func Test_WithResponseSh(t *testing.T) {
	responseVariations := [][]string{
		{"--response-sh", "path/to/some/script.sh"},
		{"--shell-script", "path/to/some/script.sh"},
	}

	for _, responseVariation := range responseVariations {
		assert.Equal(
			t,
			[]types.EndpointConfig{
				{
					Route:              "foo/bar",
					Method:             "",
					Response:           []byte("sh:path/to/some/script.sh"),
					ResponseStatusCode: 0,
					Headers:            nil,
					ResponseIf:         nil,
					HeadersBase:        nil,
				},
			},
			args2config.Parse([]string{
				"--route",
				"foo/bar",
				responseVariation[0],
				responseVariation[1],
			}),
		)
	}

}

func Test_WithResponseExec(t *testing.T) {
	responseVariations := [][]string{
		{
			"--route",
			"foo/bar",
			"--exec",
			"some command | some other command",
		},
		{
			"--route",
			"foo/bar",
			"--response-exec",
			"some command | some other command",
		},
	}

	for _, responseVariation := range responseVariations {
		assert.Equal(
			t,
			[]types.EndpointConfig{
				{
					Route:              "foo/bar",
					Method:             "",
					Response:           []byte("exec:some command | some other command"),
					ResponseStatusCode: 0,
					Headers:            nil,
					ResponseIf:         nil,
					HeadersBase:        nil,
				},
			},
			args2config.Parse(responseVariation),
		)
	}
}

func Test_WithHeaders(t *testing.T) {
	assert.Equal(
		t,
		[]types.EndpointConfig{
			{
				Route:              "endpoint/one",
				Method:             "",
				Response:           []byte("Endpoint one's response."),
				ResponseStatusCode: 0,
				Headers: map[string]string{
					"Header-One": "This is the 1st header.",
					"Header-Two": "This is the 2nd header.",
				},
				ResponseIf:  nil,
				HeadersBase: nil,
			},
			{
				Route:              "endpoint/two",
				Method:             "",
				Response:           []byte("Endpoint two's response."),
				ResponseStatusCode: 0,
				Headers: map[string]string{
					"Header-Three": "This is the 3rd header.",
				},
				ResponseIf:  nil,
				HeadersBase: nil,
			},
		},
		args2config.Parse([]string{
			"--route",
			"endpoint/one",
			"--header",
			"Header-One: This is the 1st header.",
			"--header",
			"Header-Two: This is the 2nd header.",
			"--response",
			"Endpoint one's response.",
			"--route",
			"endpoint/two",
			"--response",
			"Endpoint two's response.",
			"--header",
			"Header-Three: This is the 3rd header.",
		}),
	)
}

func Test_WithIrrelevantFlags(t *testing.T) {
	assert.Equal(
		t,
		[]types.EndpointConfig{},
		args2config.Parse([]string{
			"--foo",
			"bar",
			"--hello",
			"world",
		}),
	)
}
