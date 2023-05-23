package mock

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strings"

	"github.com/dhuan/mock/internal/types"
	"github.com/dhuan/mock/internal/utils"
	. "github.com/dhuan/mock/pkg/mock"
)

func assertHeaderMatch(requestRecord *types.RequestRecord, requestRecords []types.RequestRecord, assert *Condition) ([]ValidationError, error) {
	validationErrors := make([]ValidationError, 0)
	keyValues := getKeyValuePairsFromAssertionOptions(assert)

	for _, key := range utils.GetSortedKeys(keyValues) {
		value := keyValues[key]

		valueFromRequestRecord, ok := requestRecord.Headers[key]
		if !ok {
			validationErrors = append(validationErrors, ValidationError{
				Code: ValidationErrorCode_HeaderNotIncluded,
				Metadata: map[string]string{
					"missing_header_key": key,
				},
			})

			continue
		}

		if value != strings.Join(valueFromRequestRecord, "") {
			validationErrors = append(validationErrors, ValidationError{
				Code: ValidationErrorCode_HeaderValueMismatch,
				Metadata: map[string]string{
					"header_key":             key,
					"header_value_requested": strings.Join(valueFromRequestRecord, ""),
					"header_value_expected":  value.(string),
				},
			})
		}
	}

	return validationErrors, nil
}

func filterRequestRecordsMatchingRouteAndMethod(
	requestRecords []types.RequestRecord,
	route,
	method string,
) []types.RequestRecord {
	result := make([]types.RequestRecord, 0)

	for i := range requestRecords {
		if requestRecords[i].Route == route && requestRecords[i].Method == method {
			result = append(result, requestRecords[i])
		}
	}

	return result
}

func assertNth(requestRecord *types.RequestRecord, requestRecords []types.RequestRecord, assert *Condition) ([]ValidationError, error) {
	filteredRequestRecords := filterRequestRecordsMatchingRouteAndMethod(requestRecords, requestRecord.Route, requestRecord.Method)
	currentRequestNthNumber := len(filteredRequestRecords) + 1
	currentRequestNth := fmt.Sprint(currentRequestNthNumber)

	if utils.EndsWith(string(assert.Value), "+") {
		return assertNthWithPlus(filteredRequestRecords, currentRequestNthNumber, assert)
	}

	if string(currentRequestNth) != string(assert.Value) {
		return []ValidationError{
			{Code: ValidationErrorCode_NthMismatch, Metadata: map[string]string{
				"nth_requested": string(currentRequestNth),
				"nth_expected":  string(assert.Value),
			}},
		}, nil
	}

	return []ValidationError{}, nil

}

func assertRouteParamMatch(requestRecord *types.RequestRecord, requestRecords []types.RequestRecord, assert *Condition) ([]ValidationError, error) {
	validationErrors := make([]ValidationError, 0)
	keyValues := make(map[string]string)
	hasKeyAndValue := assert.Key != "" && assert.Value != ""

	if hasKeyAndValue {
		keyValues[assert.Key] = string(assert.Value)
	}

	for key := range assert.KeyValues {
		keyValues[key] = assert.KeyValues[key].(string)
	}

	for key := range keyValues {
		exists, equals, value := utils.MapContainsX(requestRecord.RouteParams, key, keyValues[key], "")
		fmt.Println("!!!!!!!!!!!!!")
		fmt.Println(exists)
		fmt.Println(equals)
		fmt.Println(keyValues[key])
		fmt.Println(value)

		if !exists {
			validationErrors = append(
				validationErrors,
				ValidationError{
					Code: ValidationErrorCode_RouteParamDoesNotExistInEndpoint,
					Metadata: map[string]string{
						"route_param_key": key,
					},
				},
			)
		}

		if exists && !equals {
			validationErrors = append(
				validationErrors,
				ValidationError{
					Code: ValidationErrorCode_RouteParamValueMismatch,
					Metadata: map[string]string{
						"route_param_key":             key,
						"route_param_value_expected":  keyValues[key],
						"route_param_value_requested": value,
					},
				},
			)
		}
	}

	return validationErrors, nil
}

