package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var getPayloadCmd = &cobra.Command{
	Use: "get-payload",
	Run: func(cmd *cobra.Command, args []string) {
		responseShellUtilWrapper("get-payload", args, &responseShellUtilOptions{
			argCountMax: 1,
		}, func(request *http.Request, rf *responseFiles) {
			fileContent, err := io.ReadAll(request.Body)
			if err != nil {
				panic(err)
			}

			if len(args) > 0 {
				fieldName := args[0]

				if searchHeader(request, "content-type", "application/json") {
					var data map[string]interface{}
					err = json.Unmarshal(fileContent, &data)

					value, ok := data[fieldName]
					if !ok {
						os.Exit(1)

						return
					}

					fmt.Printf("%s\n", value)
				}

				if searchHeader(request, "content-type", "application/x-www-form-urlencoded") {
					query, err := url.ParseQuery(string(fileContent))
					if err != nil {
						os.Exit(1)
					}

					value, ok := query[fieldName]
					if !ok {
						os.Exit(1)
					}

					fmt.Printf("%s\n", strings.Join(value, ","))

					return
				}

				return
			}

			fmt.Printf(string(fileContent))
		})
	},
}

func searchHeader(request *http.Request, key, value string) bool {
	for headerKey := range request.Header {
		headerValue := strings.Join(request.Header[headerKey], "")

		if strings.ToLower(headerKey) == key && headerValue == value {
			return true
		}
	}

	return false
}
