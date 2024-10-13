package cmd

import (
	"fmt"
	"net/http"
	"strings"
)

type responseShellUtilOptions struct {
	argCountMustMatch int
}

func responseShellUtilWrapper(
	commandName string,
	args []string,
	options *responseShellUtilOptions,
	f func(request *http.Request, rf *responseFiles),
) {
	if options.argCountMustMatch > 0 && len(args) != options.argCountMustMatch {
		exitWithError(fmt.Sprintf(`"%s" allows only 2 paramaters.`, commandName))
	}

	request,
		valid,
		envValidationErrors,
		rf,
		err := buildRequestFromMockEnvVars()
	if !valid {
		fmt.Println(strings.Join(envValidationErrors, "\n"))

		exitWithError(
			fmt.Sprintf("Something went wrong. \"%s\" is supposed to be used within Response Shell Scripts. Check the manual for more details.", commandName),
		)
	}
	if err != nil {
		panic(err)
	}

	f(request, rf)
}
