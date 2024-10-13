package cmd

import (
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var writeCmd = &cobra.Command{
	Use: "write",
	Run: func(cmd *cobra.Command, args []string) {
		responseShellUtilWrapper("write", args, &responseShellUtilOptions{}, func(request *http.Request, rf *responseFiles) {
			stdin, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				exitWithError("Failed to read stdin!")
			}

			err = os.WriteFile(rf.body, stdin, 0644)
			if err != nil {
				exitWithError(err.Error())
			}
		})
	},
}
