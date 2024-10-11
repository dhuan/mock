package cmd

import (
	"fmt"
	"net/http"
	"strings"
)

func responseShellUtilWrapper(
	commandName string,
	f func(request *http.Request, rf *responseFiles),
) {
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
