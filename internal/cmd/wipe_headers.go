package cmd

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dhuan/mock/internal/utils"
)

var wipeHeadersCmd = &cobra.Command{
	Use: "wipe-headers",
	Run: func(cmd *cobra.Command, args []string) {
		responseShellUtilWrapper("wipe-headers", args, &responseShellUtilOptions{}, func(request *http.Request, rf *responseFiles) {
			for i := range args {
				strings.ToLower(args[i])
			}

			err := utils.MapFilterFileLines(rf.headers, func(line string) (string, bool) {
				key, _, ok := utils.ParseHeaderLine(line)
				if !ok {
					return line, true
				}

				if utils.IndexOf(args, strings.ToLower(key)) > -1 || (flagRegex && utils.IndexOfRegex(args, strings.ToLower(key)) > -1) {
					return "", false
				}

				return line, true
			})
			if err != nil {
				log.Print(err)

				os.Exit(1)
			}
		})
	},
}
