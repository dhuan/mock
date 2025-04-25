package mock

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

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
	cache map[string]string,
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
			cache,
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
		cache,
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
	cache map[string]string,
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
		return execResponse(rrd, readFile, exec, cache, endpointConfig)
	}

	if endpointConfigContentType == types.Endpoint_content_type_fileserver {
		return fileServerResponse(rrd, readFile, exec)
	}

	if endpointConfigContentType == types.Endpoint_content_type_json {
		return jsonResponse(rrd, readFile, exec)
	}

	return &Response{[]byte(""), types.Endpoint_content_type_unknown, responseStatusCode, headers}, errorMetadata, nil
}

func extractFilePathFromResponseString(responseStr, configFolderPath string) (string, string) {
	splitResult := strings.Split(responseStr, ":")

	if len(splitResult) < 2 {
		return "unknown", ""
	}

	filePath := splitResult[1]

	splitResult = strings.Split(filePath, ".")

	extension := ""
	if len(splitResult) > 1 {
		extension = splitResult[len(splitResult)-1]
	}

	if utils.BeginsWith(filePath, "/") {
		return filePath, extension
	}

	return fmt.Sprintf("%s/%s", configFolderPath, filePath), extension
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

	log.Printf("Output from program execution:\n\n%s\n", output)
}

type HandlerFiles struct {
	Headers            string
	Body               string
	ResponseBody       string
	ResponseHeaders    string
	ResponseStatusCode string
}

func buildMockVariablesForPlainTextResponse(requestBody []byte) map[string]string {
	return map[string]string{
		"MOCK_REQUEST_BODY": string(requestBody),
	}
}

func BuildHandlerFiles(
	requestBody []byte,
	requestRecord *types.RequestRecord,
	responseBody []byte,
	responseHeaders *http.Header,
	responseStatusCode int,
) (map[string]string, *HandlerFiles, error) {
	bodyFile, err := utils.CreateTempFile(requestBody)
	if err != nil {
		return map[string]string{}, &HandlerFiles{}, err
	}

	responseStatusCodeFile, err := utils.CreateTempFile([]byte(strconv.Itoa(responseStatusCode)))
	if err != nil {
		return map[string]string{}, &HandlerFiles{}, err
	}

	headersText := ""
	if len(*responseHeaders) > 0 {
		headersText = utils.ToHeadersText(*responseHeaders)
	}

	responseHeadersFile, err := utils.CreateTempFile([]byte(headersText))
	if err != nil {
		return map[string]string{}, &HandlerFiles{}, err
	}

	responseBodyFile, err := utils.CreateTempFile(responseBody)
	if err != nil {
		return map[string]string{}, &HandlerFiles{}, err
	}

	headersFile, err := utils.CreateTempFile([]byte(utils.ToHeadersText(requestRecord.Headers)))
	if err != nil {
		return map[string]string{}, &HandlerFiles{}, err
	}

	return map[string]string{
			"MOCK_REQUEST_HEADERS":      headersFile,
			"MOCK_REQUEST_BODY":         bodyFile,
			"MOCK_RESPONSE_BODY":        responseBodyFile,
			"MOCK_RESPONSE_HEADERS":     responseHeadersFile,
			"MOCK_RESPONSE_STATUS_CODE": responseStatusCodeFile,
		}, &HandlerFiles{
			headersFile,
			bodyFile,
			responseBodyFile,
			responseHeadersFile,
			responseStatusCodeFile,
		}, nil
}

