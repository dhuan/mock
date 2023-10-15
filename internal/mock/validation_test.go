package mock_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dhuan/mock/internal/mock"
	"github.com/dhuan/mock/internal/types"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

type readFileMock struct {
	testifymock.Mock
}

func (rfm *readFileMock) ReadFile(name string) ([]byte, error) {
	args := rfm.Called(name)

	return args.Get(0).([]byte), args.Get(1).(error)
}

var readFileMockInstance = readFileMock{}

var configDirPath = "path/to/somewhere"

func Test_ValidateEndpointConfigs_Duplicates(t *testing.T) {
	endpointConfigs := []types.EndpointConfig{
		{
			Route:    "foo/bar",
			Method:   "get",
			Response: []byte(`{"foo":"bar"}`),
		},
		{
			Route:    "hello/world",
			Method:   "get",
			Response: []byte(`{"foo":"bar"}`),
		},
		{
			Route:    "foo/bar",
			Method:   "get",
			Response: []byte(`{"hello":"world"}`),
		},
	}

	validationErrors, _ := mock.ValidateEndpointConfigs(endpointConfigs, readFileMockInstance.ReadFile, configDirPath)

	assert.Equal(
		t,
		[]mock.EndpointConfigError{
			{
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
		{
			Route:    "foo/bar",
			Method:   "get",
			Response: []byte(`{"foo":"bar"}`),
		},
		{
			Route:    "hello/world",
			Method:   "foobar",
			Response: []byte(`{"foo":"bar"}`),
		},
	}

	validationErrors, _ := mock.ValidateEndpointConfigs(endpointConfigs, readFileMockInstance.ReadFile, configDirPath)

	assert.Equal(
		t,
		[]mock.EndpointConfigError{
			{
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
		{
			Route:    "foo/bar",
			Method:   "get",
			Response: []byte(`{"foo":"bar"}`),
		},
		{
			Route:    "hello/world?foo=bar",
			Method:   "get",
			Response: []byte(`{"foo":"bar"}`),
		},
	}

	validationErrors, _ := mock.ValidateEndpointConfigs(endpointConfigs, readFileMockInstance.ReadFile, configDirPath)

	assert.Equal(
		t,
		[]mock.EndpointConfigError{
			{
				Code:          mock.EndpointConfigErrorCode_RouteWithQuerystring,
				EndpointIndex: 1,
				Metadata:      map[string]string{},
			},
		},
		validationErrors,
	)
}

func Test_ValidateEndpointConfigs_FailingToReadFile_WithTxtFile(t *testing.T) {
	testFailingToReadFile(t, "file:some_file.txt", "some_file.txt")
}

func Test_ValidateEndpointConfigs_FailingToReadFile_WithShellScriptFile(t *testing.T) {
	testFailingToReadFile(t, "sh:some_file.txt", "some_file.txt")
}

func testFailingToReadFile(t *testing.T, responseContent, fileName string) {
	endpointConfigs := []types.EndpointConfig{
		{
			Route:    "foo/bar",
			Method:   "get",
			Response: []byte(responseContent),
		},
	}

	readFileMockInstance.On(
		"ReadFile",
		fmt.Sprintf("path/to/somewhere/%s", fileName),
	).Return([]byte(""), errors.New("Some error."))

	validationErrors, _ := mock.ValidateEndpointConfigs(endpointConfigs, readFileMockInstance.ReadFile, configDirPath)

	assert.Equal(
		t,
		[]mock.EndpointConfigError{
			{
				Code:          mock.EndpointConfigErrorCode_FileUnreadable,
				EndpointIndex: 0,
				Metadata: map[string]string{
					"file_path": fileName,
				},
			},
		},
		validationErrors,
	)
}
