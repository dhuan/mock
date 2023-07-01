package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dhuan/mock/internal/mock"
	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
)

func RunMiddleware(
	exec mock.ExecFunc,
	readFile types.ReadFileFunc,
	configPath string,
	middlewareConfigs []types.MiddlewareConfig,
	responseBody []byte,
	responseHeaders map[string]string,
	responseStatusCode int,
	request *http.Request,
	endpointParams map[string]string,
	vars map[string]string,
) ([]byte, map[string]string, int, error) {
	headersParsed := make(map[string]string)
	resultResponseBody := make([]byte, len(responseBody))
	copy(resultResponseBody, responseBody)

	resultResponseHeaders := []byte(utils.ToHeadersText(toHttpHeaders(responseHeaders)) + "\n")

	resultResponseStatusCode := []byte(fmt.Sprintf("%d", responseStatusCode))

	if len(middlewareConfigs) == 0 {
		return responseBody, headersParsed, responseStatusCode, nil
	}

	for i := range middlewareConfigs {
		responseBodyFile, err := utils.CreateTempFile(resultResponseBody)
		if err != nil {
			return responseBody, headersParsed, responseStatusCode, err
		}

		responseHeadersFile, err := utils.CreateTempFile(resultResponseHeaders)
		if err != nil {
			return responseBody, headersParsed, responseStatusCode, err
		}

		responseStatusCodeFile, err := utils.CreateTempFile(resultResponseStatusCode)
		if err != nil {
			return responseBody, headersParsed, responseStatusCode, err
		}

		envVars := map[string]string{
			"MOCK_RESPONSE_BODY":        responseBodyFile,
			"MOCK_RESPONSE_HEADERS":     responseHeadersFile,
			"MOCK_RESPONSE_STATUS_CODE": responseStatusCodeFile,
		}

		for key := range endpointParams {
			endpointParamKey := fmt.Sprintf("MOCK_ROUTE_PARAM_%s", strings.ToUpper(key))

			envVars[endpointParamKey] = endpointParams[key]
		}

		for key := range vars {
			envVars[key] = vars[key]
		}

		_, err = exec(
			middlewareConfigs[i].Exec,
			&mock.ExecOptions{
				Env:        envVars,
				WorkingDir: configPath,
			},
		)
		if err != nil {
			return responseBody, headersParsed, responseStatusCode, err
		}

		resultResponseBody, err = readFile(responseBodyFile)
		if err != nil {
			return responseBody, headersParsed, responseStatusCode, err
		}

		resultResponseHeaders, err = readFile(responseHeadersFile)
		if err != nil {
			return responseBody, headersParsed, responseStatusCode, err
		}

		resultResponseStatusCode, err = readFile(responseStatusCodeFile)
		if err != nil {
			return responseBody, headersParsed, responseStatusCode, err
		}
	}

	headersParsed = utils.ExtractHeadersFromText(resultResponseHeaders)

	responseStatusCodeParsed := bytesToInt(resultResponseStatusCode, responseStatusCode)

	return resultResponseBody, headersParsed, responseStatusCodeParsed, nil
}

func bytesToInt(data []byte, fallback int) int {
	statusCodeParsed, err := strconv.Atoi(string(data))
	if err != nil {
		return fallback
	}

	return statusCodeParsed
}

func toHttpHeaders(m map[string]string) http.Header {
	result := make(http.Header)

	for key := range m {
		result[key] = []string{m[key]}
	}

	return result
}

func GetMiddlewareForRequest(middlewareConfigs []types.MiddlewareConfig, r *http.Request) []types.MiddlewareConfig {
	middlewares := make([]types.MiddlewareConfig, 0)

	for i := range middlewareConfigs {
		if routeMatch(r, &middlewareConfigs[i]) {
			middlewares = append(middlewares, middlewareConfigs[i])
		}
	}

	return middlewares
}

func routeMatch(r *http.Request, middlewareConfig *types.MiddlewareConfig) bool {
	if middlewareConfig.RouteMatch == "*" || middlewareConfig.RouteMatch == "" {
		return true
	}

	requestRoute := utils.ReplaceRegex(r.URL.Path, []string{"^/"}, "")

	return utils.RegexTest(middlewareConfig.RouteMatch, requestRoute)
}
