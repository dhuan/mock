package cmd

import (
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var replaceCmd = &cobra.Command{
	Use: "replace",
	Run: func(cmd *cobra.Command, args []string) {
		responseShellUtilWrapper("write", func(request *http.Request, rf *responseFiles) {
			fileContent, err := os.ReadFile(rf.body)
			if err != nil {
				panic(err)
			}

			result := strings.Replace(string(fileContent), args[0], args[1], 1)

			err = os.WriteFile(rf.body, []byte(result), 0644)
			if err != nil {
				panic(err)
			}
		})
	},
}
