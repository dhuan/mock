package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dhuan/mock/internal/args2config"
	mockMiddleware "github.com/dhuan/mock/internal/middleware"
	"github.com/dhuan/mock/internal/mock"
	"github.com/dhuan/mock/internal/mockfs"
	"github.com/dhuan/mock/internal/record"
	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		endpointsFromCommandLine := args2config.ParseEndpoints(os.Args)

		config, err := resolveConfig(flagConfig)
		if err != nil {
			exitWithError(err.Error())
		}

		hasBaseApi, baseApi := resolveBaseApi(flagBaseApi, config)

		if hasBaseApi && !baseApiIsValid(baseApi) {
			exitWithError(fmt.Sprintf("Base API is not valid: %s\nSet it as a valid domain name such as google.com", baseApi))
		}

		if flagConfig == "" && len(endpointsFromCommandLine) == 0 && !hasBaseApi {
			exitWithError(cmd.UsageString())
		}

		allEndpoints := mergeEndpoints(config.Endpoints, endpointsFromCommandLine)

		config.Endpoints = allEndpoints

		middlewaresFromCommandLine := args2config.ParseMiddlewares(os.Args)
		mergeMiddlewares(config, middlewaresFromCommandLine)

		router := chi.NewRouter()
		router.Use(middleware.Logger)
		router.MethodNotAllowed(onMethodNotAllowed(flagCors))

		prepareConfig(config)

		endpointConfigErrors, err := mock.ValidateEndpointConfigs(
			config.Endpoints,
			readFile,
			filepath.Dir(flagConfig),
		)
		if err != nil {
			panic(err)
		}

		if len(endpointConfigErrors) > 0 {
			log.Println("mock can't be started. The following errors were found in your configuration:")
			log.Println("")
			displayEndpointConfigErrors(endpointConfigErrors, config.Endpoints)

			os.Exit(1)
		}

		tempDir, err := utils.MktempDir()
		log.Printf("Temporary folder created for Request Records: %s", tempDir)
		if err != nil {
			panic(err)
		}
		state := &types.State{
			RequestRecordDirectoryPath: tempDir,
			ConfigFolderPath:           filepath.Dir(flagConfig),
			ListenPort:                 flagPort,
		}
		mockFs := mockfs.MockFs{State: state}

		router.Use(handleOptions(flagCors, state, config, mockFs, router))

		router.NotFound(onNotFound(flagCors, hasBaseApi, baseApi, state, config, mockFs))

		for i := range config.Endpoints {
			endpointConfig := config.Endpoints[i]
			route := fmt.Sprintf("/%s", endpointConfig.Route)

			if endpointConfig.Method == "get" || endpointConfig.Method == "" {
				router.Get(route, newEndpointHandler(state, config.Middlewares, &endpointConfig, mockFs, flagDelay, config))
			}

			if endpointConfig.Method == "post" {
				router.Post(route, newEndpointHandler(state, config.Middlewares, &endpointConfig, mockFs, flagDelay, config))
			}

			if endpointConfig.Method == "patch" {
				router.Patch(route, newEndpointHandler(state, config.Middlewares, &endpointConfig, mockFs, flagDelay, config))
			}

			if endpointConfig.Method == "put" {
				router.Put(route, newEndpointHandler(state, config.Middlewares, &endpointConfig, mockFs, flagDelay, config))
			}

			if endpointConfig.Method == "delete" {
				router.Delete(route, newEndpointHandler(state, config.Middlewares, &endpointConfig, mockFs, flagDelay, config))
			}
		}

		router.Post("/__mock__/assert", mockApiHandler(mockFs, state, config))
		router.Post("/__mock__/reset", resetApiHandler(mockFs, state, config))

		port, errorMessage := resolvePort(flagPort)
		if errorMessage != "" {
			exitWithError(errorMessage)
		}

		log.Printf("Starting server on port %d.", port)

		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Println("An error occurred while starting up the server.")
			log.Fatalln(err)
		}

		log.Println("Server started.")

		if err = http.Serve(listener, router); err != nil {
			log.Println("An error occurred while starting up the server.")
			log.Fatalln(err)
		}
	},
}

