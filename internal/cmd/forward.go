package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dhuan/mock/internal/mock"
	"github.com/dhuan/mock/internal/utils"
)

var forwardCmd = &cobra.Command{
	Use: "forward",
	Run: func(cmd *cobra.Command, args []string) {
		request, valid, err, envValidationErrors := buildRequestFromMockEnvVars()
		if !valid {
			fmt.Println(strings.Join(envValidationErrors, "\n"))

			exitWithError(
				"Something went wrong. \"forward\" is supposed to be used within Response Shell Scripts. Check the manual for more details.",
			)
		}
		if err != nil {
			panic(err)
		}

		log.Printf("Forwarding request to Base API: %s %s\n", request.Method, request.RequestURI)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			panic(err)
		}

		log.Printf("Got response from Base API: %d\n", response.StatusCode)

		responseVars, err := mock.BuildResponseVars(response)
		if err != nil {
			panic(err)
		}

		responseVarKeys := utils.GetSortedKeys(responseVars)
		for _, key := range responseVarKeys {
			fmt.Printf("%s=%s\n", key, responseVars[key])
		}
	},
}

var validHttpMethods []string = []string{
	"get",
	"post",
	"delete",
	"patch",
	"put",
	"options",
}

func buildRequestFromMockEnvVars() (*http.Request, bool, error, []string) {
	var baseApiUrl string
	var headersPlainText string
	var method string
	var endpoint string
	var querystring string

	envValid, errorMessages := validateEnv(map[string]*envValidationConfig{
		"MOCK_BASE_API":            {variable: &baseApiUrl, f: isStringWithText},
		"MOCK_REQUEST_HEADERS":     {variable: &headersPlainText, f: pointsToFile},
		"MOCK_REQUEST_METHOD":      {variable: &method, f: isStringAny(validHttpMethods)},
		"MOCK_REQUEST_ENDPOINT":    {variable: &endpoint, f: isStringWithText},
		"MOCK_REQUEST_QUERYSTRING": {variable: &querystring, f: optionalString},
	})
	if !envValid {
		return nil, false, nil, errorMessages
	}

	method = strings.ToUpper(method)
	url := fmt.Sprintf("%s/%s", baseApiUrl, endpoint)
	if querystring != "" {
		url = fmt.Sprintf("%s?%s", url, querystring)
	}

	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, true, err, []string{}
	}

	headers, err := mock.ExtractHeadersFromFile(os.Getenv("MOCK_REQUEST_HEADERS"), readFile)
	if err != nil {
		return nil, true, err, []string{}
	}
	for headerKey, headerValue := range headers {
		request.Header.Add(headerKey, headerValue)
	}

	return request, true, nil, []string{}
}

type envValidationConfig struct {
	variable *string
	f        func(string, bool) (string, bool, string)
}

func validateEnv(envConfig map[string]*envValidationConfig) (bool, []string) {
	errorMessages := make([]string, 0)

	keys := utils.GetSortedKeys(envConfig)
	for _, key := range keys {
		value := os.Getenv(key)
		exists := value != ""

		newValue, valid, errorMessage := envConfig[key].f(value, exists)

		if !valid {
			errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", key, errorMessage))

			continue
		}

		*envConfig[key].variable = newValue
	}

	return len(errorMessages) == 0, errorMessages
}

func isStringAny(list []string) func(string, bool) (string, bool, string) {
	return func(value string, exists bool) (string, bool, string) {
		return value, exists && (utils.IndexOf(list, value) > -1), fmt.Sprintf(
			"is not set as any of: %s", strings.Join(list, ","),
		)
	}
}

func isStringWithText(value string, exists bool) (string, bool, string) {
	return value, exists && strings.TrimSpace(value) != "", "is not a string with text"
}

func pointsToFile(value string, exists bool) (string, bool, string) {
	if !exists {
		return value, false, "does not exist"
	}

	fileContent, err := os.ReadFile(value)
	if err != nil {
		return value, false, "failed to read"
	}

	return string(fileContent), true, ""
}

func optionalString(value string, exists bool) (string, bool, string) {
	return value, true, ""
}
