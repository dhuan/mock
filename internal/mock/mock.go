package mock

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dhuan/mock/internal/types"
	. "github.com/dhuan/mock/pkg/mock"
)

type JsonValidate func(jsonA map[string]interface{}, jsonB map[string]interface{}) bool

func ParseAssertRequest(req *http.Request) (*AssertConfig, error) {
	var assertConfig AssertConfig
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&assertConfig)

	return &assertConfig, err
}

func Validate(
	mockFs types.MockFs,
	jsonValidate JsonValidate,
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

	return validate(requestRecord, assertConfig.Assert, jsonValidate)
}

func validate(requestRecord *types.RequestRecord, assert *AssertOptions, jsonValidate JsonValidate) (*[]ValidationError, error) {
	hasAnd := assert.And != nil
	hasOr := assert.Or != nil
	validationErrors := make([]ValidationError, 0)
	assertFunc := resolveAssertTypeFunc(assert.Type, jsonValidate)
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
		furtherValidationErrors, err := validate(requestRecord, assert.And, jsonValidate)
		if err != nil {
			return &validationErrors, err
		}

		validationErrors = append(*furtherValidationErrors, validationErrorsCurrent...)
	}

	if !success && hasOr {
		furtherValidationErrors, err := validate(requestRecord, assert.Or, jsonValidate)
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
	assertType AssertType,
	jsonValidate JsonValidate,
) func(requestRecord *types.RequestRecord, assert *AssertOptions) ([]ValidationError, error) {
	if assertType == AssertType_HeaderMatch {
		return assertHeaderMatch
	}

	if assertType == AssertType_MethodMatch {
		return assertMethodMatch
	}

	if assertType == AssertType_JsonBodyMatch {
		return assertJsonBodyMatch(jsonValidate)
	}

	if assertType == AssertType_FormMatch {
		return assertFormMatch
	}

	if assertType == AssertType_QuerystringMatch {
		return assertQuerystringMatch
	}

	if assertType == AssertType_QuerystringExactMatch {
		return assertQuerystringExactMatch
	}

	panic(fmt.Sprintf("Failed to resolve assert type: %d", assertType))
}

func getRequestRecordMatchingRoute(mockFs types.MockFs, route string) ([]*types.RequestRecord, error) {
	requestRecords, err := mockFs.GetRecordsMatchingRoute(route)
	if err != nil {
		return requestRecords, err
	}

	if len(requestRecords) == 0 {
		return requestRecords, err
	}

	return requestRecords, nil
}
