package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
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

			if strings.ToLower(endpointConfig.Method) == "get" {
				router.Get(route, newEndpointHandler(state, &endpointConfig, mockFs))
			}

			if strings.ToLower(endpointConfig.Method) == "post" {
				router.Post(route, newEndpointHandler(state, &endpointConfig, mockFs))
			}

			if strings.ToLower(endpointConfig.Method) == "patch" {
				router.Patch(route, newEndpointHandler(state, &endpointConfig, mockFs))
			}

			if strings.ToLower(endpointConfig.Method) == "put" {
				router.Put(route, newEndpointHandler(state, &endpointConfig, mockFs))
			}
		}

		router.Post("/__mock__/assert", mockApiHandler(mockFs, state, config))

		fmt.Println(fmt.Sprintf("Mock server is listening on port %s.", flagPort))

		http.ListenAndServe(fmt.Sprintf(":%s", flagPort), router)
	},
}

func newEndpointHandler(state *types.State, endpointConfig *types.EndpointConfig, mockFs types.MockFs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseContent, endpointContentType, err := mock.ResolveEndpointResponse(state, endpointConfig)
		if err != nil {
			panic(err)
		}
		if endpointContentType == types.Endpoint_content_type_unknown {
			fmt.Println(fmt.Sprintf("Failed to resolve endpoint content type for route %s", endpointConfig.Route))

			return
		}
		if endpointContentType == types.Endpoint_content_type_json {
			w.Header().Add("Content-Type", "application/json")
		}

		addHeaders(w, endpointConfig)

		err = mockFs.StoreRequestRecord(r, endpointConfig)
		if err != nil {
			panic(err)
		}

		w.Write(responseContent)
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

func prepareConfig(mockConfig *MockConfig) {
	for i, endpoint := range mockConfig.Endpoints {
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

func addHeaders(w http.ResponseWriter, endpointConfig *types.EndpointConfig) {
	for headerKey, _ := range endpointConfig.Headers {
		w.Header().Add(headerKey, endpointConfig.Headers[headerKey])
	}
}
