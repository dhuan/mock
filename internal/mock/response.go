package mock

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
	. "github.com/dhuan/mock/pkg/mock"
)

type ExecOptions struct {
	Env        map[string]string
	WorkingDir string
}

type ExecFunc = func(command string, options *ExecOptions) (*ExecResult, error)

type Response struct {
	Body                []byte
	EndpointContentType types.Endpoint_content_type
	StatusCode          int
	Headers             map[string]string
}

type ExecResult struct {
	Output []byte
}

type responseResolverData struct {
	state                     *types.State
	envVars                   map[string]string
	requestBody               []byte
	responseStr               string
	requestVariables          map[string]string
	endpointParams            map[string]string
	endpointConfigContentType types.Endpoint_content_type
	responseStatusCode        int
	headers                   map[string]string
	errorMetadata             map[string]string
	requestRecord             *types.RequestRecord
}

func ResolveEndpointResponse(
	readFile types.ReadFileFunc,
	exec ExecFunc,
	requestBody []byte,
	state *types.State,
	endpointConfig *types.EndpointConfig,
	envVars map[string]string,
	endpointParams map[string]string,
	requestRecord *types.RequestRecord,
	requestRecords []types.RequestRecord,
	baseApi string,
) (*Response, map[string]string, error) {
	hasResponseIf := len(endpointConfig.ResponseIf) > 0
	matchingResponseIf := &types.ResponseIf{}

	if hasResponseIf {
		matchingResponseIfB, foundMatchingResponseIf := resolveResponseIf(requestRecord, endpointConfig, requestRecords)
		matchingResponseIf = matchingResponseIfB
		hasResponseIf = foundMatchingResponseIf
	}

	if hasResponseIf {
		return resolveEndpointResponseInternal(
			requestRecord,
			requestRecords,
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
			baseApi,
		)
	}

	return resolveEndpointResponseInternal(
		requestRecord,
		requestRecords,
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
		baseApi,
	)
}

func resolveResponseStatusCode(statusCode int) int {
	if statusCode < 1 {
		return 200
	}

	return statusCode
}

func resolveResponseIf(requestRecord *types.RequestRecord, endpointConfig *types.EndpointConfig, requestRecords []types.RequestRecord) (*types.ResponseIf, bool) {
	matchingResponseIfs := make([]int, 0)

	for responseIfKey := range endpointConfig.ResponseIf {
		responseIf := endpointConfig.ResponseIf[responseIfKey]
		matches := VerifyCondition(requestRecord, responseIf.Condition, requestRecords)

		if matches {
			matchingResponseIfs = append(matchingResponseIfs, responseIfKey)
		}
	}

	if len(matchingResponseIfs) == 0 {
		return &types.ResponseIf{}, false
	}

	return &endpointConfig.ResponseIf[matchingResponseIfs[0]], true
}

func VerifyCondition(requestRecord *types.RequestRecord, condition *Condition, requestRecords []types.RequestRecord) bool {
	conditionFunction := resolveAssertTypeFunc(condition.Type, requestRecords)
	validationErrors, err := conditionFunction(requestRecord, requestRecords, condition)
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
		return VerifyCondition(requestRecord, condition.And, requestRecords)
	}

	if !result && hasOr {
		return VerifyCondition(requestRecord, condition.Or, requestRecords)
	}

	if !result && !hasOr {
		return false
	}

	return false
}

