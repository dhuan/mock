package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dhuan/mock/internal/mock"
	"github.com/dhuan/mock/internal/mockfs"
	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nsf/jsondiff"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := resolveConfig(flagConfig)
		if err != nil {
			panic(err)
		}

		router := chi.NewRouter()
		router.Use(middleware.Logger)

		prepareConfig(config)

		endpointConfigErrors, err := mock.ValidateEndpointConfigs(config.Endpoints)
		if err != nil {
			panic(err)
		}

		if len(endpointConfigErrors) > 0 {
			fmt.Println("mock can't be started. The following errors were found in your configuration:")
			fmt.Println("")
			displayEndpointConfigErrors(endpointConfigErrors, config.Endpoints)

			os.Exit(1)
		}

		tempDir, err := utils.MktempDir()
		fmt.Println(fmt.Sprintf("Temporary folder created for Request Records: %s", tempDir))
		if err != nil {
			panic(err)
		}
		state := &types.State{
			RequestRecordDirectoryPath: tempDir,
			ConfigFolderPath:           filepath.Dir(flagConfig),
		}
		mockFs := mockfs.MockFs{State: state}

		for i, _ := range config.Endpoints {
			endpointConfig := config.Endpoints[i]
			route := fmt.Sprintf("/%s", endpointConfig.Route)

			if endpointConfig.Method == "get" {
				router.Get(route, newEndpointHandler(state, &endpointConfig, mockFs))
			}

			if endpointConfig.Method == "post" {
				router.Post(route, newEndpointHandler(state, &endpointConfig, mockFs))
			}

			if endpointConfig.Method == "patch" {
				router.Patch(route, newEndpointHandler(state, &endpointConfig, mockFs))
			}

			if endpointConfig.Method == "put" {
				router.Put(route, newEndpointHandler(state, &endpointConfig, mockFs))
			}

			if endpointConfig.Method == "delete" {
				router.Delete(route, newEndpointHandler(state, &endpointConfig, mockFs))
			}
		}

		router.Post("/__mock__/assert", mockApiHandler(mockFs, state, config))
		router.Post("/__mock__/reset", resetApiHandler(mockFs, state, config))

		fmt.Println(fmt.Sprintf("Mock server is listening on port %s.", flagPort))

		http.ListenAndServe(fmt.Sprintf(":%s", flagPort), router)
	},
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
		return fmt.Sprintf(
			"Routes cannot have querystrings. Read about \"response_if\" in the documentation to learn how to set Conditional Responses based on querystrings.",
		)
	}

	panic("Failed to resolve endpoint error description.")
}

func displayEndpointConfigErrors(endpointConfigErrors []mock.EndpointConfigError, endpointConfigs []types.EndpointConfig) {
	for i, endpointConfigError := range endpointConfigErrors {
		endpointRoute := endpointConfigs[endpointConfigError.EndpointIndex].Route
		endpointMethod := endpointConfigs[endpointConfigError.EndpointIndex].Method

		fmt.Println(
			fmt.Sprintf(
				"%d: Endpoint #%d (%s %s):\n%s\n",
				i+1,
				endpointConfigError.EndpointIndex+1,
				endpointMethod,
				endpointRoute,
				resolveEndpointErrorDescription(&endpointConfigError),
			),
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

func newEndpointHandler(state *types.State, endpointConfig *types.EndpointConfig, mockFs types.MockFs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		response, err, errorMetadata := mock.ResolveEndpointResponse(readFile, r, requestBody, state, endpointConfig)
		if errors.Is(err, mock.ErrResponseFileDoesNotExist) {
			log.Println(fmt.Sprintf("Tried to read file that does not exist: %s", errorMetadata["file"]))
			w.WriteHeader(400)
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
			fmt.Println(fmt.Sprintf("Failed to resolve endpoint content type for route %s", endpointConfig.Route))

			return
		}
		if response.EndpointContentType == types.Endpoint_content_type_json {
			w.Header().Add("Content-Type", "application/json")
		}

		addHeaders(w, response)

		err = mockFs.StoreRequestRecord(r, requestBody, endpointConfig)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(response.StatusCode)
		w.Write(response.Body)
	}
}

func mockApiHandler(mockFs types.MockFs, state *types.State, config *MockConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assertConfig, err := mock.ParseAssertRequest(r)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(400)

			return
		}

		validationErrors, err := mock.Validate(mockFs, jsonValidate, assertConfig)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)

			return
		}

		response := MockApiResponse{validationErrors}
		responseJson, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
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

func jsonValidate(jsonA map[string]interface{}, jsonB map[string]interface{}) bool {
	a, err := json.Marshal(jsonA)
	if err != nil {
		return false
	}

	b, err := json.Marshal(jsonB)
	if err != nil {
		return false
	}

	options := jsondiff.DefaultJSONOptions()
	result, _ := jsondiff.Compare(a, b, &options)

	return result == jsondiff.FullMatch
}

func resolveConfig(configPath string) (*MockConfig, error) {
	var mockConfig MockConfig

	configFileContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		return &mockConfig, err
	}

	if err = json.Unmarshal(configFileContent, &mockConfig); err != nil {
		return &mockConfig, err
	}

	return &mockConfig, nil
}

func addHeaders(w http.ResponseWriter, response *mock.Response) {
	for headerKey, _ := range response.Headers {
		w.Header().Add(headerKey, response.Headers[headerKey])
	}
}
