package mock

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dhuan/mock/internal/record"
	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
	. "github.com/dhuan/mock/pkg/mock"
)

type ReadFileFunc = func(name string) ([]byte, error)

type ExecFunc = func(command string, env map[string]string) (*ExecResult, error)

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
	envVars map[string]string,
	endpointParams map[string]string,
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
			request,
			requestRecord,
			requestBody,
			readFile,
			exec,
			state,
			matchingResponseIf.Response,
			resolveResponseStatusCode(matchingResponseIf.ResponseStatusCode),
			endpointConfig,
			matchingResponseIf,
			hasResponseIf,
			envVars,
			endpointParams,
		)
	}

	return resolveEndpointResponseInternal(
		request,
		requestRecord,
		requestBody,
		readFile,
		exec,
		state,
		endpointConfig.Response,
		resolveResponseStatusCode(endpointConfig.ResponseStatusCode),
		endpointConfig,
		matchingResponseIf,
		hasResponseIf,
		envVars,
		endpointParams,
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

	for responseIfKey := range endpointConfig.ResponseIf {
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
	request *http.Request,
	requestRecord *types.RequestRecord,
	requestBody []byte,
	readFile ReadFileFunc,
	exec ExecFunc,
	state *types.State,
	response types.EndpointConfigResponse,
	responseStatusCode int,
	endpointConfig *types.EndpointConfig,
	responseIf *types.ResponseIf,
	hasResponseIf bool,
	envVars map[string]string,
	endpointParams map[string]string,
) (*Response, error, map[string]string) {
	errorMetadata := make(map[string]string)
	endpointConfigContentType := resolveEndpointConfigContentType(response)
	headers := make(map[string]string)
	utils.JoinMap(headers, endpointConfig.Headers)
	utils.JoinMap(headers, endpointConfig.HeadersBase)
	responseStr := utils.ReplaceVars(string(response), endpointParams, utils.ToDolarSignWithWrapVariablePlaceHolder)
	responseStr = utils.ReplaceVars(responseStr, envVars, utils.ToDolarSignWithWrapVariablePlaceHolder)

	if hasResponseIf {
		headers = make(map[string]string)
		utils.JoinMap(headers, endpointConfig.HeadersBase)
		utils.JoinMap(headers, responseIf.Headers)
	}

	if endpointConfigContentType == types.Endpoint_content_type_unknown {
		return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, nil, errorMetadata
	}

	if endpointConfigContentType == types.Endpoint_content_type_plaintext {
		return &Response{[]byte(utils.Unquote(responseStr)), endpointConfigContentType, responseStatusCode, headers}, nil, errorMetadata
	}

	requestVariables, err := buildVars(
		state,
		endpointConfigContentType,
		responseStatusCode,
		requestRecord,
		request,
		requestBody,
	)
	if err != nil {
		return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
	}

	if endpointConfigContentType == types.Endpoint_content_type_file {
		responseFile := fmt.Sprintf(
			"%s/%s",
			state.ConfigFolderPath,
			strings.Replace(responseStr, "file:", "", -1),
		)
		fileContent, err := readFile(responseFile)
		if errors.Is(err, ErrResponseFileDoesNotExist) {
			errorMetadata["file"] = responseFile
		}
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		responseContent := utils.ReplaceVars(string(fileContent), requestVariables, utils.ToDolarSignWithWrapVariablePlaceHolder)
		responseContent = utils.ReplaceVars(responseContent, endpointParams, utils.ToDolarSignWithWrapVariablePlaceHolder)
		responseContent = utils.ReplaceVars(responseContent, envVars, utils.ToDolarSignWithWrapVariablePlaceHolder)

		return &Response{
			[]byte(responseContent),
			endpointConfigContentType,
			responseStatusCode,
			headers}, nil, errorMetadata
	}

	if endpointConfigContentType == types.Endpoint_content_type_shell {
		scriptFilePath := fmt.Sprintf(
			"%s/%s",
			state.ConfigFolderPath,
			strings.Replace(responseStr, "sh:", "", -1),
		)

		if len(endpointParams) > 0 {
			addUrlParamsToRequestVariables(requestVariables, endpointParams)
		}

		bodyFile, err := utils.CreateTempFile(string(requestBody))
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		responseStatusCodeFile, err := utils.CreateTempFile("")
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		responseHeadersFile, err := utils.CreateTempFile("")
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		headersFile, err := utils.CreateTempFile(utils.ToHeadersText(requestRecord.Headers))
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		fileVars := map[string]string{
			"MOCK_REQUEST_HEADERS":      headersFile,
			"MOCK_REQUEST_BODY":         bodyFile,
			"MOCK_RESPONSE_HEADERS":     responseHeadersFile,
			"MOCK_RESPONSE_STATUS_CODE": responseStatusCodeFile,
		}

		utils.JoinMap(requestVariables, fileVars)

		execResult, err := exec(fmt.Sprintf("sh %s", scriptFilePath), requestVariables)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		extraHeaders, err := extractHeadersFromFile(responseHeadersFile, readFile)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		extraHeadersKeys := utils.GetSortedKeys(extraHeaders)
		for _, headerKey := range extraHeadersKeys {
			headers[headerKey] = extraHeaders[headerKey]
		}

		statusCode, err := extractStatusCodeFromFile(responseStatusCodeFile, readFile)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		response := &Response{execResult.Output, endpointConfigContentType, responseStatusCode, headers}

		if statusCode != -1 {
			response.StatusCode = statusCode
		}

		return response, nil, errorMetadata
	}

	if endpointConfigContentType == types.Endpoint_content_type_exec {
		execCommand := strings.Replace(responseStr, "exec:", "", -1)
		tempShellScriptFile, err := utils.CreateTempFile(execCommand)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		if len(endpointParams) > 0 {
			addUrlParamsToRequestVariables(requestVariables, endpointParams)
		}

		bodyFile, err := utils.CreateTempFile(string(requestBody))
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		responseStatusCodeFile, err := utils.CreateTempFile("")
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		responseHeadersFile, err := utils.CreateTempFile("")
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		headersFile, err := utils.CreateTempFile(utils.ToHeadersText(requestRecord.Headers))
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		fileVars := map[string]string{
			"MOCK_REQUEST_HEADERS":      headersFile,
			"MOCK_REQUEST_BODY":         bodyFile,
			"MOCK_RESPONSE_HEADERS":     responseHeadersFile,
			"MOCK_RESPONSE_STATUS_CODE": responseStatusCodeFile,
		}

		utils.JoinMap(requestVariables, fileVars)

		execResult, err := exec(fmt.Sprintf("sh %s", tempShellScriptFile), requestVariables)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		extraHeaders, err := extractHeadersFromFile(responseHeadersFile, readFile)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		extraHeadersKeys := utils.GetSortedKeys(extraHeaders)
		for _, headerKey := range extraHeadersKeys {
			headers[headerKey] = extraHeaders[headerKey]
		}

		statusCode, err := extractStatusCodeFromFile(responseStatusCodeFile, readFile)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		response := &Response{execResult.Output, endpointConfigContentType, responseStatusCode, headers}

		if statusCode != -1 {
			response.StatusCode = statusCode
		}

		return response, nil, errorMetadata
	}

	if endpointConfigContentType == types.Endpoint_content_type_fileserver {
		staticFilesPath := fmt.Sprintf(
			"%s/%s",
			state.ConfigFolderPath,
			strings.Replace(responseStr, "fs:", "", -1),
		)

		fileRequested, ok := endpointParams["*"]
		if !ok {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, errors.New("Failed to capture file name."), errorMetadata
		}

		filePath := fmt.Sprintf("%s/%s", staticFilesPath, fileRequested)

		fileContent, err := readFile(filePath)
		if err != nil {
			errorMetadata["file"] = fileRequested

			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		return &Response{fileContent, endpointConfigContentType, responseStatusCode, headers}, nil, errorMetadata
	}

	if endpointConfigContentType == types.Endpoint_content_type_json {
		var jsonParsed interface{}
		err := json.Unmarshal([]byte(responseStr), &jsonParsed)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		jsonEncoded, err := json.Marshal(jsonParsed)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, err, errorMetadata
		}

		jsonEncodedModified := []byte(utils.ReplaceVars(string(jsonEncoded), requestVariables, utils.ToDolarSignWithWrapVariablePlaceHolder))

		return &Response{jsonEncodedModified, endpointConfigContentType, responseStatusCode, headers}, nil, errorMetadata
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

	if utils.BeginsWith(string(response), "exec:") {
		return types.Endpoint_content_type_exec
	}

	if utils.BeginsWith(string(response), "fs:") {
		return types.Endpoint_content_type_fileserver
	}

	if utils.BeginsWith(string(response), "{") {
		return types.Endpoint_content_type_json
	}

	return types.Endpoint_content_type_plaintext
}

func extractHeadersFromFile(filePath string, readFile ReadFileFunc) (map[string]string, error) {
	headers := make(map[string]string)

	fileContent, err := readFile(filePath)

	if err != nil {
		return headers, err
	}

	fileContentText := utils.RemoveEmptyLines(string(fileContent))

	if fileContentText == "" {
		return headers, nil
	}

	headerLines := strings.Split(fileContentText, "\n")

	for i := range headerLines {
		headerKey, headerValue, ok := parseHeaderLine(headerLines[i])
		if !ok {
			continue
		}

		headers[headerKey] = headerValue
	}

	return headers, nil
}

func extractStatusCodeFromFile(filePath string, readFile ReadFileFunc) (int, error) {
	fileContent, err := readFile(filePath)

	if err != nil {
		return -1, err
	}

	fileContentText := utils.RemoveEmptyLines(string(fileContent))

	if fileContentText == "" {
		return -1, nil
	}

	statusCodeParsed, err := strconv.Atoi(fileContentText)
	if err != nil {
		return -1, err
	}

	return statusCodeParsed, nil
}

func parseHeaderLine(text string) (string, string, bool) {
	splitResult := strings.Split(text, ":")

	if len(splitResult) < 2 {
		return "", "", false
	}

	return splitResult[0], strings.Join(splitResult[1:], ":"), true
}

func addUrlParamsToRequestVariables(requestVariables, endpointParams map[string]string) {
	endpointParamKeys := utils.GetSortedKeys(endpointParams)

	for _, key := range endpointParamKeys {
		keyTransformed := fmt.Sprintf("MOCK_ROUTE_PARAM_%s", strings.ToUpper(key))

		requestVariables[keyTransformed] = endpointParams[key]
	}
}

func buildVars(
	state *types.State,
	endpointConfigContentType types.Endpoint_content_type,
	responseStatusCode int,
	requestRecord *types.RequestRecord,
	request *http.Request,
	requestBody []byte,
) (map[string]string, error) {
	endpoint := utils.ReplaceRegex(request.URL.Path, []string{"^/"}, "")
	mockHost := fmt.Sprintf("localhost:%s", state.ListenPort)
	querystring := requestRecord.Querystring
	protocol := "http://"
	if request.TLS != nil {
		protocol = "https://"
	}

	return map[string]string{
		"MOCK_HOST":                mockHost,
		"MOCK_REQUEST_HOST":        request.Host,
		"MOCK_REQUEST_URL":         fmt.Sprintf("%s%s%s", protocol, request.Host, request.URL.Path),
		"MOCK_REQUEST_ENDPOINT":    endpoint,
		"MOCK_REQUEST_METHOD":      request.Method,
		"MOCK_REQUEST_QUERYSTRING": querystring,
	}, nil
}
