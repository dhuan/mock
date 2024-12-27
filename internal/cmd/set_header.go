package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dhuan/mock/internal/utils"
)

var setHeaderCmd = &cobra.Command{
	Use: "set-header",
	Run: func(cmd *cobra.Command, args []string) {
		responseShellUtilWrapper("set-header", args, &responseShellUtilOptions{
			argCountMustMatch: 2,
		}, func(request *http.Request, rf *responseFiles) {
			headerKey := args[0]
			headerValue := args[1]
			headerExists := false
			newHeaderLine := fmt.Sprintf("%s: %s", headerKey, headerValue)

			err := utils.MapFilterFileLines(rf.headers, func(line string) (string, bool) {
				key, _, ok := utils.ParseHeaderLine(line)
				if !ok {
					return line, true
				}

				if strings.EqualFold(key, headerKey) {
					headerExists = true

					return newHeaderLine, true
				}

				return line, true
			})

			if headerExists {
				os.Exit(0)
			}

			if err = utils.AddLineToFile(rf.headers, newHeaderLine); err != nil {
				log.Print(err)

				os.Exit(1)
			}

			os.Exit(0)
		})
	},
}
