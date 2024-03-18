package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dhuan/mock/internal/mock"
	"github.com/dhuan/mock/internal/utils"
)

var forwardCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		request, valid, err := buildRequestFromMockEnvVars()
		if !valid {
			exitWithError("Something went wrong. \"forward\" is supposed to be used within Response Shell Scripts. Check the manual for more details.")
		}
		if err != nil {
			panic(err)
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			panic(err)
		}

		responseVars := mock.BuildResponseVars(response)
		responseVarKeys := utils.GetSortedKeys(responseVars)
		for _, key := range responseVarKeys {
			fmt.Printf("%s=%s\n", key, responseVars[key])
		}
	},
}

var validHttpMethods []string = []string{
	"GET",
	"POST",
	"DELETE",
	"PATCH",
	"PUT",
	"OPTIONS",
}

func buildRequestFromMockEnvVars() (*http.Request, bool, error) {
	var headersPlainText string
	var method string
	var endpoint string
	var querystring string

	envValid := validateEnv(map[string]*envValidationConfig{
		"MOCK_REQUEST_HEADERS":     &envValidationConfig{variable: &headersPlainText, f: pointsToFile},
		"MOCK_REQUEST_METHOD":      &envValidationConfig{variable: &method, f: isStringAny(validHttpMethods)},
		"MOCK_REQUEST_ENDPOINT":    &envValidationConfig{variable: &endpoint, f: isStringWithText},
		"MOCK_REQUEST_QUERYSTRING": &envValidationConfig{variable: &querystring, f: optionalString},
	})
	if !envValid {
		return nil, false, nil
	}

	request, err := http.NewRequest(method)
	if err != nil {
		return nil, true, err
	}

	return request, true, nil
}

type envValidationConfig struct {
	variable *string
	f        func(string, bool) (string, bool)
}

func validateEnv(envConfig map[string]*envValidationConfig) bool {
	keys := utils.GetSortedKeys(envConfig)
	for _, key := range keys {
		value := os.Getenv(key)
		exists := value != ""

		newValue, valid := envConfig[key].f(value, exists)

		if !valid {
			return false
		}

		*envConfig[key].variable = newValue
	}

	return true
}

func isStringAny(list []string) func(string, bool) (string, bool) {
	return func(value string, exists bool) (string, bool) {
		return value, exists && (utils.IndexOf(list, value) > -1)
	}
}

func isStringWithText(value string, exists bool) (string, bool) {
	return value, exists && strings.TrimSpace(value) != ""
}

func pointsToFile(value string, exists bool) (string, bool) {
	if !exists {
		return value, false
	}

	fileContent, err := os.ReadFile(value)
	if err != nil {
		return value, false
	}

	return string(fileContent), true
}

func optionalString(value string, exists bool) (string, bool) {
	return value, true
}
