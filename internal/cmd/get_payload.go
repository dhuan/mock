package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/dhuan/mock/internal/map_extract"
	"github.com/dhuan/mock/internal/utils"

	"github.com/spf13/cobra"
)

var getPayloadCmd = &cobra.Command{
	Use: "get-payload",
	Run: func(cmd *cobra.Command, args []string) {
		responseShellUtilWrapper("get-payload", args, &responseShellUtilOptions{
			argCountMax: 1,
		}, func(request *http.Request, rf *responseFiles) {
			requestCloned, err := utils.CloneRequest(request)
			if err != nil {
				panic(err)
			}

			fileContent, err := io.ReadAll(request.Body)
			if err != nil {
				panic(err)
			}

			if len(args) > 0 {
				fieldName := args[0]

				getField, ok := resolveGetFieldFunc(requestCloned)
				if !ok {
					return
				}

				value, ok := getField(requestCloned, fileContent, fieldName)
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

func getPayloadField_Json(request *http.Request, payload []byte, fieldName string) (string, bool) {
	var data interface{}
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return "", false
	}

	nestedQuery := strings.Contains(fieldName, ".") || strings.Contains(fieldName, "[")

	payloadIsObject := reflect.TypeOf(data).Kind() == reflect.Map

	if !nestedQuery && payloadIsObject {
		dataParsed, ok := data.(map[string]interface{})
		if !ok {
			return "", false
		}

		value, ok := dataParsed[fieldName]
		if !ok {
			return "", false
		}

		return fmt.Sprintf("%+v", value), true
	}

	value, ok := map_extract.Extract(data, fieldName)

	return fmt.Sprintf("%+v", value), ok
}

func getPayloadField_UrlEncoded(request *http.Request, payload []byte, fieldName string) (string, bool) {
	query, err := url.ParseQuery(string(payload))
	if err != nil {
		return "", false
	}

	value, ok := query[fieldName]
	if !ok {
		return "", false
	}

	return strings.Join(value, ","), true
}

func getPayloadField_Multipart(request *http.Request, payload []byte, fieldName string) (string, bool) {
	err := request.ParseMultipartForm(1024 * 1024 * 16)
	if err != nil {
		panic(err)
	}

	value, ok := request.MultipartForm.Value[fieldName]
	if !ok {
		value, ok := request.MultipartForm.File[fieldName]
		if !ok {
			return "", false
		}

		if len(value) == 0 {
			return "", false
		}

		file, err := value[0].Open()
		if err != nil {
			return "", false
		}

		fileContent, err := io.ReadAll(file)
		if err != nil {
			return "", false
		}

		return string(fileContent), true
	}

	return strings.Join(value, ","), true
}

func resolveGetFieldFunc(request *http.Request) (func(request *http.Request, payload []byte, fieldName string) (string, bool), bool) {
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

	if utils.RegexTest("^multipart/form-data", contentType) {
		return getPayloadField_Multipart, true
	}

	return getPayloadField_Json, false
}
