package args2config

import "github.com/dhuan/mock/internal/types"
import "strconv"
import "fmt"
import "log"
import "strings"

func Parse(args []string) []types.EndpointConfig {
	endpoints := make([]types.EndpointConfig, 0)
	endpointCurrent := -1

	for i, arg := range args {
		if startingNew(arg) {
			endpoints = append(endpoints, types.EndpointConfig{})
			endpointCurrent = endpointCurrent + 1
		}

		routeName, isRoute := parseParamString("--route", arg, args, i)
		if isRoute {
			endpoints[endpointCurrent].Route = routeName
		}

		method, isMethod := parseParamString("--method", arg, args, i)
		if isMethod {
			endpoints[endpointCurrent].Method = method
		}

		response, isResponse := parseParamString("--response", arg, args, i)
		if isResponse {
			endpoints[endpointCurrent].Response = types.EndpointConfigResponse(response)
		}

		responseFile, isResponseFile := parseParamString("--response-file", arg, args, i)
		if isResponseFile {
			endpoints[endpointCurrent].Response = types.EndpointConfigResponse(fmt.Sprintf("file:%s", responseFile))
		}

		responseFileServer, isResponseFileServer := parseParamString("--response-file-server", arg, args, i)
		if isResponseFileServer {
			endpoints[endpointCurrent].Response = types.EndpointConfigResponse(fmt.Sprintf("fs:%s", responseFileServer))
		}

		header, isHeader := parseParamString("--header", arg, args, i)
		if isHeader {
			headerKey, headerValue, headerOk := parseHeaderLine(header)

			if headerOk {
				if len(endpoints[endpointCurrent].Headers) == 0 {
					endpoints[endpointCurrent].Headers = map[string]string{}
				}

				endpoints[endpointCurrent].Headers[headerKey] = headerValue
			}
		}

		statusCode, isStatusCode := parseParamString("--status-code", arg, args, i)
		if isStatusCode {
			statusCodeParsed, err := strconv.Atoi(statusCode)
			if err != nil {
				log.Fatalln(fmt.Sprintf("Failed to parse %s", statusCode))

				return endpoints
			}

			endpoints[endpointCurrent].ResponseStatusCode = statusCodeParsed
		}
	}

	return endpoints
}

func startingNew(arg string) bool {
	return arg == "--route"
}

func parseParamString(paramName string, arg string, args []string, i int) (string, bool) {
	if arg != paramName {
		return "", false
	}

	if i == (len(args) - 1) {
		return "", false
	}

	return args[i+1], true
}

func parseHeaderLine(text string) (string, string, bool) {
	splitResult := strings.Split(text, ":")

	if len(splitResult) < 2 {
		return "", "", false
	}

	return splitResult[0], strings.TrimSpace(strings.Join(splitResult[1:], ":")), true
}
