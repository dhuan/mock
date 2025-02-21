package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/dhuan/mock/internal/utils"

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

			if flagJson {
				if flagAppend {
					fmt.Printf("--json cannot be used with --append.\n")

					os.Exit(1)
				}

				jsonStr, err := formatJson(contentToWrite)
				if err != nil || string(jsonStr) == "null" {
					os.Exit(1)
				}

				if err = utils.AddLineToFile(rf.headers, "Content-Type: application/json"); err != nil {
					log.Print(err)

					os.Exit(1)
				}

				contentToWrite = jsonStr
			}

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

func formatJson(data []byte) ([]byte, error) {
	var result interface{}
	json.Unmarshal(data, &result)

	return json.Marshal(result)
}
