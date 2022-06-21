package mock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
)

type ReadFileFunc = func(name string) ([]byte, error)

func ResolveEndpointResponse(
	readFile ReadFileFunc,
	request *http.Request,
	state *types.State,
	endpointConfig *types.EndpointConfig,
) ([]byte, types.Endpoint_content_type, int, error) {
	hasResponseIf := len(endpointConfig.ResponseIf) > 0
	matchingResponseIf := &types.ResponseIf{}

	if hasResponseIf {
		matchingResponseIfB, foundMatchingResponseIf := resolveResponseIf(request, endpointConfig)
		matchingResponseIf = matchingResponseIfB
		hasResponseIf = foundMatchingResponseIf
	}

	if hasResponseIf {
		return resolveEndpointResponseInternal(
			readFile,
			state,
			matchingResponseIf.Response,
			resolveResponseStatusCode(matchingResponseIf.ResponseStatusCode),
		)
	}

	return resolveEndpointResponseInternal(
		readFile,
		state,
		endpointConfig.Response,
		resolveResponseStatusCode(endpointConfig.ResponseStatusCode),
	)
}

func resolveResponseStatusCode(statusCode int) int {
	if statusCode < 1 {
		return 200
	}

	return statusCode
}

func resolveResponseIf(request *http.Request, endpointConfig *types.EndpointConfig) (*types.ResponseIf, bool) {
	matchingResponseIfs := make([]int, 0)

	for responseIfKey, _ := range endpointConfig.ResponseIf {
		responseIf := endpointConfig.ResponseIf[responseIfKey]
		matches := resolveSingleResponseIf(request, responseIf.Condition)

		if matches {
			matchingResponseIfs = append(matchingResponseIfs, responseIfKey)
		}
	}

	if len(matchingResponseIfs) == 0 {
		return &types.ResponseIf{}, false
	}

	return &endpointConfig.ResponseIf[matchingResponseIfs[0]], true
}

func resolveSingleResponseIf(request *http.Request, condition *types.Condition) bool {
	conditionFunction := resolveConditionFunction(condition)
	result := conditionFunction(request, condition.Key, condition.Value)
	hasAnd := condition.And != nil
	hasOr := condition.Or != nil

	if result && !hasAnd {
		return true
	}

	if result && hasAnd {
		return resolveSingleResponseIf(request, condition.And)
	}

	if !result && hasOr {
		return resolveSingleResponseIf(request, condition.Or)
	}

	if !result && !hasOr {
		return false
	}

	return false
}

func resolveConditionFunction(condition *types.Condition) func(request *http.Request, key, value string) bool {
	if condition.Type == types.ConditionType_QuerystringMatch {
		return conditionQuerystringMatch
	}

	panic("Failed to resolve condition func!")
}

func conditionQuerystringMatch(request *http.Request, key, value string) bool {
	query := request.URL.Query()

	if !query.Has(key) {
		return false
	}

	return value == query.Get(key)
}

func resolveEndpointResponseInternal(
	readFile ReadFileFunc,
	state *types.State,
	response types.EndpointConfigResponse,
	responseStatusCode int,
) ([]byte, types.Endpoint_content_type, int, error) {
	endpointConfigContentType := resolveEndpointConfigContentType(response)

	if endpointConfigContentType == types.Endpoint_content_type_unknown {
		return []byte(""), endpointConfigContentType, responseStatusCode, nil
	}

	if endpointConfigContentType == types.Endpoint_content_type_plaintext {
		return []byte(utils.Unquote(string(response))), endpointConfigContentType, responseStatusCode, nil
	}

	if endpointConfigContentType == types.Endpoint_content_type_file {
		responseFile := fmt.Sprintf(
			"%s/%s",
			state.ConfigFolderPath,
			strings.Replace(string(response), "file:", "", -1),
		)
		fileContent, err := readFile(responseFile)
		if err != nil {
			return []byte(""), endpointConfigContentType, responseStatusCode, err
		}

		return fileContent, endpointConfigContentType, responseStatusCode, nil
	}

	if endpointConfigContentType == types.Endpoint_content_type_json {
		var jsonParsed interface{}
		err := json.Unmarshal(response, &jsonParsed)
		if err != nil {
			return []byte(""), endpointConfigContentType, responseStatusCode, err
		}

		jsonEncoded, err := json.Marshal(jsonParsed)
		if err != nil {
			return []byte(""), endpointConfigContentType, responseStatusCode, err
		}

		return jsonEncoded, endpointConfigContentType, responseStatusCode, nil
	}

	return []byte(""), types.Endpoint_content_type_unknown, responseStatusCode, nil
}

func resolveEndpointConfigContentType(response types.EndpointConfigResponse) types.Endpoint_content_type {
	if utils.BeginsWith(string(response), "file:") {
		return types.Endpoint_content_type_file
	}

	if utils.BeginsWith(string(response), "{") {
		return types.Endpoint_content_type_json
	}

	return types.Endpoint_content_type_plaintext
}
