package args2config

import "github.com/dhuan/mock/internal/types"
import "strconv"

func Parse(args []string) []types.EndpointConfig {
	endpoints := make([]types.EndpointConfig, 0)
	endpointCurrent := -1

	for i, arg := range args {
		if startingNew(arg) {
			endpoints = append(endpoints, types.EndpointConfig{})
			endpointCurrent = endpointCurrent + 1
		}

		routeName, isRoute := parseRoute(arg, args, i)
		if isRoute {
			endpoints[endpointCurrent].Route = routeName
		}

		method, isMethod := parseMethod(arg, args, i)
		if isMethod {
			endpoints[endpointCurrent].Method = method
		}

		response, isResponse := parseResponse(arg, args, i)
		if isResponse {
			endpoints[endpointCurrent].Response = response
		}

		statusCode, isStatusCode := parseStatusCode(arg, args, i)
		if isStatusCode {
			endpoints[endpointCurrent].ResponseStatusCode = statusCode
		}
	}

	return endpoints
}

func startingNew(arg string) bool {
	return arg == "--route"
}

func parseRoute(arg string, args []string, i int) (string, bool) {
	if arg != "--route" {
		return "", false
	}

	if i == (len(args) - 1) {
		return "", false
	}

	return args[i+1], true
}

func parseMethod(arg string, args []string, i int) (string, bool) {
	if arg != "--method" {
		return "", false
	}

	if i == (len(args) - 1) {
		return "", false
	}

	return args[i+1], true
}

func parseStatusCode(arg string, args []string, i int) (int, bool) {
	if arg != "--status-code" {
		return 0, false
	}

	if i == (len(args) - 1) {
		return 0, false
	}

	statusCode, err := strconv.Atoi(args[i+1])
	if err != nil {
		return 0, false
	}

	return statusCode, true
}

func parseResponse(arg string, args []string, i int) ([]byte, bool) {
	if arg != "--response" {
		return []byte(""), false
	}

	if i == (len(args) - 1) {
		return []byte(""), false
	}

	return []byte(args[i+1]), true
}