func resolvePort(flagPort string) (int, string) {
	portIsDefinedByUser := flagPort != "UNSET"

	if !portIsDefinedByUser {
		return utils.GetFreePort(), ""
	}

	port, err := strconv.Atoi(flagPort)
	if err != nil {
		return -1, fmt.Sprintf("Port %s is not valid!", flagPort)
	}

	portIsFree := utils.IsPortFree(port)

	if !portIsFree {
		return -1, fmt.Sprintf("Port %d is not available! Start mock without specifying a port number to get a random available port.", port)
	}

	if portIsFree {
		return port, ""
	}

	return -1, "Failed to resolve port!"
}

func resolveEndpointErrorDescription(endpointConfigError *mock.EndpointConfigError) string {
	if endpointConfigError.Code == mock.EndpointConfigErrorCode_EndpointDuplicate {
		duplicateIndexParsed, err := strconv.Atoi(endpointConfigError.Metadata["duplicate_index"])
		if err != nil {
			panic(err)
		}
		duplicateIndex := duplicateIndexParsed + 1

		return fmt.Sprintf(
			"This endpoint has a duplicate (Endpoint #%d). A combination of route and method must be unique. If you're looking to define different responses for the same endpoint/method, look for \"Conditional Responses\" in the documentation.",
			duplicateIndex,
		)
	}

	if endpointConfigError.Code == mock.EndpointConfigErrorCode_InvalidMethod {
		return fmt.Sprintf(
			"The given method, \"%s\" , is invalid. The available HTTP Methods you can use are POST, GET, PUT, PATCH, and DELETE.",
			endpointConfigError.Metadata["method"],
		)
	}

	if endpointConfigError.Code == mock.EndpointConfigErrorCode_RouteWithQuerystring {
		return "Routes cannot have querystrings. Read about \"response_if\" in the documentation to learn how to set Conditional Responses based on querystrings."
	}

	if endpointConfigError.Code == mock.EndpointConfigErrorCode_FileUnreadable {
		return fmt.Sprintf(
			"The file provided for the endpoint's response is either unreadable or does not exist: %s",
			endpointConfigError.Metadata["file_path"],
		)
	}

	panic("Failed to resolve endpoint error description.")
}

func displayEndpointConfigErrors(endpointConfigErrors []mock.EndpointConfigError, endpointConfigs []types.EndpointConfig) {
	for i, endpointConfigError := range endpointConfigErrors {
		endpointRoute := endpointConfigs[endpointConfigError.EndpointIndex].Route
		endpointMethod := endpointConfigs[endpointConfigError.EndpointIndex].Method

		log.Printf(
			"%d: Endpoint #%d (%s %s):\n%s\n\n",
			i+1,
			endpointConfigError.EndpointIndex+1,
			endpointMethod,
			endpointRoute,
			resolveEndpointErrorDescription(&endpointConfigError),
		)
	}
}

func readFile(name string) ([]byte, error) {
	_, err := os.Stat(name)
	if errors.Is(err, os.ErrNotExist) {
		return []byte(""), mock.ErrResponseFileDoesNotExist
	}

	return os.ReadFile(name)
}

