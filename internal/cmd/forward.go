package cmd

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
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
		responseShellUtilWrapper("forward", args, &responseShellUtilOptions{}, func(request *http.Request, rf *responseFiles) {
			log.Printf("Forwarding request to Base API: %s %s\n", request.Method, request.RequestURI)

			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				panic(err)
			}

			log.Printf("Got response from Base API: %d\n", response.StatusCode)

			responseBody, err := io.ReadAll(response.Body)
			if err != nil {
				panic(err)
			}

			responseBody, err = uncompressIfNeeded(response, responseBody)
			if err != nil {
				panic(err)
			}

			if err = writeFile(rf.body, responseBody); err != nil {
				panic(err)
			}

			if err = writeFile(rf.statusCode, []byte(fmt.Sprintf("%d", response.StatusCode))); err != nil {
				panic(err)
			}

			for key := range response.Header {
				if strings.ToLower(key) == "content-length" {
					response.Header.Del(key)
				}
			}

			if err = writeFile(rf.headers, []byte(utils.ToHeadersText(response.Header))); err != nil {
				panic(err)
			}
		})
	},
}

func uncompressIfNeeded(response *http.Response, responseBody []byte) ([]byte, error) {
	encoding := response.Header.Get("Content-Encoding")
	if encoding != "gzip" {
		return responseBody, nil
	}

	responseBodyReader := bytes.NewReader(responseBody)

	uncompressReader, err := gzip.NewReader(responseBodyReader)
	if err != nil {
		return nil, err
	}

	response.Header.Del("Content-Encoding")

	return io.ReadAll(uncompressReader)
}

func writeFile(filePath string, data []byte) error {
	return os.WriteFile(filePath, data, 0644)
}

var validHttpMethods []string = []string{
	"get",
	"post",
	"delete",
	"patch",
	"put",
	"options",
}

type responseFiles struct {
	headers    string
	statusCode string
	body       string
}

func buildRequestFromMockEnvVars() (*http.Request, bool, []string, *responseFiles, error) {
	var baseApiUrl string
	var headersPlainText string
	var method string
	var endpoint string
	var querystring string
	var tlsStr string
	var tls bool
	var responseFileHeaders string
	var responseFileBody string
	var responseFileStatusCode string

	envValid, errorMessages := validateEnv(map[string]*envValidationConfig{
		"MOCK_BASE_API":             {variable: &baseApiUrl, f: optionalString},
		"MOCK_REQUEST_HEADERS":      {variable: &headersPlainText, f: pointsToFile},
		"MOCK_REQUEST_METHOD":       {variable: &method, f: isStringAny(validHttpMethods)},
		"MOCK_REQUEST_ENDPOINT":     {variable: &endpoint, f: isStringWithText},
		"MOCK_REQUEST_QUERYSTRING":  {variable: &querystring, f: optionalString},
		"MOCK_REQUEST_HTTPS":        {variable: &tlsStr, f: isBoolString},
		"MOCK_RESPONSE_HEADERS":     {variable: &responseFileHeaders, f: pointsToFile},
		"MOCK_RESPONSE_BODY":        {variable: &responseFileBody, f: pointsToFile},
		"MOCK_RESPONSE_STATUS_CODE": {variable: &responseFileStatusCode, f: pointsToFile},
	})
	if !envValid {
		return nil, false, errorMessages, nil, nil
	}

	tls = tlsStr == "true"

	method = strings.ToUpper(method)
	url := fmt.Sprintf("%s/%s", baseApiUrl, endpoint)
	if querystring != "" {
		url = fmt.Sprintf("%s?%s", url, querystring)
	}

	if !utils.RegexTest("^https?://", url) {
		protocol := "http://"
		if tls {
			protocol = "https://"
		}

		url = fmt.Sprintf("%s%s", protocol, url)
	}

	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, true, []string{}, nil, err
	}

	headers, err := mock.ExtractHeadersFromFile(os.Getenv("MOCK_REQUEST_HEADERS"), readFile)
	if err != nil {
		return nil, true, []string{}, nil, err
	}
	for headerKey, headerValue := range headers {
		request.Header.Add(headerKey, headerValue)
	}

	return request, true, []string{}, &responseFiles{
		headers:    os.Getenv("MOCK_RESPONSE_HEADERS"),
		statusCode: os.Getenv("MOCK_RESPONSE_STATUS_CODE"),
		body:       os.Getenv("MOCK_RESPONSE_BODY"),
	}, nil
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

func isBoolString(value string, exists bool) (string, bool, string) {
	return strings.ToLower(value),
		exists && utils.AnyEquals([]string{"true", "false"}, strings.ToLower(value)),
		"is not a string with text"
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
