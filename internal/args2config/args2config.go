package args2config

import "github.com/dhuan/mock/internal/types"
import "strconv"
import "fmt"
import "log"
import "strings"

func ParseEndpoints(args []string) []types.EndpointConfig {
	endpoints := make([]types.EndpointConfig, 0)
	endpointCurrent := -1

	for i, arg := range args {
		if startingNewEndpoint(arg) {
			endpoints = append(endpoints, types.EndpointConfig{})
			endpointCurrent = endpointCurrent + 1
		}

		if endpointCurrent == -1 {
			continue
		}

		parseParam([]string{"--route"}, arg, args, i, assignRoute(&endpoints[endpointCurrent]))
		parseParam([]string{"--method"}, arg, args, i, assignMethod(&endpoints[endpointCurrent]))
		parseParam([]string{"--response"}, arg, args, i, assignResponse(&endpoints[endpointCurrent], ""))
		parseParam([]string{"--response-file"}, arg, args, i, assignResponse(&endpoints[endpointCurrent], "file:"))
		parseParam([]string{"--response-file-server", "--file-server"}, arg, args, i, assignResponse(&endpoints[endpointCurrent], "fs:"))
		parseParam([]string{"--response-sh", "--shell-script"}, arg, args, i, assignResponse(&endpoints[endpointCurrent], "sh:"))
		parseParam([]string{"--exec", "--response-exec"}, arg, args, i, assignResponse(&endpoints[endpointCurrent], "exec:"))
		parseParam([]string{"--status-code"}, arg, args, i, assignStatusCode(&endpoints[endpointCurrent]))
		parseParam([]string{"--header"}, arg, args, i, assignHeader(&endpoints[endpointCurrent]))
	}

	return endpoints
}

func ParseMiddlewares(args []string) []types.MiddlewareConfig {
	middlewares := make([]types.MiddlewareConfig, 0)
	middlewareCurrent := -1

	for i, arg := range args {
		if startingNewMiddleware(arg) {
			middlewares = append(middlewares, types.MiddlewareConfig{RouteMatch: "*"})
			middlewareCurrent = middlewareCurrent + 1
		}

		if middlewareCurrent == -1 {
			continue
		}

		parseParam([]string{"--middleware"}, arg, args, i, assignMiddlewareType(&middlewares[middlewareCurrent], types.MiddlewareType_BeforeResponse))
		parseParam([]string{"--route-match"}, arg, args, i, assignMiddlewareRouteMatch(&middlewares[middlewareCurrent]))
	}

	return middlewares
}

func startingNewEndpoint(arg string) bool {
	return arg == "--route"
}

func startingNewMiddleware(arg string) bool {
	return arg == "--middleware"
}

func assignMiddlewareType(middlewareConfig *types.MiddlewareConfig, middlewareType types.MiddlewareType) func(scriptPath string) {
	return func(scriptPath string) {
		middlewareConfig.Type = middlewareType
		middlewareConfig.Exec = scriptPath
	}
}

func assignMiddlewareRouteMatch(middlewareConfig *types.MiddlewareConfig) func(routeMatch string) {
	return func(routeMatch string) {
		middlewareConfig.RouteMatch = routeMatch
	}
}

func assignRoute(endpointConfig *types.EndpointConfig) func(route string) {
	return func(route string) {
		endpointConfig.Route = route
	}
}

func assignMethod(endpointConfig *types.EndpointConfig) func(method string) {
	return func(method string) {
		endpointConfig.Method = method
	}
}

func assignResponse(endpointConfig *types.EndpointConfig, prefix string) func(response string) {
	return func(response string) {
		endpointConfig.Response = types.EndpointConfigResponse(fmt.Sprintf("%s%s", prefix, response))
	}
}

func assignStatusCode(endpointConfig *types.EndpointConfig) func(statusCode string) {
	return func(statusCode string) {
		statusCodeParsed, err := strconv.Atoi(statusCode)
		if err != nil {
			log.Fatalf("Failed to parse %s", statusCode)

			return
		}

		endpointConfig.ResponseStatusCode = statusCodeParsed
	}
}

func assignHeader(endpointConfig *types.EndpointConfig) func(header string) {
	return func(header string) {
		headerKey, headerValue, headerOk := parseHeaderLine(header)

		if headerOk {
			if len(endpointConfig.Headers) == 0 {
				endpointConfig.Headers = map[string]string{}
			}

			endpointConfig.Headers[headerKey] = headerValue
		}
	}
}

func parseParam(
	paramNames []string,
	arg string,
	args []string,
	i int,
	f func(value string),
) {
	if !anyEquals(paramNames, arg) {
		return
	}

	if i == (len(args) - 1) {
		return
	}

	f(args[i+1])
}

func parseHeaderLine(text string) (string, string, bool) {
	splitResult := strings.Split(text, ":")

	if len(splitResult) < 2 {
		return "", "", false
	}

	return splitResult[0], strings.TrimSpace(strings.Join(splitResult[1:], ":")), true
}

func anyEquals[T comparable](list []T, value T) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}

	return false
}
