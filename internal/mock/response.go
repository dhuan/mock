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
	endpointConfigContentType := resolveEndpointConfigContentType(response)
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

	if endpointConfigContentType == types.Endpoint_content_type_unknown {
		return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, errorMetadata, nil
	}

	if endpointConfigContentType == types.Endpoint_content_type_plaintext {
		mockVars := buildMockVariablesForPlainTextResponse(requestBody)

		utils.JoinMap(requestVariables, mockVars)

		addUrlParamsToRequestVariables(requestVariables, endpointParams)
		response := utils.Unquote(responseStr)
		response = utils.ReplaceVars(response, requestVariables, utils.ToDolarSignWithWrapVariablePlaceHolder)
		response = utils.ReplaceVars(response, endpointParams, utils.ToDolarSignWithWrapVariablePlaceHolder)

		return &Response{[]byte(response), endpointConfigContentType, responseStatusCode, headers}, errorMetadata, nil
	}

	if endpointConfigContentType == types.Endpoint_content_type_file {
		responseFile := extractFilePathFromResponseString(responseStr, state.ConfigFolderPath)

		fileContent, err := readFile(responseFile)
		if errors.Is(err, ErrResponseFileDoesNotExist) {
			errorMetadata["file"] = responseFile
		}
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, errorMetadata, err
		}

		responseContent := utils.ReplaceVars(string(fileContent), requestVariables, utils.ToDolarSignWithWrapVariablePlaceHolder)
		responseContent = utils.ReplaceVars(responseContent, endpointParams, utils.ToDolarSignWithWrapVariablePlaceHolder)
		responseContent = utils.ReplaceVars(responseContent, envVars, utils.ToDolarSignWithWrapVariablePlaceHolder)

		return &Response{
			[]byte(responseContent),
			endpointConfigContentType,
			responseStatusCode,
			headers}, errorMetadata, nil
	}

	if endpointConfigContentType == types.Endpoint_content_type_shell {
		scriptFilePath := extractFilePathFromResponseString(responseStr, state.ConfigFolderPath)

		if len(endpointParams) > 0 {
			addUrlParamsToRequestVariables(requestVariables, endpointParams)
		}

		fileVars, hf, err := buildHandlerFiles(requestBody, requestRecord)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, errorMetadata, err
		}

		utils.JoinMap(requestVariables, fileVars)

		log.Printf("Executing shell script located in %s", scriptFilePath)

		execResult, err := exec(
			fmt.Sprintf("sh %s", scriptFilePath),
			&ExecOptions{
				Env: requestVariables,
			},
		)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, errorMetadata, err
		}

		printOutExecOutputIfNecessary(execResult)

		response, err := extractModifiedResponse(hf, readFile, endpointConfigContentType, headers, responseStatusCode)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, errorMetadata, err
		}

		return response, errorMetadata, nil
	}

	if endpointConfigContentType == types.Endpoint_content_type_exec {
		execCommand := strings.Replace(responseStr, "exec:", "", -1)
		tempShellScriptFile, err := utils.CreateTempFile([]byte(execCommand))
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, errorMetadata, err
		}

		if len(endpointParams) > 0 {
			addUrlParamsToRequestVariables(requestVariables, endpointParams)
		}

		fileVars, hf, err := buildHandlerFiles(requestBody, requestRecord)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, errorMetadata, err
		}

		utils.JoinMap(requestVariables, fileVars)

		log.Printf("Executing command: %s", execCommand)

		execResult, err := exec(fmt.Sprintf("sh %s", tempShellScriptFile), &ExecOptions{
			Env: requestVariables,
		})
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, errorMetadata, err
		}

		printOutExecOutputIfNecessary(execResult)

		response, err := extractModifiedResponse(hf, readFile, endpointConfigContentType, headers, responseStatusCode)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, errorMetadata, err
		}

		return response, errorMetadata, nil
	}

	if endpointConfigContentType == types.Endpoint_content_type_fileserver {
		staticFilesPath := extractFilePathFromResponseString(responseStr, state.ConfigFolderPath)

		fileRequested, ok := endpointParams["*"]
		if !ok {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, errorMetadata, errors.New("Failed to capture file name.")
		}

		filePath := fmt.Sprintf("%s/%s", staticFilesPath, fileRequested)

		fileContent, err := readFile(filePath)
		if err != nil {
			errorMetadata["file"] = fileRequested

			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, errorMetadata, err
		}

		return &Response{fileContent, endpointConfigContentType, responseStatusCode, headers}, errorMetadata, nil
	}

	if endpointConfigContentType == types.Endpoint_content_type_json {
		var jsonParsed interface{}
		err := json.Unmarshal([]byte(responseStr), &jsonParsed)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, errorMetadata, err
		}

		jsonEncoded, err := json.Marshal(jsonParsed)
		if err != nil {
			return &Response{[]byte(""), endpointConfigContentType, responseStatusCode, headers}, errorMetadata, err
		}

		jsonEncodedModified := []byte(utils.ReplaceVars(string(jsonEncoded), requestVariables, utils.ToDolarSignWithWrapVariablePlaceHolder))

		return &Response{jsonEncodedModified, endpointConfigContentType, responseStatusCode, headers}, errorMetadata, nil
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

	return splitResult[0], strings.Join(splitResult[1:], ":"), true
}

func addUrlParamsToRequestVariables(requestVariables, endpointParams map[string]string) {
	endpointParamKeys := utils.GetSortedKeys(endpointParams)

	for _, key := range endpointParamKeys {
		keyTransformed := fmt.Sprintf("MOCK_ROUTE_PARAM_%s", strings.ToUpper(key))

		requestVariables[keyTransformed] = endpointParams[key]
	}
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