func assertNthWithPlus(matchingRequestRecords []types.RequestRecord, currentRequestNth int, assert *Condition) ([]ValidationError, error) {
	expected, err := utils.ExtractNumbersFromString(string(assert.Value))
	if err != nil {
		return []ValidationError{}, err
	}

	if expected <= len(matchingRequestRecords) {
		return []ValidationError{}, nil
	}

	return []ValidationError{
		{Code: ValidationErrorCode_NthMismatch, Metadata: map[string]string{
			"nth_requested": fmt.Sprint(currentRequestNth),
			"nth_expected":  string(assert.Value),
		}},
	}, nil
}

func assertMethodMatch(requestRecord *types.RequestRecord, requestRecords []types.RequestRecord, assert *Condition) ([]ValidationError, error) {
	validationErrors := make([]ValidationError, 0)

	if requestRecord.Method != strings.ToLower(string(assert.Value)) {
		validationErrors = append(validationErrors, ValidationError{
			Code: ValidationErrorCode_MethodMismatch,
			Metadata: map[string]string{
				"method_requested": requestRecord.Method,
				"method_expected":  string(assert.Value),
			},
		})
	}

	return validationErrors, nil
}

func assertFormMatch(requestRecord *types.RequestRecord, requestRecords []types.RequestRecord, assert *Condition) ([]ValidationError, error) {
	validationErrors := make([]ValidationError, 0)
	requestBody := string(*requestRecord.Body)

	parsedForm, err := parseForm(requestBody)
	if err != nil {
		panic(err)
	}

	for i := range assert.KeyValues {
		value, ok := parsedForm[i]
		if !ok {
			validationErrors = append(
				validationErrors,
				ValidationError{
					Code: ValidationErrorCode_FormKeyDoesNotExist,
					Metadata: map[string]string{
						"form_key": i,
					},
				},
			)

			continue
		}

		if value != assert.KeyValues[i] {
			validationErrors = append(
				validationErrors,
				ValidationError{
					Code: ValidationErrorCode_FormValueMismatch,
					Metadata: map[string]string{
						"form_key":             i,
						"form_value_requested": value,
						"form_value_expected":  assert.KeyValues[i].(string),
					},
				},
			)
		}
	}

	return validationErrors, nil
}

func assertQuerystringMatch(requestRecord *types.RequestRecord, requestRecords []types.RequestRecord, assert *Condition) ([]ValidationError, error) {
	validationErrors := make([]ValidationError, 0)

	if requestRecord.Querystring == "" {
		return []ValidationError{{
			Code:     ValidationErrorCode_RequestHasNoQuerystring,
			Metadata: map[string]string{},
		}}, nil
	}

	parsedQuery, err := url.ParseQuery(requestRecord.Querystring)
	if err != nil {
		return validationErrors, err
	}

	expectedKeyValuePairs := getKeyValuePairsFromAssertionOptions(assert)

	for _, key := range utils.GetSortedKeys(expectedKeyValuePairs) {
		_, ok := parsedQuery[key]
		if !ok {
			validationErrors = append(
				validationErrors,
				ValidationError{
					Code: ValidationErrorCode_QuerystringKeyNotSet,
					Metadata: map[string]string{
						"querystring_key": key,
					},
				},
			)

			continue
		}

		if expectedKeyValuePairs[key] != parsedQuery[key][0] {
			validationErrors = append(
				validationErrors,
				ValidationError{
					Code: ValidationErrorCode_QuerystringMismatch,
					Metadata: map[string]string{
						"querystring_key":             key,
						"querystring_value_expected":  expectedKeyValuePairs[key].(string),
						"querystring_value_requested": parsedQuery[key][0],
					},
				},
			)
		}
	}

	return validationErrors, nil
}

