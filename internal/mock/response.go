package mock

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
)

type ReadFileFunc = func(name string) ([]byte, error)

func ResolveEndpointResponse(
	readFile ReadFileFunc,
	state *types.State,
	endpointConfig *types.EndpointConfig,
) ([]byte, types.Endpoint_content_type, error) {
	endpointConfigContentType := resolveEndpointConfigContentType(endpointConfig)

	if endpointConfigContentType == types.Endpoint_content_type_unknown {
		return []byte(""), endpointConfigContentType, nil
	}

	if endpointConfigContentType == types.Endpoint_content_type_file {
		responseFile := fmt.Sprintf(
			"%s/%s",
			state.ConfigFolderPath,
			strings.Replace(string(endpointConfig.Content), "file:", "", -1),
		)
		fileContent, err := readFile(responseFile)
		if err != nil {
			return []byte(""), endpointConfigContentType, err
		}

		return fileContent, endpointConfigContentType, nil
	}

	if endpointConfigContentType == types.Endpoint_content_type_json {
		var jsonParsed interface{}
		err := json.Unmarshal(endpointConfig.Content, &jsonParsed)
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

func resolveEndpointConfigContentType(endpointConfig *types.EndpointConfig) types.Endpoint_content_type {
	if utils.BeginsWith(string(endpointConfig.Content), "file:") {
		return types.Endpoint_content_type_file
	}

	if utils.BeginsWith(string(endpointConfig.Content), "{") {
		return types.Endpoint_content_type_json
	}

	return types.Endpoint_content_type_unknown
}