func execute(command string, options *mock.ExecOptions) (*mock.ExecResult, error) {
	commandStrings := utils.ToCommandStrings(command)
	commandName, commandParams := utils.ToCommandParams(commandStrings)
	cmd := exec.Command(commandName, commandParams...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, utils.ParseEnv(options.Env)...)

	if options.WorkingDir != "" {
		cmd.Dir = options.WorkingDir
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Run()

	return &mock.ExecResult{
		Output: out.Bytes(),
	}, nil
}

func newEndpointHandler(
	state *types.State,
	middlewareConfigs []types.MiddlewareConfig,
	endpointConfig *types.EndpointConfig,
	mockFs types.MockFs,
	delay int64,
	config *MockConfig,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		envVars := getAllEnvVars()
		endpointParams := getEndpointParams(r)

		requestRoute := utils.ReplaceRegex(r.URL.Path, []string{"^/"}, "")
		requestRecords, err := mockFs.GetRecordsMatchingRoute(requestRoute)
		if err != nil {
			panic(err)
		}

		requestRecord, requestBody, err := record.BuildRequestRecord(r, endpointParams)
		if err != nil {
			panic(err)
		}

		response, errorMetadata, err := mock.ResolveEndpointResponse(
			readFile,
			execute,
			requestBody,
			state,
			endpointConfig,
			envVars,
			endpointParams,
			requestRecord,
			requestRecords,
			config.Base,
		)
		if errors.Is(err, mock.ErrResponseFileDoesNotExist) {
			log.Printf("Tried to read file that does not exist: %s", errorMetadata["file"])
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(
				"File does not exist: %s",
				errorMetadata["file"],
			)))

			return
		}
		if err != nil {
			panic(err)
		}
		if response.EndpointContentType == types.Endpoint_content_type_unknown {
			log.Printf("Failed to resolve endpoint content type for route %s", endpointConfig.Route)

			return
		}
		if response.EndpointContentType == types.Endpoint_content_type_json {
			w.Header().Add("Content-Type", "application/json")
		}

		if flagCors {
			setCorsHeadersToResponse(response)
		}

		err = mockFs.StoreRequestRecord(requestRecord, endpointConfig)
		if err != nil {
			panic(err)
		}

		if delay > 0 {
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}

		response = handleMiddleware(
			state,
			r,
			response,
			endpointParams,
			config,
			requestRecord,
			requestRecords,
			requestBody,
			map[string]string{},
		)

		addHeaders(w, response)

		w.WriteHeader(response.StatusCode)
		w.Write(response.Body)
	}
}

func mockApiHandler(mockFs types.MockFs, state *types.State, config *MockConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assertOptions, err := mock.ParseAssertRequest(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(400)

			return
		}

		validationErrors, err := mock.Validate(mockFs, assertOptions)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)

			return
		}

		response := MockApiResponse{validationErrors}
		responseJson, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
			w.Header().Add("Content-Type", "application/json")
			w.Write(responseJson)

			return
		}

		statusCode := 200
		if len(*validationErrors) > 0 {
			statusCode = 400
		}
		w.WriteHeader(statusCode)

		w.Header().Add("Content-Type", "application/json")
		w.Write(responseJson)
	}
}

func resetApiHandler(mockFs types.MockFs, state *types.State, config *MockConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := mockFs.RemoveAllRequestRecords()
		if err != nil {
			log.Println("Failed to remove Request Records.")
			log.Println(err)

			w.WriteHeader(400)

			return
		}

		w.WriteHeader(200)
	}
}

func prepareConfig(mockConfig *MockConfig) {
	for i, endpoint := range mockConfig.Endpoints {
		mockConfig.Endpoints[i].Method = strings.ToLower(mockConfig.Endpoints[i].Method)
		mockConfig.Endpoints[i].Route = utils.ReplaceRegex(endpoint.Route, []string{`^\/`}, "")
	}
}

func resolveConfig(configPath string) (*MockConfig, error) {
	if configPath == "" {
		return &MockConfig{
			Endpoints: []types.EndpointConfig{},
		}, nil
	}

	var mockConfig MockConfig

	configFileContent, err := os.ReadFile(configPath)
	if err != nil {
		return &mockConfig, fmt.Errorf("Unable to read configuration file \"%s\". Make sure it exists and/or is readable, then try again.", configPath)
	}

	if err = json.Unmarshal(configFileContent, &mockConfig); err != nil {
		return &mockConfig, err
	}

	return &mockConfig, nil
}

func addHeaders(w http.ResponseWriter, response *mock.Response) {
	for headerKey := range response.Headers {
		w.Header().Add(headerKey, response.Headers[headerKey])
	}
}

