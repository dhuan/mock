package mock

import (
	"fmt"
	"strings"

	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
)

type EndpointConfigErrorCode int

const (
	EndpointConfigErrorCode_Unknown EndpointConfigErrorCode = iota
	EndpointConfigErrorCode_EndpointDuplicate
	EndpointConfigErrorCode_FileUnreadable
	EndpointConfigErrorCode_InvalidMethod
	EndpointConfigErrorCode_RouteWithQuerystring
)

type EndpointConfigError struct {
	Code          EndpointConfigErrorCode
	EndpointIndex int
	Metadata      map[string]string
}

var available_http_methods = []string{
	"post",
	"get",
	"put",
	"patch",
	"delete",
}

func ValidateEndpointConfigs(
	endpointConfigs []types.EndpointConfig,
	readFile ReadFileFunc,
	configDirPath string,
) ([]EndpointConfigError, error) {
	endpointConfigErrors := make([]EndpointConfigError, 0)

	for i, endpointConfig := range endpointConfigs {
		newEndpointConfigErrors, err := validateEndpointConfig(&endpointConfig, i, endpointConfigs, endpointConfigErrors, readFile, configDirPath)
		if err != nil {
			return endpointConfigErrors, err
		}

		endpointConfigErrors = append(endpointConfigErrors, newEndpointConfigErrors...)
	}

	return endpointConfigErrors, nil
}

func validateEndpointConfig(
	endpointConfig *types.EndpointConfig,
	endpointConfigIndex int,
	endpointConfigs []types.EndpointConfig,
	currentEndpointConfigErrors []EndpointConfigError,
	readFile ReadFileFunc,
	configDirPath string,
) ([]EndpointConfigError, error) {
	endpointConfigErrors := make([]EndpointConfigError, 0)

	shouldSkipDuplicateFind := hasConfigErrorMatching(
		currentEndpointConfigErrors,
		EndpointConfigErrorCode_EndpointDuplicate,
		"duplicate_index",
		fmt.Sprint(endpointConfigIndex),
	)

	duplicates := []int{}
	if !shouldSkipDuplicateFind {
		duplicates = findDuplicates(endpointConfig, endpointConfigIndex, endpointConfigs)
		if len(duplicates) > 0 {
			endpointConfigErrors = append(endpointConfigErrors, EndpointConfigError{
				EndpointIndex: endpointConfigIndex,
				Code:          EndpointConfigErrorCode_EndpointDuplicate,
				Metadata: map[string]string{
					"duplicate_index": fmt.Sprint(duplicates[0]),
				},
			})
		}
	}

	if !utils.AnyEquals(available_http_methods, endpointConfig.Method) && endpointConfig.Method != "" {
		endpointConfigErrors = append(endpointConfigErrors, EndpointConfigError{
			EndpointIndex: endpointConfigIndex,
			Code:          EndpointConfigErrorCode_InvalidMethod,
			Metadata: map[string]string{
				"method": endpointConfig.Method,
			},
		})
	}

	if strings.Contains(endpointConfig.Route, "?") {
		endpointConfigErrors = append(endpointConfigErrors, EndpointConfigError{
			EndpointIndex: endpointConfigIndex,
			Code:          EndpointConfigErrorCode_RouteWithQuerystring,
			Metadata:      map[string]string{},
		})
	}

	fileReferences := getFileReferences(endpointConfig)
	if len(fileReferences) > 0 {
		newEndpointConfigErrors, err := validateFiles(fileReferences, endpointConfigIndex, readFile, configDirPath)
		if err != nil {
			return endpointConfigErrors, err
		}

		if len(newEndpointConfigErrors) > 0 {
			endpointConfigErrors = append(endpointConfigErrors, newEndpointConfigErrors...)
		}
	}

	return endpointConfigErrors, nil
}

func validateFiles(
	filePaths []string,
	endpointIndex int,
	readFile ReadFileFunc,
	configDirPath string,
) ([]EndpointConfigError, error) {
	errors := make([]EndpointConfigError, 0)

	for i := range filePaths {
		filePath := filePaths[i]

		if filePathIsDynamic(filePath) {
			continue
		}

		_, err := readFile(fmt.Sprintf("%s/%s", configDirPath, filePath))
		if err != nil {
			errors = append(errors, EndpointConfigError{
				Code:          EndpointConfigErrorCode_FileUnreadable,
				EndpointIndex: endpointIndex,
				Metadata: map[string]string{
					"file_path": filePaths[i],
				},
			})
		}
	}

	return errors, nil
}

func filePathIsDynamic(filePath string) bool {
	return strings.Contains(filePath, "${") && strings.Contains(filePath, "}")
}

func getFileReferences(endpointConfig *types.EndpointConfig) []string {
	filePaths := make([]string, 0)

	fileReference, hasFileReference := getFileReferenceFromResponseObject(endpointConfig.Response)
	if hasFileReference {
		filePaths = append(filePaths, fileReference)
	}

	return filePaths
}

func getFileReferenceFromResponseObject(response types.EndpointConfigResponse) (string, bool) {
	responseStr := string(response)
	isFileReference := utils.BeginsWith(responseStr, "file:") || utils.BeginsWith(responseStr, "sh:")

	if !isFileReference {
		return "", false
	}

	return utils.GetWord(0, utils.ReplaceRegex(
		responseStr,
		[]string{"^file:", "^sh:"},
		"",
	), ""), true
}

func hasConfigErrorMatching(
	errors []EndpointConfigError,
	errorCode EndpointConfigErrorCode,
	metadataKey,
	metadataValue string,
) bool {
	for _, configError := range errors {
		if configError.Code == errorCode && utils.MapContains(configError.Metadata, metadataKey, metadataValue) {
			return true
		}
	}

	return false
}

func findDuplicates(
	endpointConfig *types.EndpointConfig,
	endpointConfigIndex int,
	endpointConfigs []types.EndpointConfig,
) []int {
	duplicates := make([]int, 0)

	for i := range endpointConfigs {
		if i == endpointConfigIndex {
			continue
		}

		if endpointConfig.Route == endpointConfigs[i].Route && endpointConfig.Method == endpointConfigs[i].Method {
			duplicates = append(duplicates, i)
		}
	}

	return duplicates
}
