package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/dhuan/mock/internal/utils"
)

var getQueryCmd = &cobra.Command{
	Use: "get-query",
	Run: func(cmd *cobra.Command, args []string) {
		responseShellUtilWrapper("get-query", args, &responseShellUtilOptions{
			argCountMax: 1,
		}, func(request *http.Request, rf *responseFiles) {
			querystringSerialized := os.Getenv("MOCK_REQUEST_QUERYSTRING_SERIALIZED")
			if querystringSerialized == "" {
				os.Exit(1)

				return
			}

			querystringParsed, err := utils.DecodeBase64Json(querystringSerialized)
			if len(querystringParsed) == 0 {
				fmt.Println(err)

				os.Exit(1)

				return
			}

			if len(args) == 0 {
				fmt.Printf("%s", toQuerystring(querystringParsed))

				os.Exit(0)
			}

			value, ok := querystringParsed[args[0]]
			if !ok {
				os.Exit(1)
			}

			fmt.Printf("%s", value)
		})
	},
}

func toQuerystring(data map[string]interface{}) string {
	result := ""

	for _, key := range utils.GetSortedKeys(data) {
		result = fmt.Sprintf("%s&%s=%s", result, key, data[key])
	}

	return result[1:]
}
