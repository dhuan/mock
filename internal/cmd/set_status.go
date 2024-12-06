package cmd

import (
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var setStatusCmd = &cobra.Command{
	Use: "set-status",
	Run: func(cmd *cobra.Command, args []string) {
		responseShellUtilWrapper("set-status", args, &responseShellUtilOptions{
			argCountMustMatch: 1,
		}, func(request *http.Request, rf *responseFiles) {
			statusCode := args[0]

			err := os.WriteFile(rf.statusCode, []byte(statusCode), 0644)
			if err != nil {
				exitWithError(err.Error())
			}
		})
	},
}
