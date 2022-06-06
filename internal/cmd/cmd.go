package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

type MockConfig struct {
	Endpoints []types.EndpointConfig `json:"endpoints"`
}

type MockApiResponse struct {
	Pass             bool                    `json:"pass"`
	ValidationErrors *[]mock.ValidationError `json:"validation_errors"`
}

var (
	flagConfig string
	flagPort   string
)

var rootCmd = &cobra.Command{
	Use: "mock",
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

		fmt.Println(fmt.Sprintf("Mock server is listening on port %s.", flagPort))

		http.ListenAndServe(fmt.Sprintf(":%s", flagPort), router)
	},
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

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&flagConfig, "config", "c", "", "configuration file")
	rootCmd.PersistentFlags().StringVarP(&flagPort, "port", "p", "3000", "port to listen on")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newEndpointHandler(state *types.State, endpointConfig *types.EndpointConfig, mockFs types.MockFs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseFile := fmt.Sprintf("%s/%s", state.ConfigFolderPath, endpointConfig.Content)
		fileContent, err := os.ReadFile(responseFile)
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
