package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type MockConfig struct {
	Endpoints []EndpointConfig `json:"endpoints"`
}

type State struct {
	RequestRecordDirectoryPath string
}

type EndpointConfig struct {
	Route   string `json:"route"`
	Method  string `json:"method"`
	Content string `json:"content"`
}

type RequestRecord struct {
	Headers http.Header `json:"headers"`
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	config := &MockConfig{
		Endpoints: []EndpointConfig{
			EndpointConfig{"/foobar", "POST", "./foobar.txt"},
		},
	}

	tempDir, err := mktempDir()
	fmt.Println(fmt.Sprintf("!!!!!!!!!!!!!! %s", tempDir))
	if err != nil {
		panic(err)
	}
	state := &State{tempDir}

	for _, endpointConfig := range config.Endpoints {
		if strings.ToLower(endpointConfig.Method) == "get" {
			router.Get(endpointConfig.Route, newEndpointHandler(state, &endpointConfig))
		}

		if strings.ToLower(endpointConfig.Method) == "post" {
			router.Post(endpointConfig.Route, newEndpointHandler(state, &endpointConfig))
		}

		if strings.ToLower(endpointConfig.Method) == "patch" {
			router.Patch(endpointConfig.Route, newEndpointHandler(state, &endpointConfig))
		}

		if strings.ToLower(endpointConfig.Method) == "put" {
			router.Put(endpointConfig.Route, newEndpointHandler(state, &endpointConfig))
		}
	}

	http.ListenAndServe(":3000", router)
}

func newEndpointHandler(state *State, endpointConfig *EndpointConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileContent, err := os.ReadFile(endpointConfig.Content)
		if err != nil {
			panic(err)
		}

		requestRecord := buildRequestRecord(r)
		requestRecordJson, err := buildRequestRecordJson(requestRecord)
		if err != nil {
			panic(err)
		}
		requestRecordFileName := fmt.Sprintf("%s_%s", nowStr(), buildEndpointId(endpointConfig))
		requestRecordFilePath := fmt.Sprintf("%s/%s", state.RequestRecordDirectoryPath, requestRecordFileName)
		if err = writeNewFile(requestRecordFilePath, requestRecordJson); err != nil {
			panic(err)
		}

		w.Write(fileContent)
	}
}

func buildEndpointId(endpointConfig *EndpointConfig) string {
	return strings.ReplaceAll(endpointConfig.Route, "/", "__")
}

func mktempDir() (string, error) {
	result, err := exec.Command("mktemp", "-d").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(result), "\n"), nil
}

func nowStr() string {
	now := time.Now()

	return fmt.Sprint(now.Unix())
}

func writeNewFile(filePath string, fileContent []byte) error {
	err := os.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		return err
	}

	return nil
}

func buildRequestRecord(r *http.Request) *RequestRecord {
	return &RequestRecord{
		Headers: r.Header,
	}
}

func buildRequestRecordJson(requestRecord *RequestRecord) ([]byte, error) {
	return json.Marshal(requestRecord)
}
