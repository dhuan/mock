package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var getRouteParamCmd = &cobra.Command{
	Use: "get-route-param",
	Run: func(cmd *cobra.Command, args []string) {
		responseShellUtilWrapper("get-route-param", args, &responseShellUtilOptions{
			argCountMustMatch: 1,
		}, func(request *http.Request, rf *responseFiles) {
			jsonEncoded, err := base64.StdEncoding.DecodeString(os.Getenv("MOCK_ROUTE_PARAMS"))
			if err != nil {
				exitWithError(fmt.Sprintf("Failed to decode route params data: %s", err.Error()))
			}

			var routeParams map[string]interface{}
			err = json.Unmarshal(jsonEncoded, &routeParams)
			if err != nil {
				os.Exit(1)
			}

			routeParamValue, ok := routeParams[args[0]]
			if !ok {
				os.Exit(1)
			}

			_, ok = routeParamValue.(string)
			if !ok {
				exitWithError(fmt.Sprintf("Failed to convert route param value to string: %+v", routeParamValue))
			}

			fmt.Printf("%s", routeParamValue)

			os.Exit(0)
		})
	},
}
