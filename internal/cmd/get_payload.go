package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var getPayloadCmd = &cobra.Command{
	Use: "get-payload",
	Run: func(cmd *cobra.Command, args []string) {
		responseShellUtilWrapper("get-payload", args, &responseShellUtilOptions{
			argCountMax: 1,
		}, func(request *http.Request, rf *responseFiles) {
			fileContent, err := io.ReadAll(request.Body)
			if err != nil {
				panic(err)
			}

			if len(args) > 0 {
				fieldName := args[0]

				getField, ok := resolveGetFieldFunc(request)
				if !ok {
					return
				}

				value, ok := getField(fileContent, fieldName)
				if !ok {
					os.Exit(1)
				}

				fmt.Printf("%s\n", value)

				return
			}

			fmt.Printf(string(fileContent))
		})
	},
}

func getHeader(request *http.Request, key string) (string, bool) {
	for headerKey := range request.Header {
		headerValue := strings.Join(request.Header[headerKey], "")

		if strings.ToLower(headerKey) == key {
			return headerValue, true
		}
	}

	return "", false
}

func getPayloadField_Json(payload []byte, fieldName string) (string, bool) {
	var data map[string]interface{}
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return "", false
	}

	value, ok := data[fieldName]
	if !ok {
		return "", false
	}

	return fmt.Sprintf("%+v", value), true
}

func getPayloadField_UrlEncoded(payload []byte, fieldName string) (string, bool) {
	query, err := url.ParseQuery(string(payload))
	if err != nil {
		return "", false
	}

	value, ok := query[fieldName]
	if !ok {
		return "", false
	}

	return fmt.Sprintf("%s", strings.Join(value, ",")), true
}

func resolveGetFieldFunc(request *http.Request) (func(payload []byte, fieldName string) (string, bool), bool) {
	contentType, ok := getHeader(request, "content-type")
	if !ok {
		return getPayloadField_Json, false
	}

	if contentType == "application/json" {
		return getPayloadField_Json, true
	}

	if contentType == "application/x-www-form-urlencoded" {
		return getPayloadField_UrlEncoded, true
	}

	return getPayloadField_Json, false
}
