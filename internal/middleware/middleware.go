package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dhuan/mock/internal/mock"
	. "github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
	. "github.com/dhuan/mock/pkg/mock"
)

type MiddlewareRunResult struct {
	Body       []byte
	Headers    map[string]string
	StatusCode int
}

func RunMiddleware(
	exec mock.ExecFunc,
	readFile ReadFileFunc,
	configPath string,
	middlewareConfigs []MiddlewareConfig,
	responseBody []byte,
	responseHeaders map[string]string,
	responseStatusCode int,
	request *http.Request,
	endpointParams map[string]string,
	vars map[string]string,
	createTempFile func([]byte) (string, error),
) (*MiddlewareRunResult, error) {
	result := &MiddlewareRunResult{}
	result.Body = responseBody
	result.Headers = make(map[string]string)
	result.StatusCode = responseStatusCode

	if len(middlewareConfigs) == 0 {
		return result, nil
	}

	responseFiles, err := buildResponseFiles(responseBody, responseHeaders, responseStatusCode)
	if err != nil {
		return result, err
	}

	for i := range middlewareConfigs {
		envVars := map[string]string{
			"MOCK_RESPONSE_BODY":        responseFiles.body,
			"MOCK_RESPONSE_HEADERS":     responseFiles.headers,
			"MOCK_RESPONSE_STATUS_CODE": responseFiles.statusCode,
		}

		for key := range endpointParams {
			endpointParamKey := fmt.Sprintf("MOCK_ROUTE_PARAM_%s", strings.ToUpper(key))

			envVars[endpointParamKey] = endpointParams[key]
		}

		for key := range vars {
			envVars[key] = vars[key]
		}

		tempScriptFilePath, err := createTempFile([]byte(middlewareConfigs[i].Exec))
		if err != nil {
			return result, err
		}

		execResult, err := exec(
			fmt.Sprintf("sh %s", tempScriptFilePath),
			&mock.ExecOptions{
				Env:        envVars,
				WorkingDir: configPath,
			},
		)

		if len(execResult.Output) > 0 {
			log.Printf("Middleware execution output:\n%s", string(execResult.Output))
		}

		if err != nil {
			return result, err
		}
	}

	return readResponseFiles(responseFiles, readFile)
}

func readResponseFiles(
	rf *responseFiles,
	readFile ReadFileFunc,
) (*MiddlewareRunResult, error) {
	result := &MiddlewareRunResult{}

	resultResponseBody, err := readFile(rf.body)
	if err != nil {
		return result, err
	}

	resultResponseHeaders, err := readFile(rf.headers)
	if err != nil {
		return result, err
	}

	resultResponseStatusCode, err := readFile(rf.statusCode)
	if err != nil {
		return result, err
	}

	responseStatusCodeParsed := bytesToInt(resultResponseStatusCode, 200)

	result.Body = resultResponseBody
	result.StatusCode = responseStatusCodeParsed
	result.Headers = utils.ExtractHeadersFromText(resultResponseHeaders)

	return result, nil
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

func GetMiddlewareForRequest(
	middlewareConfigs []MiddlewareConfig,
	r *http.Request,
	requestRecord *RequestRecord,
	requestRecords []RequestRecord,
	verifyCondition func(requestRecord *RequestRecord, condition *Condition, requestRecords []RequestRecord) bool,
) []MiddlewareConfig {
	middlewares := make([]MiddlewareConfig, 0)

	for i := range middlewareConfigs {
		if middlewareConfigs[i].Condition != nil && !verifyCondition(requestRecord, middlewareConfigs[i].Condition, requestRecords) {
			continue
		}

		if routeMatch(r, &middlewareConfigs[i]) {
			middlewares = append(middlewares, middlewareConfigs[i])
		}
	}

	return middlewares
}

func routeMatch(r *http.Request, middlewareConfig *MiddlewareConfig) bool {
	if middlewareConfig.RouteMatch == "*" || middlewareConfig.RouteMatch == "" {
		return true
	}

	requestRoute := utils.ReplaceRegex(r.URL.Path, []string{"^/"}, "")

	return utils.RegexTest(middlewareConfig.RouteMatch, requestRoute)
}

type responseFiles struct {
	body       string
	headers    string
	statusCode string
}

func buildResponseFiles(
	responseBody []byte,
	responseHeaders map[string]string,
	responseStatusCode int,
) (*responseFiles, error) {
	result := &responseFiles{}
	resultResponseHeaders := []byte(utils.ToHeadersText(toHttpHeaders(responseHeaders)) + "\n")
	resultResponseBody := make([]byte, len(responseBody))
	copy(resultResponseBody, responseBody)
	resultResponseStatusCode := []byte(fmt.Sprintf("%d", responseStatusCode))

	responseBodyFile, err := utils.CreateTempFile(resultResponseBody)
	if err != nil {
		return result, err
	}

	responseHeadersFile, err := utils.CreateTempFile(resultResponseHeaders)
	if err != nil {
		return result, err
	}

	responseStatusCodeFile, err := utils.CreateTempFile(resultResponseStatusCode)
	if err != nil {
		return result, err
	}

	result.body = responseBodyFile
	result.headers = responseHeadersFile
	result.statusCode = responseStatusCodeFile

	return result, nil
}