func extractModifiedResponse(
	hf *HandlerFiles,
	readFile types.ReadFileFunc,
	endpointConfigContentType types.Endpoint_content_type,
	headers map[string]string,
	fallbackStatusCode int,
) (*Response, error) {
	extraHeaders, err := ExtractHeadersFromFile(hf.ResponseHeaders, readFile)
	if err != nil {
		return &Response{[]byte(""), types.Endpoint_content_type_unknown, fallbackStatusCode, nil}, err
	}

	extraHeadersKeys := utils.GetSortedKeys(extraHeaders)
	for _, headerKey := range extraHeadersKeys {
		headers[headerKey] = extraHeaders[headerKey]
	}

	statusCode, err := extractStatusCodeFromFile(hf.ResponseStatusCode, readFile)
	if err != nil {
		return &Response{[]byte(""), types.Endpoint_content_type_unknown, fallbackStatusCode, nil}, err
	}

	bodyContent, err := readFile(hf.ResponseBody)
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
	responseFile, extension := extractFilePathFromResponseString(data.responseStr, data.state.ConfigFolderPath)

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

	if extension == "json" {
		data.headers["content-type"] = "application/json"
	}

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
	cache map[string]string,
	ec *types.EndpointConfig,
) (*Response, map[string]string, error) {
	return execWrapper(data, readFile, exec, func() (string, error) {
		var err error
		cacheKey := fmt.Sprintf("tmp_exec_file_%s_%s", ec.Method, ec.Route)

		cachedScriptFile, ok := cache[cacheKey]
		if !ok {
			execCommand := strings.Replace(data.responseStr, "exec:", "", -1)

			cachedScriptFile, err = utils.CreateTempFile([]byte(execCommand))
			if err != nil {
				return "", err
			}

			cache[cacheKey] = cachedScriptFile
		}

		command := fmt.Sprintf("sh %s", cachedScriptFile)

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
		scriptFilePath, _ := extractFilePathFromResponseString(data.responseStr, data.state.ConfigFolderPath)

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

	fileVars, hf, err := BuildHandlerFiles(data.requestBody, data.requestRecord, []byte(""), &http.Header{}, 200)
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

func fileServerDirectoryResponse(
	data *responseResolverData,
	staticFilesPath string,
	fileRequested string,
) (*Response, map[string]string, error) {
	directoryPath := fmt.Sprintf("%s/%s", staticFilesPath, fileRequested)

	files, err := scanDir(directoryPath)
	if err != nil {
		return nil, nil, err
	}

	iconDirectory := svg("icon_directory", 24, 24)
	iconFile := svg("icon_file", 24, 24)

	sort.Slice(files, func(i, j int) bool {
		if files[i].dir != !files[j].dir {
			return files[i].dir
		}

		return true
	})

	fileList := ""
	for i := range files {
		fileName := files[i].name

		href := fmt.Sprintf("/%s/%s", data.requestRecord.Route, fileName)
		href = strings.Replace(href, "//", "/", -1)

		modified := files[i].modified

		icon := iconFile
		if files[i].dir {
			icon = iconDirectory
		}

		fileList = fmt.Sprintf(`
%s
<tr>
	<td>%s </td>
	<td><a href="%s">%s</a></td>
	<!-- TD FILE SIZE -->
	<td>%s</td>
	<!-- TD FILE MODIFIED -->
	<td>%s</td>
</tr>`,
			fileList, icon, href, fileName, files[i].size, modified)
	}

	indexName := fmt.Sprintf("/%s", fileRequested)
	hasPrevDir := fileRequested != ""

	prevDirRow := ""
	if hasPrevDir {
		prevDirRow = fmt.Sprintf(`
<tr>
	<td>%s </td>
	<td><a href="%s">Parent directory</a></td>
	<td></td>
	<td></td>
</tr>
`, svg("icon_back", 24, 24), previousDirHref(data.requestRecord.Route))
	}

	result := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>mock directory list</title>
<style>
a {color: blue;}

table{width: 100%%; border-collapse: collapse;}

table tbody tr td:nth-child(1) {
	display: flex;
	justify-content: center;
}

table tr {
	border-bottom: 1px dashed #dadada;
}

table tr:hover td {
	background: #feffec;
}

table tr td {
	padding: 10px 0;
}

.footer {
	padding-top: 20px;
}
</style>
</head>
<body>
<div style="display: none">
%s
%s
%s
</div>
<h1>Index of %s</h1>
<table>
	<thead>
		<tr>
			<td width="50px"></td>
			<td>Name</td>
			<td>Size</td>
			<td>Modified</td>
		</tr>
	</thead>
	<tbody>
%s
%s
	</tbody>
</table>

<div class="footer">
Served with mock __MOCK_VERSION__
</div>
</body>
</html>
`, svg_back, svg_directory, svg_file, indexName, prevDirRow, fileList)

	return &Response{[]byte(result), data.endpointConfigContentType, data.responseStatusCode, data.headers}, data.errorMetadata, nil
}

func fileServerResponse(
	data *responseResolverData,
	readFile types.ReadFileFunc,
	exec ExecFunc,
) (*Response, map[string]string, error) {
	staticFilesPath, _ := extractFilePathFromResponseString(data.responseStr, data.state.ConfigFolderPath)

	fileRequested, ok := data.endpointParams["*"]
	if !ok {
		fileRequested = ""
	}

	filePath := fmt.Sprintf("%s/%s", staticFilesPath, fileRequested)

	fileRequestedIsDir := isDir(filePath)

	if fileRequested == "" || fileRequestedIsDir {
		return fileServerDirectoryResponse(data, staticFilesPath, fileRequested)
	}

	fileContent, err := readFile(filePath)
	if err != nil {
		data.errorMetadata["file"] = fileRequested

		return &Response{[]byte(""), data.endpointConfigContentType, data.responseStatusCode, data.headers}, data.errorMetadata, err
	}

	data.headers["Content-Type"] = resolveContentType(filePath)

	return &Response{fileContent, data.endpointConfigContentType, data.responseStatusCode, data.headers}, data.errorMetadata, nil
}

func resolveContentType(filePath string) string {
	splitResult := strings.Split(filePath, ".")

	if len(splitResult) < 2 {
		return "text/plain; charset=utf-8"
	}

	extension := splitResult[len(splitResult)-1]

	contentType, ok := content_type[extension]
	if !ok {
		return "text/plain; charset=utf-8"
	}

	return contentType
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

type fileEntry struct {
	name     string
	dir      bool
	modified string
	size     string
}

func scanDir(walkFrom string) ([]fileEntry, error) {
	files := make([]fileEntry, 0)

	i := -1
	err := filepath.Walk(walkFrom, func(path string, info os.FileInfo, err error) error {
		i++

		if walkFrom == path {
			return nil
		}

		files = append(files, fileEntry{
			filepath.Base(path),
			info.IsDir(),
			info.ModTime().Format(time.RFC822),
			fmt.Sprintf("%d", info.Size()),
		})

		if info.IsDir() && i > 0 {
			return filepath.SkipDir
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func svg(id string, w, h int) string {
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg"  width="%d"  height="%d"><use href="#%s"></use></svg>`, w, h, id)
}

func previousDirHref(route string) string {
	route = utils.ReplaceRegex(route, []string{`^\/`}, "")

	splitResult := strings.Split(route, "/")

	return fmt.Sprintf("/%s", strings.Join(splitResult[0:(len(splitResult)-1)], "/"))
}