func exitWithError(errorMessage string) {
	fmt.Println(errorMessage)

	os.Exit(1)
}

func onNotFound(
	corsEnabled,
	hasBaseApi bool,
	baseApi string,
	state *types.State,
	config *MockConfig,
	mockFs types.MockFs,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpointConfig, ok := findEndpointConfigByRoute(config.Endpoints, fmt.Sprintf("%s/*", utils.ReplaceRegex(r.URL.Path, []string{"^/"}, "")))
		if ok {
			newEndpointHandler(state, config.Middlewares, endpointConfig, mockFs, flagDelay, config)(w, r)

			return
		}

		endpointParams := getEndpointParams(r)

		envVars := make(map[string]string)
		envVars["MOCK_REQUEST_NOT_FOUND"] = "true"

		if !hasBaseApi {
			requestRecord, requestBody, err := record.BuildRequestRecord(r, endpointParams)
			if err != nil {
				panic(err)
			}

			requestRoute := utils.ReplaceRegex(r.URL.Path, []string{"^/"}, "")
			requestRecords, err := mockFs.GetRecordsMatchingRoute(requestRoute)

			response := &mock.Response{
				Body:                []byte(""),
				EndpointContentType: types.Endpoint_content_type_unknown,
				StatusCode:          405,
				Headers:             nil,
			}

			if corsEnabled {
				setCorsHeadersToResponse(response)
			}

			response = handleMiddleware(
				state,
				r,
				response,
				endpointParams,
				config,
				requestRecord,
				requestRecords,
				requestBody,
				envVars,
			)

			addHeaders(w, response)

			w.WriteHeader(response.StatusCode)
			w.Write(response.Body)

			return
		}

		r2, err := utils.CloneRequest(r)
		if err != nil {
			panic(err)
		}

		response, _, err := sendRequestForBaseApi(baseApi, r2)
		if err != nil {
			panic(err)
		}

		requestRecord, requestBody, err := record.BuildRequestRecord(r, endpointParams)
		if err != nil {
			panic(err)
		}

		requestRoute := utils.ReplaceRegex(r.URL.Path, []string{"^/"}, "")
		requestRecords, err := mockFs.GetRecordsMatchingRoute(requestRoute)
		if err != nil {
			panic(err)
		}

		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		mockResponse := buildResponse(response, responseBody)

		if corsEnabled {
			setCorsHeadersToResponse(mockResponse)
		}

		for key := range mockResponse.Headers {
			if strings.ToLower(key) == "content-length" {
				delete(mockResponse.Headers, key)
			}
		}

		middlewareHandlerExtraVars := map[string]string{
			"MOCK_BASE_API_RESPONSE": "true",
		}

		mockResponse = handleMiddleware(
			state,
			r,
			mockResponse,
			endpointParams,
			config,
			requestRecord,
			requestRecords,
			requestBody,
			middlewareHandlerExtraVars,
		)

		forwardResponse(mockResponse, w)
	}
}

func buildResponse(response *http.Response, responseBody []byte) *mock.Response {
	headers := make(map[string]string)
	for key := range response.Header {
		headers[key] = strings.Join(response.Header[key], " ")
	}

	return &mock.Response{
		Body:                responseBody,
		EndpointContentType: types.Endpoint_content_type_unknown,
		StatusCode:          response.StatusCode,
		Headers:             headers,
	}
}