func resolveEndpointResponseInternal(
	requestRecord *types.RequestRecord,
	requestRecords []types.RequestRecord,
	requestBody []byte,
	readFile types.ReadFileFunc,
	exec ExecFunc,
	state *types.State,
	response types.EndpointConfigResponse,
	responseStatusCode int,
	endpointConfig *types.EndpointConfig,
	responseIf *types.ResponseIf,
	hasResponseIf bool,
	envVars map[string]string,
	endpointParams map[string]string,
	baseApi string,
) (*Response, map[string]string, error) {
	errorMetadata := make(map[string]string)
	endpointConfigContentType := resolveEndpointConfigContentType(string(response))
	headers := make(map[string]string)
	utils.JoinMap(headers, endpointConfig.Headers)
	utils.JoinMap(headers, endpointConfig.HeadersBase)
	responseStr := utils.ReplaceVars(string(response), endpointParams, utils.ToDolarSignWithWrapVariablePlaceHolder)
	responseStr = utils.ReplaceVars(responseStr, envVars, utils.ToDolarSignWithWrapVariablePlaceHolder)

	requestVariables, err := BuildVars(
		state,
		responseStatusCode,
		requestRecord,
		requestRecords,
		requestBody,
		baseApi,
	)
	if err != nil {
		return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, errorMetadata, err
	}

	if hasResponseIf {
		headers = make(map[string]string)
		utils.JoinMap(headers, endpointConfig.HeadersBase)
		utils.JoinMap(headers, responseIf.Headers)
	}

	rrd := &responseResolverData{
		state,
		envVars,
		requestBody,
		responseStr,
		requestVariables,
		endpointParams,
		endpointConfigContentType,
		responseStatusCode,
		headers,
		errorMetadata,
		requestRecord,
	}

	if endpointConfigContentType == types.Endpoint_content_type_unknown {
		return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, errorMetadata, nil
	}

	if endpointConfigContentType == types.Endpoint_content_type_plaintext {
		return plainTextResponse(rrd, readFile, exec)
	}

	if endpointConfigContentType == types.Endpoint_content_type_file {
		return fileResponse(rrd, readFile, exec)
	}

	if endpointConfigContentType == types.Endpoint_content_type_shell {
		return shellResponse(rrd, readFile, exec)
	}

	if endpointConfigContentType == types.Endpoint_content_type_exec {
		return execResponse(rrd, readFile, exec)
	}

	if endpointConfigContentType == types.Endpoint_content_type_fileserver {
		return fileServerResponse(rrd, readFile, exec)
	}

	if endpointConfigContentType == types.Endpoint_content_type_json {
		return jsonResponse(rrd, readFile, exec)
	}

	return &Response{[]byte(""), types.Endpoint_content_type_unknown, responseStatusCode, headers}, errorMetadata, nil
}

func extractFilePathFromResponseString(responseStr, configFolderPath string) string {
	splitResult := strings.Split(responseStr, ":")

	if len(splitResult) < 2 {
		return "unknown"
	}

	filePath := splitResult[1]

	if utils.BeginsWith(filePath, "/") {
		return filePath
	}

	return fmt.Sprintf("%s/%s", configFolderPath, filePath)
}

func resolveEndpointConfigContentType(response string) types.Endpoint_content_type {
	if utils.BeginsWith(response, "file:") {
		return types.Endpoint_content_type_file
	}

	if utils.BeginsWith(response, "sh:") {
		return types.Endpoint_content_type_shell
	}

	if utils.BeginsWith(response, "exec:") {
		return types.Endpoint_content_type_exec
	}

	if utils.BeginsWith(response, "fs:") {
		return types.Endpoint_content_type_fileserver
	}

	if (utils.BeginsWith(response, "{") || utils.BeginsWith(response, "[")) && isValidJson(response) {
		return types.Endpoint_content_type_json
	}

	return types.Endpoint_content_type_plaintext
}

func isValidJson(txt string) bool {
	var jsonParsed interface{}

	return json.Unmarshal([]byte(txt), &jsonParsed) == nil
}

