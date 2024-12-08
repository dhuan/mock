package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

				if isJsonRequest(request) {
					var data map[string]interface{}
					err = json.Unmarshal(fileContent, &data)

					value, ok := data[fieldName]
					if !ok {
						os.Exit(1)

						return
					}

					fmt.Printf("%s\n", value)
				}

				return
			}

			fmt.Printf(string(fileContent))
		})
	},
}

func isJsonRequest(request *http.Request) bool {
	for headerKey := range request.Header {
		headerValue := strings.Join(request.Header[headerKey], "")

		if strings.ToLower(headerKey) == "content-type" && headerValue == "application/json" {
			return true
		}
	}

	return false
}
