package record

import (
	"net/http"
	"strings"

	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
)

func BuildRequestRecord(r *http.Request, requestBody []byte) (*types.RequestRecord, error) {
	route := utils.ReplaceRegex(r.RequestURI, []string{`^\/`}, "")
	headers := buildHeadersForRequestRecord(&r.Header)
	routeParsed, querystring := parseRoute(route)
	requestRecord := &types.RequestRecord{
		Route:       routeParsed,
		Querystring: querystring,
		Headers:     *headers,
	}

	requestRecord.Body = &requestBody

	requestRecord.Method = strings.ToLower(r.Method)

	requestRecord.Host = r.Host

	https := false
	if r.TLS != nil {
		https = true
	}
	requestRecord.Https = https

	return requestRecord, nil
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
