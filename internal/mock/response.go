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
) ([]byte, types.Endpoint_content_type, error) {
	hasResponseIf := len(endpointConfig.ResponseIf) > 0
	matchingResponseIf := &types.ResponseIf{}

	if hasResponseIf {
		matchingResponseIfB, foundMatchingResponseIf := resolveResponseIf(request, endpointConfig)
		matchingResponseIf = matchingResponseIfB
		hasResponseIf = foundMatchingResponseIf
	}

	if hasResponseIf {
		return resolveEndpointResponseInternal(readFile, state, matchingResponseIf.Response)
	}

	return resolveEndpointResponseInternal(readFile, state, endpointConfig.Content)
}

func resolveResponseIf(request *http.Request, endpointConfig *types.EndpointConfig) (*types.ResponseIf, bool) {
	matchingResponseIfs := make([]int, 0)

	for responseIfKey, _ := range endpointConfig.ResponseIf {
		querystringConditions := endpointConfig.ResponseIf[responseIfKey].QuerystringMatches

		if len(querystringConditions) == 0 {
			continue
		}

		querystringMatch := querystringConditionsMatches(request, querystringConditions)

		if querystringMatch {
			matchingResponseIfs = append(matchingResponseIfs, responseIfKey)
		}
	}

	if len(matchingResponseIfs) == 0 {
		return &types.ResponseIf{}, false
	}

	return &endpointConfig.ResponseIf[matchingResponseIfs[0]], true
}

func querystringConditionsMatches(request *http.Request, querystringConditions []types.QuerystringMatches) bool {
	querystring := request.URL.Query()

	for i, _ := range querystringConditions {
		if !querystring.Has(querystringConditions[i].Key) {
			return false
		}

		if querystring.Get(querystringConditions[i].Key) != querystringConditions[i].Value {
			return false
		}
	}

	return true
}

func resolveEndpointResponseInternal(readFile ReadFileFunc, state *types.State, response types.EndpointConfigResponse) ([]byte, types.Endpoint_content_type, error) {
	endpointConfigContentType := resolveEndpointConfigContentType(response)

	if endpointConfigContentType == types.Endpoint_content_type_unknown {
		return []byte(""), endpointConfigContentType, nil
	}

	if endpointConfigContentType == types.Endpoint_content_type_plaintext {
		return []byte(utils.Unquote(string(response))), endpointConfigContentType, nil
	}

	if endpointConfigContentType == types.Endpoint_content_type_file {
		responseFile := fmt.Sprintf(
			"%s/%s",
			state.ConfigFolderPath,
			strings.Replace(string(response), "file:", "", -1),
		)
		fileContent, err := readFile(responseFile)
		if err != nil {
			return []byte(""), endpointConfigContentType, err
		}

		return fileContent, endpointConfigContentType, nil
	}

	if endpointConfigContentType == types.Endpoint_content_type_json {
		var jsonParsed interface{}
		err := json.Unmarshal(response, &jsonParsed)
		if err != nil {
			return []byte(""), endpointConfigContentType, err
		}

		jsonEncoded, err := json.Marshal(jsonParsed)
		if err != nil {
			return []byte(""), endpointConfigContentType, err
		}

		return jsonEncoded, endpointConfigContentType, nil
	}

	return []byte(""), types.Endpoint_content_type_unknown, nil
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
