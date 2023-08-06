package record

import (
	"net/http"
	"strings"

	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
)

func BuildRequestRecord(r *http.Request, requestBody []byte, routeParams map[string]string) (*types.RequestRecord, error) {
	route := utils.ReplaceRegex(r.RequestURI, []string{`^\/`}, "")
	headers := buildHeadersForRequestRecord(&r.Header)
	routeParsed, querystring := parseRoute(route)
	requestRecord := &types.RequestRecord{
		Route:             routeParsed,
		Querystring:       querystring,
		QuerystringParsed: parseQuerystring(r),
		Headers:           *headers,
	}

	requestRecord.Body = &requestBody

	requestRecord.Method = strings.ToLower(r.Method)

	requestRecord.Host = r.Host

	https := false
	if r.TLS != nil {
		https = true
	}
	requestRecord.Https = https

	requestRecord.RouteParams = routeParams

	return requestRecord, nil
}

func parseQuerystring(r *http.Request) map[string]string {
	result := make(map[string]string)
	query := r.URL.Query()

	for key := range query {
		result[key] = query[key][0]
	}

	return result
}

func parseRoute(route string) (string, string) {
	splitResult := strings.Split(route, "?")

	if len(splitResult) == 1 {
		return splitResult[0], ""
	}

	return splitResult[0], splitResult[1]
}

func buildHeadersForRequestRecord(headers *http.Header) *http.Header {
	headersNew := make(http.Header)

	for key, value := range *headers {
		headersNew[strings.ToLower(key)] = value
	}

	return &headersNew
}
