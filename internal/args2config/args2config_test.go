package args2config_test

import (
	"testing"

	"github.com/dhuan/mock/internal/args2config"
	"github.com/dhuan/mock/internal/types"
	"github.com/stretchr/testify/assert"
)

func Test_ParseEndpoints_WithEmptyArgs(t *testing.T) {
	assert.Equal(
		t,
		[]types.EndpointConfig{},
		args2config.ParseEndpoints([]string{}),
	)
}

func Test_ParseEndpoints_WithOneEndpoint(t *testing.T) {
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
		args2config.ParseEndpoints([]string{
			"--route",
			"foo/bar",
			"--method",
			"post",
			"--response",
			"Hello world!",
		}),
	)
}

func Test_ParseEndpoints_WithTwoEndpoints(t *testing.T) {
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
		args2config.ParseEndpoints([]string{
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

func Test_ParseEndpoints_WithStatusCode(t *testing.T) {
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
		args2config.ParseEndpoints([]string{
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

func Test_ParseEndpoints_WithResponseFile(t *testing.T) {
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
		args2config.ParseEndpoints([]string{
			"--route",
			"hello/world",
			"--response-file",
			"path/to/some/file.txt",
		}),
	)
}

func Test_ParseEndpoints_WithResponseFileServer(t *testing.T) {
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
			args2config.ParseEndpoints([]string{
				"--route",
				"public/*",
				responseVariation[0],
				responseVariation[1],
			}),
		)
	}
}

func Test_ParseEndpoints_WithResponseSh(t *testing.T) {
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
			args2config.ParseEndpoints([]string{
				"--route",
				"foo/bar",
				responseVariation[0],
				responseVariation[1],
			}),
		)
	}

}

func Test_ParseEndpoints_WithResponseExec(t *testing.T) {
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
			args2config.ParseEndpoints(responseVariation),
		)
	}
}

func Test_ParseEndpoints_WithHeaders(t *testing.T) {
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
		args2config.ParseEndpoints([]string{
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

func Test_ParseEndpoints_WithIrrelevantFlags(t *testing.T) {
	assert.Equal(
		t,
		[]types.EndpointConfig{},
		args2config.ParseEndpoints([]string{
			"--foo",
			"bar",
			"--hello",
			"world",
		}),
	)
}

func Test_ParseMiddlewares_WithEmptyArgs(t *testing.T) {
	assert.Equal(
		t,
		[]types.MiddlewareConfig{},
		args2config.ParseMiddlewares([]string{}),
	)
}

func Test_ParseMiddlewares_WithOneMiddleware_WithoutRouteMatch(t *testing.T) {
	assert.Equal(
		t,
		[]types.MiddlewareConfig{
			{
				Exec:       "sh /path/to/some/script.sh",
				Type:       types.MiddlewareType_BeforeResponse,
				RouteMatch: "*",
			},
		},
		args2config.ParseMiddlewares([]string{
			"--middleware-before-response",
			"sh /path/to/some/script.sh",
		}),
	)
}

func Test_ParseMiddlewares_WithOneMiddleware_WithRouteMatch(t *testing.T) {
	assert.Equal(
		t,
		[]types.MiddlewareConfig{
			{
				Exec:       "sh /path/to/some/script.sh",
				Type:       types.MiddlewareType_BeforeResponse,
				RouteMatch: "foobar",
			},
		},
		args2config.ParseMiddlewares([]string{
			"--middleware-before-response",
			"sh /path/to/some/script.sh",
			"--route-match",
			"foobar",
		}),
	)
}

func Test_ParseMiddlewares_WithMultipleMiddlewares(t *testing.T) {
	assert.Equal(
		t,
		[]types.MiddlewareConfig{
			{
				Exec:       "sh /path/to/some/script.sh",
				Type:       types.MiddlewareType_BeforeResponse,
				RouteMatch: "*",
			},
			{
				Exec:       "sh /path/to/another/script.sh",
				Type:       types.MiddlewareType_BeforeRequest,
				RouteMatch: "some_regex",
			},
		},
		args2config.ParseMiddlewares([]string{
			"--middleware-before-response",
			"sh /path/to/some/script.sh",
			"--middleware-before-request",
			"sh /path/to/another/script.sh",
			"--route-match",
			"some_regex",
		}),
	)
}

func Test_ParseMiddlewares_WithIrrelevantFlags(t *testing.T) {
	assert.Equal(
		t,
		[]types.MiddlewareConfig{},
		args2config.ParseMiddlewares([]string{
			"--foo",
			"bar",
			"--hello",
			"world",
		}),
	)
}
