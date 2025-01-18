package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dhuan/mock/internal/utils"
)

var getHeaderCmd = &cobra.Command{
	Use: "get-header",
	Run: func(cmd *cobra.Command, args []string) {
		responseShellUtilWrapper("get-header", args, &responseShellUtilOptions{
			argCountMax: 1,
		}, func(request *http.Request, rf *responseFiles) {
			headersFileContent, err := os.ReadFile(os.Getenv("MOCK_REQUEST_HEADERS"))
			if err != nil {
				panic(err)
			}

			if len(args) == 0 {
				fmt.Println(string(headersFileContent))

				os.Exit(0)

				return
			}

			headerSearch := args[0]

			headers := utils.ExtractHeadersFromText(headersFileContent)
			headerKeys := utils.GetSortedKeys(headers)
			matches := make([]int, 0)

			for i, headerKey := range headerKeys {
				if flagRegex {
					if utils.RegexTest(strings.ToLower(headerSearch), strings.ToLower(headerKey)) {
						matches = append(matches, i)
					}

					continue
				}

				if headerSearch == strings.ToLower(headerKey) {
					matches = append(matches, i)
				}
			}

			if len(matches) == 0 {
				os.Exit(1)

				return
			}

			for _, headerIndex := range matches {
				if headerIndex > (len(headerKeys) - 1) {
					panic(fmt.Errorf("Something went wrong while searching for headers."))
				}

				headerKey := headerKeys[headerIndex]

				headerValue, ok := headers[headerKey]
				if !ok {
					panic(fmt.Errorf("Something went wrong while searching for headers."))
				}

				headerValue = strings.TrimSpace(headerValue)

				if flagValueOnly {
					fmt.Printf("%s\n", headerValue)
				} else {
					fmt.Printf("%s: %s\n", headerKey, headerValue)
				}
			}
		})
	},
}