func ExtractHeadersFromFile(filePath string, readFile types.ReadFileFunc) (map[string]string, error) {
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

func extractStatusCodeFromFile(filePath string, readFile types.ReadFileFunc) (int, error) {
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

	key := splitResult[0]
	value := strings.TrimSpace(strings.Join(splitResult[1:], ":"))

	return key, value, true
}

func addUrlParamsToRequestVariables(requestVariables, endpointParams map[string]string) {
	endpointParamKeys := utils.GetSortedKeys(endpointParams)

	for _, key := range endpointParamKeys {
		keyTransformed := fmt.Sprintf("MOCK_ROUTE_PARAM_%s", strings.ToUpper(key))

		requestVariables[keyTransformed] = endpointParams[key]
	}

	routeParamsEncoded, err := utils.EncodeToJsonBase64(endpointParams)
	if err != nil {
		log.Printf("Failed to encode route params to json/base64: %s", err.Error())
	}

	requestVariables["MOCK_ROUTE_PARAMS"] = routeParamsEncoded
}

func printOutExecOutputIfNecessary(execResult *ExecResult) {
	output := string(execResult.Output)

	if strings.TrimSpace(output) == "" {
		return
	}

	log.Printf(fmt.Sprintf("Output from program execution:\n\n%s\n", output))
}

type handlerFiles struct {
	headers            string
	body               string
	responseBody       string
	responseHeaders    string
	responseStatusCode string
}

func buildMockVariablesForPlainTextResponse(requestBody []byte) map[string]string {
	return map[string]string{
		"MOCK_REQUEST_BODY": string(requestBody),
	}
}

func buildHandlerFiles(requestBody []byte, requestRecord *types.RequestRecord) (map[string]string, *handlerFiles, error) {
	bodyFile, err := utils.CreateTempFile(requestBody)
	if err != nil {
		return map[string]string{}, &handlerFiles{}, err
	}

	responseStatusCodeFile, err := utils.CreateTempFile([]byte(""))
	if err != nil {
		return map[string]string{}, &handlerFiles{}, err
	}

	responseHeadersFile, err := utils.CreateTempFile([]byte(""))
	if err != nil {
		return map[string]string{}, &handlerFiles{}, err
	}

	responseBodyFile, err := utils.CreateTempFile([]byte(""))
	if err != nil {
		return map[string]string{}, &handlerFiles{}, err
	}

	headersFile, err := utils.CreateTempFile([]byte(utils.ToHeadersText(requestRecord.Headers)))
	if err != nil {
		return map[string]string{}, &handlerFiles{}, err
	}

	return map[string]string{
			"MOCK_REQUEST_HEADERS":      headersFile,
			"MOCK_REQUEST_BODY":         bodyFile,
			"MOCK_RESPONSE_BODY":        responseBodyFile,
			"MOCK_RESPONSE_HEADERS":     responseHeadersFile,
			"MOCK_RESPONSE_STATUS_CODE": responseStatusCodeFile,
		}, &handlerFiles{
			headersFile,
			bodyFile,
			responseBodyFile,
			responseHeadersFile,
			responseStatusCodeFile,
		}, nil
}

func extractModifiedResponse(
	hf *handlerFiles,
	readFile types.ReadFileFunc,
	endpointConfigContentType types.Endpoint_content_type,
	headers map[string]string,
	fallbackStatusCode int,
) (*Response, error) {
	extraHeaders, err := ExtractHeadersFromFile(hf.responseHeaders, readFile)
	if err != nil {
		return &Response{[]byte(""), types.Endpoint_content_type_unknown, fallbackStatusCode, nil}, err
	}

	extraHeadersKeys := utils.GetSortedKeys(extraHeaders)
	for _, headerKey := range extraHeadersKeys {
		headers[headerKey] = extraHeaders[headerKey]
	}

	statusCode, err := extractStatusCodeFromFile(hf.responseStatusCode, readFile)
	if err != nil {
		return &Response{[]byte(""), types.Endpoint_content_type_unknown, fallbackStatusCode, nil}, err
	}

	bodyContent, err := readFile(hf.responseBody)
	if err != nil {
		return &Response{[]byte(""), types.Endpoint_content_type_unknown, fallbackStatusCode, nil}, err
	}

	response := &Response{bodyContent, endpointConfigContentType, statusCode, headers}

	if statusCode == -1 {
		response.StatusCode = fallbackStatusCode
	}

	return response, nil
}

func plainTextResponse(
	data *responseResolverData,
	readFile types.ReadFileFunc,
	exec ExecFunc,
) (*Response, map[string]string, error) {
	mockVars := buildMockVariablesForPlainTextResponse(data.requestBody)

	utils.JoinMap(data.requestVariables, mockVars)

	addUrlParamsToRequestVariables(data.requestVariables, data.endpointParams)
	response := utils.Unquote(data.responseStr)
	response = utils.ReplaceVars(response, data.requestVariables, utils.ToDolarSignWithWrapVariablePlaceHolder)
	response = utils.ReplaceVars(response, data.endpointParams, utils.ToDolarSignWithWrapVariablePlaceHolder)

	return &Response{[]byte(response), data.endpointConfigContentType, data.responseStatusCode, data.headers}, data.errorMetadata, nil
}

func fileResponse(
	data *responseResolverData,
	readFile types.ReadFileFunc,
	exec ExecFunc,
) (*Response, map[string]string, error) {
	responseFile := extractFilePathFromResponseString(data.responseStr, data.state.ConfigFolderPath)

	fileContent, err := readFile(responseFile)
	if errors.Is(err, ErrResponseFileDoesNotExist) {
		data.errorMetadata["file"] = responseFile
	}
	if err != nil {
		return &Response{[]byte(""), data.endpointConfigContentType, data.responseStatusCode, data.headers}, data.errorMetadata, err
	}

	responseContent := utils.ReplaceVars(string(fileContent), data.requestVariables, utils.ToDolarSignWithWrapVariablePlaceHolder)
	responseContent = utils.ReplaceVars(responseContent, data.endpointParams, utils.ToDolarSignWithWrapVariablePlaceHolder)
	responseContent = utils.ReplaceVars(responseContent, data.envVars, utils.ToDolarSignWithWrapVariablePlaceHolder)

	return &Response{
		[]byte(responseContent),
		data.endpointConfigContentType,
		data.responseStatusCode,
		data.headers}, data.errorMetadata, nil
}

func execResponse(
	data *responseResolverData,
	readFile types.ReadFileFunc,
	exec ExecFunc,
) (*Response, map[string]string, error) {
	return execWrapper(data, readFile, exec, func() (string, error) {
		execCommand := strings.Replace(data.responseStr, "exec:", "", -1)

		tempShellScriptFile, err := utils.CreateTempFile([]byte(execCommand))
		if err != nil {
			return "", err
		}

		command := fmt.Sprintf("sh %s", tempShellScriptFile)

		log.Printf("Executing command: %s", command)

		return command, nil
	})
}

func shellResponse(
	data *responseResolverData,
	readFile types.ReadFileFunc,
	exec ExecFunc,
) (*Response, map[string]string, error) {
	return execWrapper(data, readFile, exec, func() (string, error) {
		scriptFilePath := extractFilePathFromResponseString(data.responseStr, data.state.ConfigFolderPath)

		log.Printf("Executing shell script located in %s", scriptFilePath)

		return fmt.Sprintf("sh %s", scriptFilePath), nil
	})
}

func execWrapper(
	data *responseResolverData,
	readFile types.ReadFileFunc,
	exec ExecFunc,
	f func() (string, error),
) (*Response, map[string]string, error) {
	command, err := f()
	if err != nil {
		return &Response{[]byte(""), data.endpointConfigContentType, data.responseStatusCode, data.headers}, data.errorMetadata, err
	}

	if len(data.endpointParams) > 0 {
		addUrlParamsToRequestVariables(data.requestVariables, data.endpointParams)
	}

	fileVars, hf, err := buildHandlerFiles(data.requestBody, data.requestRecord)
	if err != nil {
		return &Response{[]byte(""), data.endpointConfigContentType, data.responseStatusCode, data.headers}, data.errorMetadata, err
	}

	utils.JoinMap(data.requestVariables, fileVars)

	execResult, err := exec(command, &ExecOptions{Env: data.requestVariables})
	if err != nil {
		return &Response{[]byte(""), data.endpointConfigContentType, data.responseStatusCode, data.headers}, data.errorMetadata, err
	}

	printOutExecOutputIfNecessary(execResult)

	response, err := extractModifiedResponse(hf, readFile, data.endpointConfigContentType, data.headers, data.responseStatusCode)
	if err != nil {
		return &Response{[]byte(""), data.endpointConfigContentType, data.responseStatusCode, data.headers}, data.errorMetadata, err
	}

	return response, data.errorMetadata, nil
}

func fileServerResponse(
	data *responseResolverData,
	readFile types.ReadFileFunc,
	exec ExecFunc,
) (*Response, map[string]string, error) {
	staticFilesPath := extractFilePathFromResponseString(data.responseStr, data.state.ConfigFolderPath)

	fileRequested, ok := data.endpointParams["*"]
	if !ok {
		return &Response{[]byte(""), data.endpointConfigContentType, data.responseStatusCode, data.headers}, data.errorMetadata, errors.New("Failed to capture file name.")
	}

	filePath := fmt.Sprintf("%s/%s", staticFilesPath, fileRequested)

	fileContent, err := readFile(filePath)
	if err != nil {
		data.errorMetadata["file"] = fileRequested

		return &Response{[]byte(""), data.endpointConfigContentType, data.responseStatusCode, data.headers}, data.errorMetadata, err
	}

	return &Response{fileContent, data.endpointConfigContentType, data.responseStatusCode, data.headers}, data.errorMetadata, nil
}

func jsonResponse(
	data *responseResolverData,
	readFile types.ReadFileFunc,
	exec ExecFunc,
) (*Response, map[string]string, error) {
	var jsonParsed interface{}
	err := json.Unmarshal([]byte(data.responseStr), &jsonParsed)
	if err != nil {
		return &Response{[]byte(""), data.endpointConfigContentType, data.responseStatusCode, data.headers}, data.errorMetadata, err
	}

	jsonEncoded, err := json.Marshal(jsonParsed)
	if err != nil {
		return &Response{[]byte(""), data.endpointConfigContentType, data.responseStatusCode, data.headers}, data.errorMetadata, err
	}

	jsonEncodedModified := []byte(utils.ReplaceVars(string(jsonEncoded), data.requestVariables, utils.ToDolarSignWithWrapVariablePlaceHolder))

	return &Response{jsonEncodedModified, data.endpointConfigContentType, data.responseStatusCode, data.headers}, data.errorMetadata, nil
}
