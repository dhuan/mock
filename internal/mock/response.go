package mock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
)

type ReadFileFunc = func(name string) ([]byte, error)

type Response struct {
	Body                []byte
	EndpointContentType types.Endpoint_content_type
	StatusCode          int
	Headers             map[string]string
}

func ResolveEndpointResponse(
	readFile ReadFileFunc,
	request *http.Request,
	state *types.State,
	endpointConfig *types.EndpointConfig,
) (*Response, error) {
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
			endpointConfig,
			matchingResponseIf,
			hasResponseIf,
		)
	}

	return resolveEndpointResponseInternal(
		readFile,
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
	result := conditionFunction(request, condition)
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

func resolveConditionFunction(condition *types.Condition) func(request *http.Request, condition *types.Condition) bool {
	if condition.Type == types.ConditionType_QuerystringMatch {
		return conditionQuerystringMatch
	}

	panic("Failed to resolve condition func!")
}

func conditionQuerystringMatch(request *http.Request, condition *types.Condition) bool {
	query := request.URL.Query()
	isSingle := condition.Key != "" && condition.Value != ""
	isMultiple := len(condition.KeyValues) > 0

	if isSingle {
		if !query.Has(condition.Key) {
			return false
		}

		return condition.Value == query.Get(condition.Key)
	}

	if isMultiple {
		return conditionQuerystringMatchWithMany(request, condition, query)
	}

	panic("Failed to resolve query string match!")
}

func conditionQuerystringMatchWithMany(request *http.Request, condition *types.Condition, query url.Values) bool {
	for i, _ := range condition.KeyValues {
		value := fmt.Sprint(condition.KeyValues[i])

		if !query.Has(i) {
			return false
		}

		if value != query.Get(i) {
			return false
		}
	}

	return true
}

func resolveEndpointResponseInternal(
	readFile ReadFileFunc,
	state *types.State,
	response types.EndpointConfigResponse,
	responseStatusCode int,
	endpointConfig *types.EndpointConfig,
	responseIf *types.ResponseIf,
	hasResponseIf bool,
) (*Response, error) {
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
		return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, nil
	}

	if endpointConfigContentType == types.Endpoint_content_type_plaintext {
		return &Response{[]byte(utils.Unquote(string(response))), endpointConfigContentType, responseStatusCode, headers}, nil
	}

	if endpointConfigContentType == types.Endpoint_content_type_file {
		responseFile := fmt.Sprintf(
			"%s/%s",
			state.ConfigFolderPath,
			strings.Replace(string(response), "file:", "", -1),
		)
		fileContent, err := readFile(responseFile)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err
		}

		return &Response{fileContent, endpointConfigContentType, responseStatusCode, headers}, nil
	}

	if endpointConfigContentType == types.Endpoint_content_type_json {
		var jsonParsed interface{}
		err := json.Unmarshal(response, &jsonParsed)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err
		}

		jsonEncoded, err := json.Marshal(jsonParsed)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err
		}

		return &Response{jsonEncoded, endpointConfigContentType, responseStatusCode, headers}, nil
	}

	return &Response{[]byte(""), types.Endpoint_content_type_unknown, responseStatusCode, headers}, nil
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
