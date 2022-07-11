package mock

import (
	"fmt"

	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
)

type EndpointConfigErrorCode int

const (
	EndpointConfigErrorCode_Unknown EndpointConfigErrorCode = iota
	EndpointConfigErrorCode_EndpointDuplicate
	EndpointConfigErrorCode_InvalidMethod
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

func ValidateEndpointConfigs(endpointConfigs []types.EndpointConfig) ([]EndpointConfigError, error) {
	endpointConfigErrors := make([]EndpointConfigError, 0)

	for i, endpointConfig := range endpointConfigs {
		newEndpointConfigErrors, err := validateEndpointConfig(&endpointConfig, i, endpointConfigs, endpointConfigErrors)
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

	if !utils.AnyEquals[string](available_http_methods, endpointConfig.Method) {
		endpointConfigErrors = append(endpointConfigErrors, EndpointConfigError{
			EndpointIndex: endpointConfigIndex,
			Code:          EndpointConfigErrorCode_InvalidMethod,
			Metadata: map[string]string{
				"method": endpointConfig.Method,
			},
		})
	}

	return endpointConfigErrors, nil
}

func hasConfigErrorMatching(
	errors []EndpointConfigError,
	errorCode EndpointConfigErrorCode,
	metadataKey,
	metadataValue string,
) bool {
	for _, configError := range errors {
		if configError.Code == errorCode && utils.MapContains[string, string](configError.Metadata, metadataKey, metadataValue) {
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

	for i, _ := range endpointConfigs {
		if i == endpointConfigIndex {
			continue
		}

		if endpointConfig.Route == endpointConfigs[i].Route && endpointConfig.Method == endpointConfigs[i].Method {
			duplicates = append(duplicates, i)
		}
	}

	return duplicates
}
