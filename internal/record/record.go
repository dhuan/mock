package record

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
)

func BuildRequestRecord(r *http.Request, routeParams map[string]string) (*types.RequestRecord, []byte, error) {
	r2, err := utils.CloneRequest(r)
	if err != nil {
		return nil, nil, err
	}

	r3, err := utils.CloneRequest(r)
	if err != nil {
		return nil, nil, err
	}

	route := utils.ReplaceRegex(r.RequestURI, []string{`^\/`}, "")
	headers := buildHeadersForRequestRecord(&r.Header)
	routeParsed, querystring := parseRoute(route)
	requestRecord := &types.RequestRecord{
		Route:             routeParsed,
		Querystring:       querystring,
		QuerystringParsed: parseQuerystring(r),
		Headers:           *headers,
	}

	requestBody, err := io.ReadAll(r2.Body)
	if err != nil {
		return nil, nil, err
	}

	requestRecord.Body = &requestBody

	requestRecord.Method = strings.ToLower(r.Method)

	requestRecord.Host = r.Host

	requestSerialized, err := serializeRequest(r3)
	if err != nil {
		return nil, nil, err
	}

	requestRecord.Serialized = requestSerialized

	https := false
	if r.TLS != nil {
		https = true
	}
	requestRecord.Https = https

	requestRecord.RouteParams = routeParams

	return requestRecord, requestBody, nil
}

func UnserializeRequest(serialized string) (*http.Request, error) {
	decoded, err := base64.StdEncoding.DecodeString(serialized)
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(bytes.NewReader(decoded))

	return http.ReadRequest(r)
}

func serializeRequest(r *http.Request) (string, error) {
	b := bytes.Buffer{}
	if err := r.Write(&b); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
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
