package middleware

import (
	"fmt"
	"net/http"

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
) ([]byte, error) {
	resultResponseBody := make([]byte, len(responseBody))
	copy(resultResponseBody, responseBody)

	if len(middlewareConfigs) == 0 {
		return responseBody, nil
	}

	for i := range middlewareConfigs {
		scriptPath := fmt.Sprintf("%s/%s", configPath, middlewareConfigs[i].ScriptPath)
		if utils.BeginsWith(middlewareConfigs[i].ScriptPath, "/") {
			scriptPath = middlewareConfigs[i].ScriptPath
		}

		responseBodyFile, err := utils.CreateTempFile(resultResponseBody)
		if err != nil {
			return responseBody, err
		}

		envVars := map[string]string{
			"MOCK_RESPONSE_BODY": responseBodyFile,
		}

		_, err = exec(fmt.Sprintf("sh %s", scriptPath), envVars)
		if err != nil {
			return responseBody, err
		}

		resultResponseBody, err = readFile(responseBodyFile)
		if err != nil {
			return responseBody, err
		}
	}

	return resultResponseBody, nil
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

	return utils.RegexTest(requestRoute, middlewareConfig.RouteMatch)
}
