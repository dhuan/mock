package cmd

import (
	"fmt"
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

			contentToWrite := stdin

			if flagAppend {
				fileContent, err := os.ReadFile(os.Getenv("MOCK_RESPONSE_BODY"))
				if err != nil {
					fmt.Printf("Failed to read response body file.\n")

					os.Exit(1)
				}

				contentToWrite = append(fileContent, stdin...)
			}

			err = os.WriteFile(rf.body, contentToWrite, 0644)
			if err != nil {
				exitWithError(err.Error())
			}
		})
	},
}
