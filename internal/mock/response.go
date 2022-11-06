package mock

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dhuan/mock/internal/record"
	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
	. "github.com/dhuan/mock/pkg/mock"
)

type ReadFileFunc = func(name string) ([]byte, error)

type ExecFunc = func(command string) (*ExecResult, error)

type Response struct {
	Body                []byte
	EndpointContentType types.Endpoint_content_type
	StatusCode          int
	Headers             map[string]string
}

type ExecResult struct {
	Output []byte
}

func ResolveEndpointResponse(
	readFile ReadFileFunc,
    exec ExecFunc,
	request *http.Request,
	requestBody []byte,
	state *types.State,
	endpointConfig *types.EndpointConfig,
) (*Response, error, map[string]string) {
	hasResponseIf := len(endpointConfig.ResponseIf) > 0
	matchingResponseIf := &types.ResponseIf{}
	requestRecord, err := record.BuildRequestRecord(request, requestBody)
	if err != nil {
		return &Response{[]byte(""), types.Endpoint_content_type_unknown, 0, nil}, err, nil
	}

	if hasResponseIf {
		matchingResponseIfB, foundMatchingResponseIf := resolveResponseIf(requestRecord, endpointConfig)
		matchingResponseIf = matchingResponseIfB
		hasResponseIf = foundMatchingResponseIf
	}

	if hasResponseIf {
		return resolveEndpointResponseInternal(
			readFile,
            exec,
			state,
			matchingResponseIf.Response,
			resolveResponseStatusCode(matchingResponseIf.ResponseStatusCode),
			endpointConfig,
			matchingResponseIf,
			hasResponseIf,
		)
	}

	return resolveEndpointResponseInternal(
		readFile,
        exec,
		state,
		endpointConfig.Response,
		resolveResponseStatusCode(endpointConfig.ResponseStatusCode),
		endpointConfig,
		matchingResponseIf,
		hasResponseIf,
	)
}

func resolveResponseStatusCode(statusCode int) int {
	if statusCode < 1 {
		return 200
	}

	return statusCode
}

func resolveResponseIf(requestRecord *types.RequestRecord, endpointConfig *types.EndpointConfig) (*types.ResponseIf, bool) {
	matchingResponseIfs := make([]int, 0)

	for responseIfKey, _ := range endpointConfig.ResponseIf {
		responseIf := endpointConfig.ResponseIf[responseIfKey]
		matches := resolveSingleResponseIf(requestRecord, responseIf.Condition)

		if matches {
			matchingResponseIfs = append(matchingResponseIfs, responseIfKey)
		}
	}

	if len(matchingResponseIfs) == 0 {
		return &types.ResponseIf{}, false
	}

	return &endpointConfig.ResponseIf[matchingResponseIfs[0]], true
}

func resolveSingleResponseIf(requestRecord *types.RequestRecord, condition *Condition) bool {
	conditionFunction := resolveAssertTypeFunc(condition.Type)
	validationErrors, err := conditionFunction(requestRecord, condition)
	if err != nil {
		panic(err)
	}
	result := len(validationErrors) == 0

	hasAnd := condition.And != nil
	hasOr := condition.Or != nil

	if result && !hasAnd {
		return true
	}

	if result && hasAnd {
		return resolveSingleResponseIf(requestRecord, condition.And)
	}

	if !result && hasOr {
		return resolveSingleResponseIf(requestRecord, condition.Or)
	}

	if !result && !hasOr {
		return false
	}

	return false
}

func resolveEndpointResponseInternal(
	readFile ReadFileFunc,
    exec ExecFunc,
	state *types.State,
	response types.EndpointConfigResponse,
	responseStatusCode int,
	endpointConfig *types.EndpointConfig,
	responseIf *types.ResponseIf,
	hasResponseIf bool,
) (*Response, error, map[string]string) {
	errorMetadata := make(map[string]string)
	endpointConfigContentType := resolveEndpointConfigContentType(response)
	headers := make(map[string]string)
	utils.JoinMap[string, string](headers, endpointConfig.Headers)
	utils.JoinMap[string, string](headers, endpointConfig.HeadersBase)

	if hasResponseIf {
		headers = make(map[string]string)
		utils.JoinMap[string, string](headers, endpointConfig.HeadersBase)
		utils.JoinMap[string, string](headers, responseIf.Headers)
	}

	if endpointConfigContentType == types.Endpoint_content_type_unknown {
		return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, nil, errorMetadata
	}

	if endpointConfigContentType == types.Endpoint_content_type_plaintext {
		return &Response{[]byte(utils.Unquote(string(response))), endpointConfigContentType, responseStatusCode, headers}, nil, errorMetadata
	}

	if endpointConfigContentType == types.Endpoint_content_type_file {
		responseFile := fmt.Sprintf(
			"%s/%s",
			state.ConfigFolderPath,
			strings.Replace(string(response), "file:", "", -1),
		)
		fileContent, err := readFile(responseFile)
		if errors.Is(err, ErrResponseFileDoesNotExist) {
			errorMetadata["file"] = responseFile
		}
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		return &Response{fileContent, endpointConfigContentType, responseStatusCode, headers}, nil, errorMetadata
	}

	if endpointConfigContentType == types.Endpoint_content_type_shell {
		scriptFilePath := fmt.Sprintf(
			"%s/%s",
			state.ConfigFolderPath,
			strings.Replace(string(response), "sh:", "", -1),
		)

        execResult, err := exec(fmt.Sprintf(
            "sh %s",
            scriptFilePath,
        ))
        if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
        }

		return &Response{execResult.Output, endpointConfigContentType, responseStatusCode, headers}, nil, errorMetadata
    }

	if endpointConfigContentType == types.Endpoint_content_type_json {
		var jsonParsed interface{}
		err := json.Unmarshal(response, &jsonParsed)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		jsonEncoded, err := json.Marshal(jsonParsed)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		return &Response{jsonEncoded, endpointConfigContentType, responseStatusCode, headers}, nil, errorMetadata
	}

	return &Response{[]byte(""), types.Endpoint_content_type_unknown, responseStatusCode, headers}, nil, errorMetadata
}

func resolveEndpointConfigContentType(response types.EndpointConfigResponse) types.Endpoint_content_type {
	if utils.BeginsWith(string(response), "file:") {
		return types.Endpoint_content_type_file
	}

	if utils.BeginsWith(string(response), "sh:") {
		return types.Endpoint_content_type_shell
	}

	if utils.BeginsWith(string(response), "{") {
		return types.Endpoint_content_type_json
	}

	return types.Endpoint_content_type_plaintext
}
