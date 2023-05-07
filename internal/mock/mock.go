package mock

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dhuan/mock/internal/types"
	. "github.com/dhuan/mock/pkg/mock"
)

func ParseAssertRequest(req *http.Request) (*AssertConfig, error) {
	var assertConfig AssertConfig
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&assertConfig)

	return &assertConfig, err
}

func Validate(
	mockFs types.MockFs,
	assertConfig *AssertConfig,
) (*[]ValidationError, error) {
	validationErrors := make([]ValidationError, 0)
	requestRecords, err := getRequestRecordMatchingRoute(mockFs, assertConfig.Route)
	if err != nil {
		return &validationErrors, err
	}
	if len(requestRecords) == 0 {
		validationErrors = append(
			validationErrors,
			ValidationError{
				Code:     ValidationErrorCode_NoCall,
				Metadata: map[string]string{},
			},
		)

		return &validationErrors, nil
	}

	nth := assertConfig.Nth
	if nth == 0 {
		nth = 1
	}

	if (nth - 1) > (len(requestRecords) - 1) {
		validationErrors = append(
			validationErrors,
			ValidationError{
				Code:     ValidationErrorCode_NthOutOfRange,
				Metadata: map[string]string{},
			},
		)

		return &validationErrors, nil
	}

	requestRecord := requestRecords[nth-1]

	return validate(&requestRecord, assertConfig.Assert, requestRecords)
}

func validate(requestRecord *types.RequestRecord, assert *Condition, requestRecords []types.RequestRecord) (*[]ValidationError, error) {
	hasAnd := assert.And != nil
	hasOr := assert.Or != nil
	validationErrors := make([]ValidationError, 0)
	assertFunc := resolveAssertTypeFunc(assert.Type, requestRecords)
	validationErrorsCurrent, err := assertFunc(requestRecord, assert)
	success := len(validationErrorsCurrent) == 0
	if err != nil {
		return &validationErrors, err
	}

	if !success {
		validationErrors = append(validationErrors, validationErrorsCurrent...)
	}

	if success && !hasAnd {
		return &validationErrors, nil
	}

	if success && hasAnd {
		furtherValidationErrors, err := validate(requestRecord, assert.And, requestRecords)
		if err != nil {
			return &validationErrors, err
		}

		validationErrors = append(*furtherValidationErrors, validationErrorsCurrent...)
	}

	if !success && hasOr {
		furtherValidationErrors, err := validate(requestRecord, assert.Or, requestRecords)
		if err != nil {
			return &validationErrors, err
		}

		if len(*furtherValidationErrors) == 0 {
			return furtherValidationErrors, nil
		}

		validationErrors = append(*furtherValidationErrors, validationErrorsCurrent...)
	}

	return &validationErrors, nil
}

func resolveAssertTypeFunc(
	conditionType ConditionType,
	requestRecords []types.RequestRecord,
) func(requestRecord *types.RequestRecord, assert *Condition) ([]ValidationError, error) {
	if conditionType == ConditionType_HeaderMatch {
		return assertHeaderMatch
	}

	if conditionType == ConditionType_MethodMatch {
		return assertMethodMatch
	}

	if conditionType == ConditionType_JsonBodyMatch {
		return assertJsonBodyMatch
	}

	if conditionType == ConditionType_FormMatch {
		return assertFormMatch
	}

	if conditionType == ConditionType_QuerystringMatch {
		return assertQuerystringMatch
	}

	if conditionType == ConditionType_QuerystringExactMatch {
		return assertQuerystringExactMatch
	}

	if conditionType == ConditionType_Nth {
		return assertNth(requestRecords)
	}

	panic(fmt.Sprintf("Failed to resolve assert type: %d", conditionType))
}

func getRequestRecordMatchingRoute(mockFs types.MockFs, route string) ([]types.RequestRecord, error) {
	requestRecords, err := mockFs.GetRecordsMatchingRoute(route)
	if err != nil {
		return requestRecords, err
	}

	if len(requestRecords) == 0 {
		return requestRecords, err
	}

	return requestRecords, nil
}