func handleMiddleware(
	state *types.State,
	r *http.Request,
	response *mock.Response,
	endpointParams map[string]string,
	config *MockConfig,
	requestRecord *types.RequestRecord,
	requestRecords []types.RequestRecord,
	requestBody []byte,
	extraVars map[string]string,
) *mock.Response {
	middlewareConfigsForRequest := mockMiddleware.GetMiddlewareForRequest(config.Middlewares, r, requestRecord, requestRecords, mock.VerifyCondition)
	hasMiddleware := len(middlewareConfigsForRequest) > 0

	vars, err := mock.BuildVars(state, response.StatusCode, requestRecord, requestRecords, requestBody, config.Base)
	if err != nil {
		panic(err)
	}

	for key := range extraVars {
		vars[key] = extraVars[key]
	}

	httpHeaders := toHttpHeaders(response.Headers)

	responseTransformed := response.Body
	if hasMiddleware {
		middlewareRunResult, err := mockMiddleware.RunMiddleware(
			execute,
			readFile,
			state.ConfigFolderPath,
			middlewareConfigsForRequest,
			responseTransformed,
			&httpHeaders,
			response.StatusCode,
			r,
			endpointParams,
			vars,
			utils.CreateTempFile,
			requestRecord,
		)
		if err != nil {
			panic(err)
		}

		response.Headers = middlewareRunResult.Headers
		response.StatusCode = middlewareRunResult.StatusCode

		return &mock.Response{
			Body:                middlewareRunResult.Body,
			EndpointContentType: types.Endpoint_content_type_unknown,
			StatusCode:          middlewareRunResult.StatusCode,
			Headers:             middlewareRunResult.Headers,
		}
	}

	return response
}

func sendRequestForBaseApi(baseApi string, r *http.Request) (*http.Response, []byte, error) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	route := r.URL.Path
	querystring := ""
	if r.URL.RawQuery != "" {
		querystring = fmt.Sprintf("?%s", r.URL.RawQuery)
	}
	protocol, host := parseBaseApi(r.TLS != nil, baseApi)
	url := fmt.Sprintf("%s://%s%s%s", protocol, host, route, querystring)

	requestCloned, err := http.NewRequest(r.Method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}

	requestCloned.Header = r.Header.Clone()

	response, err := client.Do(requestCloned)

	return response, requestBody, err
}

func parseBaseApi(currentRequestIsHttps bool, baseApi string) (string, string) {
	currentRequestProtocol := "http"
	if currentRequestIsHttps {
		currentRequestProtocol = "https"
	}

	if !utils.RegexTest("^http", baseApi) {
		return currentRequestProtocol, baseApi
	}

	protocol := "http"
	if utils.RegexTest("^https", baseApi) {
		protocol = "https"
	}

	domain := extractDomain(baseApi)

	return protocol, domain
}

func extractDomain(url string) string {
	split := strings.Split(url, "//")

	if len(split) < 2 {
		panic("Something went wrong while extracting domain.")
	}

	return split[1]
}

func forwardResponse(response *mock.Response, w http.ResponseWriter) {
	for key := range response.Headers {
		w.Header().Add(key, response.Headers[key])
	}

	w.WriteHeader(response.StatusCode)

	if len(response.Body) > 0 {
		w.Write(response.Body)
	}
}

func onMethodNotAllowed(corsEnabled bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !corsEnabled {
			w.WriteHeader(405)

			return
		}

		setCorsHeaders(w)

		w.WriteHeader(405)
	}
}

func matchRouteAnyMethod(router *chi.Mux, path string) bool {
	ctx := chi.NewRouteContext()

	for _, method := range strings.Split("POST,GET,PUT,PATCH,DELETE,OPTIONS", ",") {
		if router.Match(ctx, method, path) {
			return true
		}
	}

	return false
}

func handleOptions(
	corsEnabled bool,
	state *types.State,
	config *MockConfig,
	mockFs types.MockFs,
	router *chi.Mux,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			match := matchRouteAnyMethod(router, r.URL.Path)

			if strings.ToLower(r.Method) != "options" || !match || !corsEnabled {
				next.ServeHTTP(w, r)

				return
			}

			headers := make(map[string]string)

			if corsEnabled {
				for key, value := range corsHeaders {
					headers[key] = value
				}
			}

			endpointParams := getEndpointParams(r)

			requestRoute := utils.ReplaceRegex(r.URL.Path, []string{"^/"}, "")
			requestRecords, err := mockFs.GetRecordsMatchingRoute(requestRoute)
			if err != nil {
				panic(err)
			}

			requestRecord, requestBody, err := record.BuildRequestRecord(r, endpointParams)
			if err != nil {
				panic(err)
			}

			response := &mock.Response{
				Body:                nil,
				EndpointContentType: types.Endpoint_content_type_unknown,
				StatusCode:          200,
				Headers:             headers,
			}

			response = handleMiddleware(
				state,
				r,
				response,
				endpointParams,
				config,
				requestRecord,
				requestRecords,
				requestBody,
				map[string]string{},
			)

			for key := range response.Headers {
				w.Header().Set(key, response.Headers[key])
			}

			w.WriteHeader(response.StatusCode)

			if response.Body != nil {
				w.Write(response.Body)
			}

			next.ServeHTTP(w, r)
		})
	}
}

