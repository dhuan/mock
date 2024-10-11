package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var writeCmd = &cobra.Command{
	Use: "write",
	Run: func(cmd *cobra.Command, args []string) {
		_,
			valid,
			envValidationErrors,
			rf,
			err := buildRequestFromMockEnvVars()
		if !valid {
			fmt.Println(strings.Join(envValidationErrors, "\n"))

			exitWithError(
				"Something went wrong. \"write\" is supposed to be used within Response Shell Scripts. Check the manual for more details.",
			)
		}
		if err != nil {
			panic(err)
		}

		stdin, err := io.ReadAll(cmd.InOrStdin())
		if err != nil {
			exitWithError("Failed to read stdin!")
		}

		err = os.WriteFile(rf.body, stdin, 0644)
		if err != nil {
			exitWithError(err.Error())
		}
	},
}
