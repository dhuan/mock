package cmd

import (
	"net/http"

	"github.com/spf13/cobra"
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
	},
}

func buildRequestFromMockEnvVars() (*http.Request, bool, error) {
	var headersPlainText []byte
	var method string
	var endpoint string
	var querystring string

	validateEnv(map[string]foo{
		"MOCK_REQUEST_HEADERS":     pointsToFile(&headersPlainText),
		"MOCK_REQUEST_METHOD":      isStringAny(&method, validHttpMethods),
		"MOCK_REQUEST_ENDPOINT":    isStringWithText(&endpoint),
		"MOCK_REQUEST_QUERYSTRING": optionalString(&querystring),
	})

    request, err := http.NewRequest(method)
    if err != nil {
        return nil, true, err
    }

    return request, true, nil
}
