package mock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
	. "github.com/dhuan/mock/pkg/mock"
)

func ParseAssertRequest(req *http.Request) (*AssertOptions, error) {
	var assertOptions AssertOptions
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&assertOptions)

	return &assertOptions, err
}

func Validate(
	mockFs types.MockFs,
	assertOptions *AssertOptions,
) (*[]ValidationError, error) {
	validationErrors := make([]ValidationError, 0)
	requestRecords, err := getRequestRecordMatchingRoute(mockFs, assertOptions.Route)
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

	nth := assertOptions.Nth
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

	return validate(&requestRecord, assertOptions.Condition, requestRecords)
}

func validate(requestRecord *types.RequestRecord, assert *Condition, requestRecords []types.RequestRecord) (*[]ValidationError, error) {
	hasAnd := assert.And != nil
	hasOr := assert.Or != nil
	validationErrors := make([]ValidationError, 0)
	assertFunc := resolveAssertTypeFunc(assert.Type, requestRecords)
	validationErrorsCurrent, err := assertFunc(requestRecord, requestRecords, assert)
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

type asserterFunc = func(requestRecord *types.RequestRecord, requestRecords []types.RequestRecord, assert *Condition) ([]ValidationError, error)

var asserters_map map[ConditionType]asserterFunc = map[ConditionType]asserterFunc{
	ConditionType_HeaderMatch:                assertHeaderMatch,
	ConditionType_MethodMatch:                assertMethodMatch,
	ConditionType_JsonBodyMatch:              assertJsonBodyMatch,
	ConditionType_FormMatch:                  assertFormMatch,
	ConditionType_QuerystringMatch:           assertQuerystringMatch,
	ConditionType_QuerystringMatchRegex:      assertQuerystringMatchRegex,
	ConditionType_QuerystringExactMatch:      assertQuerystringExactMatch,
	ConditionType_QuerystringExactMatchRegex: assertQuerystringExactMatchRegex,
	ConditionType_Nth:                        assertNth,
	ConditionType_RouteParamMatch:            assertRouteParamMatch,
}

func resolveAssertTypeFunc(
	conditionType ConditionType,
	requestRecords []types.RequestRecord,
) asserterFunc {
	assert, ok := asserters_map[conditionType]
	if !ok {
		panic(fmt.Sprintf("Failed to resolve assert type: %d", conditionType))
	}

	return assert
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

func BuildVars(
	state *types.State,
	responseStatusCode int,
	requestRecord *types.RequestRecord,
	requestRecords []types.RequestRecord,
	requestBody []byte,
	baseApi string,
) (map[string]string, error) {
	endpoint := requestRecord.Route
	mockHost := fmt.Sprintf("localhost:%s", state.ListenPort)
	querystring := requestRecord.Querystring
	protocol := "http://"
	if requestRecord.Https {
		protocol = "https://"
	}

	nth := 1
	for i := range requestRecords {
		if requestRecords[i].Route == requestRecord.Route && requestRecords[i].Method == requestRecord.Method {
			nth = nth + 1
		}
	}

	result := map[string]string{
		"MOCK_HOST": mockHost, "MOCK_REQUEST_HOST": requestRecord.Host,
		"MOCK_REQUEST_URL":         fmt.Sprintf("%s%s/%s", protocol, requestRecord.Host, requestRecord.Route),
		"MOCK_REQUEST_ENDPOINT":    endpoint,
		"MOCK_REQUEST_METHOD":      requestRecord.Method,
		"MOCK_REQUEST_QUERYSTRING": querystring,
		"MOCK_REQUEST_NTH":         fmt.Sprintf("%d", nth),
		"MOCK_BASE_API":            baseApi,
	}

	addQuerystringParams(result, requestRecord)

	addHeaders(result, requestRecord)

	return result, nil
}

func addQuerystringParams(m map[string]string, requestRecord *types.RequestRecord) {
	for key := range requestRecord.QuerystringParsed {
		m[fmt.Sprintf("MOCK_REQUEST_QUERYSTRING_%s", strings.ToUpper(key))] = requestRecord.QuerystringParsed[key]
	}
}

func addHeaders(m map[string]string, requestRecord *types.RequestRecord) {
	for key := range requestRecord.Headers {
		headerKey := utils.ReplaceRegex(key, []string{"-"}, "_")

		m[fmt.Sprintf("MOCK_REQUEST_HEADER_%s", strings.ToUpper(headerKey))] = requestRecord.Headers[key][0]
	}
}

func BuildResponseVars(response *http.Response) (map[string]string, error) {
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	responseBodyFile, err := utils.CreateTempFile([]byte(responseBody))
	if err != nil {
		return nil, err
	}

	responseStatusCodeFile, err := utils.CreateTempFile([]byte(fmt.Sprintf("%d", response.StatusCode)))
	if err != nil {
		return nil, err
	}

	responseHeadersFile, err := utils.CreateTempFile([]byte(utils.ToHeadersText(response.Header) + "\n"))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"MOCK_RESPONSE_BODY":        responseBodyFile,
		"MOCK_RESPONSE_STATUS_CODE": responseStatusCodeFile,
		"MOCK_RESPONSE_HEADERS":     responseHeadersFile,
	}, nil
}
