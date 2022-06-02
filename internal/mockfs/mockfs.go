package mockfs

import (
	"encoding/json"
	"fmt"
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

func (this MockFs) StoreRequestRecord(r *http.Request, endpointConfig *types.EndpointConfig) error {
	requestRecord := buildRequestRecord(r)
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

func buildEndpointId(endpointConfig *types.EndpointConfig) string {
	return strings.ReplaceAll(endpointConfig.Route, "/", "__")
}

func (this MockFs) GetRecordsMatchingRoute(route string) ([]*types.RequestRecord, error) {
	requestRecords := make([]*types.RequestRecord, 0)

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
			requestRecords = append(requestRecords, requestRecord)
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

func buildRequestRecord(r *http.Request) *types.RequestRecord {
	route := utils.ReplaceRegex(r.RequestURI, []string{`^\/`}, "")
	headers := buildHeadersForRequestRecord(&r.Header)

	return &types.RequestRecord{
		Route:   route,
		Headers: *headers,
	}
}

func buildHeadersForRequestRecord(headers *http.Header) *http.Header {
	headersNew := make(http.Header)

	for key, value := range *headers {
		headersNew[strings.ToLower(key)] = value
	}

	return &headersNew
}

func buildRequestRecordJson(requestRecord *types.RequestRecord) ([]byte, error) {
	return json.Marshal(requestRecord)
}

func nowStr() string {
	now := time.Now()

	return fmt.Sprint(now.Unix())
}

func writeNewFile(filePath string, fileContent []byte) error {
	err := os.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		return err
	}

	return nil
}
