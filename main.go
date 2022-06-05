package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/dhuan/mock/internal/mock"
	"github.com/dhuan/mock/internal/mockfs"
	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
	"github.com/nsf/jsondiff"
)

type MockConfig struct {
	Endpoints []types.EndpointConfig `json:"endpoints"`
}

type MockApiResponse struct {
	Pass             bool                    `json:"pass"`
	ValidationErrors *[]mock.ValidationError `json:"validation_errors"`
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	config := &MockConfig{
		Endpoints: []types.EndpointConfig{
			types.EndpointConfig{
				Route:   "/foobar",
				Method:  "POST",
				Content: "./foobar.txt",
			},
		},
	}

	prepareConfig(config)

	tempDir, err := utils.MktempDir()
	fmt.Println(fmt.Sprintf("Temporary folder created for Request Records: %s", tempDir))
	if err != nil {
		panic(err)
	}
	state := &types.State{RequestRecordDirectoryPath: tempDir}
	mockFs := mockfs.MockFs{State: state}

	for _, endpointConfig := range config.Endpoints {
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

	router.Post("/__mock__", mockApiHandler(mockFs, state, config))

	http.ListenAndServe(":3000", router)
}

func newEndpointHandler(state *types.State, endpointConfig *types.EndpointConfig, mockFs types.MockFs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileContent, err := os.ReadFile(endpointConfig.Content)
		if err != nil {
			panic(err)
		}

		err = mockFs.StoreRequestRecord(r, endpointConfig)
		if err != nil {
			panic(err)
		}

		w.Write(fileContent)
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

		pass, validationErrors, err := mock.Validate(mockFs, jsonValidate, assertConfig)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)

			return
		}

		response := MockApiResponse{pass, validationErrors}
		responseJson, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
		}

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