var corsHeaders map[string]string = map[string]string{
	"Access-Control-Allow-Origin":      "*",
	"Access-Control-Allow-Credentials": "true",
	"Access-Control-Allow-Headers":     "*",
	"Access-Control-Allow-Methods":     "POST, GET, OPTIONS, PUT, DELETE",
}

func setCorsHeaders(w http.ResponseWriter) {
	for key, value := range corsHeaders {
		w.Header().Set(key, value)
	}
}

func setCorsHeadersToResponse(response *mock.Response) {
	if response.Headers == nil {
		response.Headers = make(map[string]string)
	}

	for key, value := range corsHeaders {
		response.Headers[key] = value
	}
}

func getEndpointParams(r *http.Request) map[string]string {
	params := make(map[string]string)
	chiUrlParams := chi.RouteContext(r.Context()).URLParams

	for i := range chiUrlParams.Keys {
		params[chiUrlParams.Keys[i]] = chiUrlParams.Values[i]
	}

	return params
}

func getAllEnvVars() map[string]string {
	result := make(map[string]string)
	envPairs := os.Environ()
	envKeys := make([]string, len(envPairs))

	for i := range envPairs {
		envKeys[i] = utils.ReplaceRegex(envPairs[i], []string{
			"=.*$",
		}, "")
	}

	for i := range envKeys {
		result[envKeys[i]] = os.Getenv(envKeys[i])
	}

	return result
}

func mergeEndpoints(a, b []types.EndpointConfig) []types.EndpointConfig {
	return append(a, b...)
}

func mergeMiddlewares(config *MockConfig, middlewares []types.MiddlewareConfig) {
	config.Middlewares = append(config.Middlewares, middlewares...)
}

func resolveBaseApi(flagBaseApi string, config *MockConfig) (bool, string) {
	if flagBaseApi != "" {
		baseApi := formatBaseApi(flagBaseApi)

		config.Base = baseApi

		return true, baseApi
	}

	if config.Base != "" {
		return true, formatBaseApi(config.Base)
	}

	return false, ""
}

func formatBaseApi(baseApi string) string {
	return utils.ReplaceRegex(baseApi, []string{"/$"}, "")
}

func baseApiIsValid(baseApi string) bool {
	baseApi = removeWebProtocolAndPort(baseApi)

	regexMustPass := []string{
		"^[a-zA-Z0-9]{1}[a-zA-Z0-9-_.]{1,}$",
		"[a-zA-Z0-9]{1}$",
	}

	for _, regex := range regexMustPass {
		if !utils.RegexTest(regex, baseApi) {
			return false
		}
	}

	return true
}

func removeWebProtocolAndPort(url string) string {
	return utils.ReplaceRegex(url, []string{
		"^https?://",
		":[0-9]{1,}$",
	}, "")
}

func toHttpHeaders(m map[string]string) http.Header {
	result := make(http.Header)

	for key := range m {
		result[key] = []string{m[key]}
	}

	return result
}

func findEndpointConfigByRoute(endpointConfigs []types.EndpointConfig, route string) (*types.EndpointConfig, bool) {
	for i := range endpointConfigs {
		if endpointConfigs[i].Route == route {
			return &endpointConfigs[i], true
		}
	}

	return &types.EndpointConfig{}, false
}