func assertJsonBodyMatch(requestRecord *types.RequestRecord, requestRecords []types.RequestRecord, assert *Condition) ([]ValidationError, error) {
	validationErrors := make([]ValidationError, 0)

	if len(*requestRecord.Body) == 0 {
		validationErrors = append(
			validationErrors,
			ValidationError{Code: ValidationErrorCode_RequestHasNoBody, Metadata: map[string]string{}},
		)

		return validationErrors, nil
	}

	var jsonResult map[string]interface{}
	err := json.Unmarshal(*requestRecord.Body, &jsonResult)
	if err != nil {
		return validationErrors, err
	}

	jsonValidationResult, err := jsonValidate(jsonResult, assert.KeyValues)
	if err != nil {
		return validationErrors, err
	}

	if !jsonValidationResult {
		assertJson, err := json.Marshal(assert.KeyValues)
		if err != nil {
			panic(err)
		}

		requestRecordReformatted, err := reformatJson(requestRecord.Body)
		if err != nil {
			panic(err)
		}

		validationErrors = append(
			validationErrors,
			ValidationError{Code: ValidationErrorCode_BodyMismatch, Metadata: map[string]string{
				"body_requested": string(requestRecordReformatted),
				"body_expected":  string(assertJson),
			}},
		)
	}

	return validationErrors, nil
}

func reformatJson(jsonEncoded *[]byte) ([]byte, error) {
	var result map[string]interface{}
	err := json.Unmarshal(*jsonEncoded, &result)
	if err != nil {
		return []byte(""), err
	}

	newJsonEncoded, err := json.Marshal(result)
	if err != nil {
		return []byte(""), err
	}

	return newJsonEncoded, nil
}

func parseForm(requestBody string) (map[string]string, error) {
	formValues := make(map[string]string)

	values, err := url.ParseQuery(requestBody)
	if err != nil {
		return formValues, err
	}

	for i := range values {
		formValues[i] = values[i][0]
	}

	return formValues, nil
}

func assertQuerystringExactMatch(requestRecord *types.RequestRecord, requestRecords []types.RequestRecord, assert *Condition) ([]ValidationError, error) {
	validationErrors, err := assertQuerystringMatch(requestRecord, requestRecords, assert)
	if err != nil {
		return validationErrors, err
	}

	parsedQuery, err := url.ParseQuery(requestRecord.Querystring)
	if err != nil {
		return validationErrors, err
	}

	if len(validationErrors) > 0 {
		return validationErrors, nil
	}

	missingKeys := make([]string, 0)
	expectedKeys := make([]string, 0)
	requestedKeys := make([]string, 0)

	if assert.Key != "" {
		expectedKeys = append(expectedKeys, assert.Key)
	}

	for key := range assert.KeyValues {
		expectedKeys = append(expectedKeys, key)
	}

	for key := range parsedQuery {
		requestedKeys = append(requestedKeys, key)

		if utils.IndexOf(expectedKeys, key) == -1 {
			missingKeys = append(missingKeys, key)
		}
	}

	sort.Strings(expectedKeys)
	sort.Strings(requestedKeys)

	if len(missingKeys) > 0 {
		validationErrors = append(
			validationErrors,
			ValidationError{Code: ValidationErrorCode_QuerystringMismatch, Metadata: map[string]string{
				"querystring_keys_expected":  strings.Join(expectedKeys, ","),
				"querystring_keys_requested": strings.Join(requestedKeys, ","),
			}},
		)
	}

	return validationErrors, nil
}

func getKeyValuePairsFromAssertionOptions(assert *Condition) map[string]interface{} {
	keys := make([]string, 0)
	keyValuePairs := make(map[string]interface{}, 0)

	if assert.Key != "" {
		keys = append(keys, assert.Key)
		keyValuePairs[assert.Key] = string(assert.Value)
	}

	for key, value := range assert.KeyValues {
		keys = append(keys, key)
		keyValuePairs[key] = value
	}

	return keyValuePairs
}

func jsonValidate(a, b map[string]interface{}) (bool, error) {
	encodedA, err := json.Marshal(a)
	if err != nil {
		return false, err
	}

	encodedB, err := json.Marshal(b)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(encodedA, encodedB), nil
}
