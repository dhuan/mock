package mockfs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
)

type MockFs struct {
	State *types.State
}

func (this MockFs) StoreRequestRecord(r *http.Request, requestBody []byte, endpointConfig *types.EndpointConfig) error {
	requestRecord, err := buildRequestRecord(r, requestBody)
	if err != nil {
		return err
	}

	requestRecordJson, err := buildRequestRecordJson(requestRecord)
	if err != nil {
		return err
	}
	requestRecordFileName := fmt.Sprintf("%s_%s", nowStr(), buildEndpointId(endpointConfig))
	requestRecordFilePath := fmt.Sprintf("%s/%s", this.State.RequestRecordDirectoryPath, requestRecordFileName)
	if err = writeNewFile(requestRecordFilePath, requestRecordJson); err != nil {
		return err
	}

	return nil
}

func (this MockFs) RemoveAllRequestRecords() error {
	walkFrom := this.State.RequestRecordDirectoryPath
	err := filepath.Walk(walkFrom, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		return os.Remove(path)
	})

	return err
}

func buildEndpointId(endpointConfig *types.EndpointConfig) string {
	return strings.ReplaceAll(endpointConfig.Route, "/", "__")
}

func (this MockFs) GetRecordsMatchingRoute(route string) ([]types.RequestRecord, error) {
	requestRecords := make([]types.RequestRecord, 0)

	walkFrom := this.State.RequestRecordDirectoryPath
	err := filepath.Walk(walkFrom, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		fileContent, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		requestRecord, err := parseRequestRecordFromFile(fileContent)
		if err != nil {
			return err
		}

		if requestRecord.Route == route {
			requestRecords = append(requestRecords, *requestRecord)
		}

		return nil
	})
	if err != nil {
		return requestRecords, err
	}

	return requestRecords, nil
}

func parseRequestRecordFromFile(requestRecordFileContent []byte) (*types.RequestRecord, error) {
	var requestRecord types.RequestRecord
	err := json.Unmarshal(requestRecordFileContent, &requestRecord)
	if err != nil {
		return &requestRecord, err
	}

	return &requestRecord, nil
}

func routeNameToRequestRecordFileRouteName(route string) string {
	return strings.ReplaceAll(route, "/", "__")
}

func buildRequestRecord(r *http.Request, requestBody []byte) (*types.RequestRecord, error) {
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

	return requestRecord, nil
}

func parseRoute(route string) (string, string) {
	splitResult := strings.Split(route, "?")

	if len(splitResult) == 1 {
		return splitResult[0], ""
	}

	return splitResult[0], splitResult[1]
}

func requestHasBody(req *http.Request) bool {
	bodyContent, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return false
	}

	return string(bodyContent) != ""
}

func buildHeadersForRequestRecord(headers *http.Header) *http.Header {
	headersNew := make(http.Header)

	for key, value := range *headers {
		headersNew[strings.ToLower(key)] = value
	}

	return &headersNew
}

func buildRequestRecordJson(requestRecord *types.RequestRecord) ([]byte, error) {
	return json.MarshalIndent(requestRecord, "", "  ")
}

func nowStr() string {
	now := time.Now()

	return fmt.Sprintf(
		"%d%d%d%d%d%d%d",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		now.Nanosecond(),
	)
}

func writeNewFile(filePath string, fileContent []byte) error {
	err := os.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		return err
	}

	return nil
}

func hasHeaderWithValue(headers *http.Header, headerKeyToSearch, headerValueToSearch string) bool {
	for headerKey, headerValues := range *headers {
		for _, headerValue := range headerValues {
			if headerKey == headerKeyToSearch && headerValue == headerValueToSearch {
				return true
			}
		}
	}

	return false
}
