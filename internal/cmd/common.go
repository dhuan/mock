package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dhuan/mock/internal/record"
)

type responseShellUtilOptions struct {
	argCountMustMatch int
	argCountMax       int
}

func responseShellUtilWrapper(
	commandName string,
	args []string,
	options *responseShellUtilOptions,
	f func(request *http.Request, rf *responseFiles),
) {
	if options.argCountMustMatch > 0 && len(args) != options.argCountMustMatch {
		exitWithError(fmt.Sprintf(`"%s" allows only 2 paramaters.`, commandName))
	}

	if options.argCountMax > 0 && len(args) > options.argCountMax {
		exitWithError(fmt.Sprintf(`"%s" cannot receive more than %d parameters.`, commandName, options.argCountMax))
	}

	request,
		valid,
		envValidationErrors,
		rf,
		err := buildRequestFromMockEnvVars()
	if !valid {
		fmt.Println(strings.Join(envValidationErrors, "\n"))

		exitWithError(
			fmt.Sprintf("Something went wrong. \"%s\" is supposed to be used within Response Shell Scripts. Check the manual for more details.", commandName),
		)
	}
	if err != nil {
		panic(err)
	}

	f(request, rf)
}

func buildRequestFromMockEnvVars() (*http.Request, bool, []string, *responseFiles, error) {
	var baseApiUrl string
	var headersPlainText string
	var method string
	var endpoint string
	var querystring string
	var requestSerialized string
	var tlsStr string
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
		"MOCK_REQUEST_SERIALIZED":   {variable: &requestSerialized, f: isStringWithText},
		"MOCK_RESPONSE_HEADERS":     {variable: &responseFileHeaders, f: pointsToFile},
		"MOCK_RESPONSE_BODY":        {variable: &responseFileBody, f: pointsToFile},
		"MOCK_RESPONSE_STATUS_CODE": {variable: &responseFileStatusCode, f: pointsToFile},
	})
	if !envValid {
		return nil, false, errorMessages, nil, nil
	}

	request, err := record.UnserializeRequest(requestSerialized)
	if err != nil {
		return nil, true, []string{}, nil, err
	}

	return request, true, []string{}, &responseFiles{
		headers:    os.Getenv("MOCK_RESPONSE_HEADERS"),
		statusCode: os.Getenv("MOCK_RESPONSE_STATUS_CODE"),
		body:       os.Getenv("MOCK_RESPONSE_BODY"),
	}, nil
}
