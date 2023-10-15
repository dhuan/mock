package mockfs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dhuan/mock/internal/types"
)

type MockFs struct {
	State *types.State
}

func (mfs MockFs) StoreRequestRecord(requestRecord *types.RequestRecord, endpointConfig *types.EndpointConfig) error {
	requestRecordJson, err := buildRequestRecordJson(requestRecord)
	if err != nil {
		return err
	}
	requestRecordFileName := fmt.Sprintf("%s_%s", nowStr(), buildEndpointId(endpointConfig))
	requestRecordFilePath := fmt.Sprintf("%s/%s", mfs.State.RequestRecordDirectoryPath, requestRecordFileName)
	if err = writeNewFile(requestRecordFilePath, requestRecordJson); err != nil {
		return err
	}

	return nil
}

func (mfs MockFs) RemoveAllRequestRecords() error {
	walkFrom := mfs.State.RequestRecordDirectoryPath
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

func (mfs MockFs) GetRecordsMatchingRoute(route string) ([]types.RequestRecord, error) {
	requestRecords := make([]types.RequestRecord, 0)

	walkFrom := mfs.State.RequestRecordDirectoryPath
	err := filepath.Walk(walkFrom, func(path string, info os.FileInfo, err2 error) error {
		if info.IsDir() {
			return nil
		}

		fileContent, err3 := os.ReadFile(path)
		if err3 != nil {
			return err3
		}

		requestRecord, err3 := parseRequestRecordFromFile(fileContent)
		if err3 != nil {
			return err3
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
