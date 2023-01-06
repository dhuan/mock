package mock_test

import (
	"testing"

	. "github.com/dhuan/mock/internal/test_unit_utils"
	. "github.com/dhuan/mock/pkg/mock"
)

func Test_Validate_NoCalls(t *testing.T) {
	RunUnitTest(
		t,
		Validate("foobar", AssertOptionsWithKeyValue(ConditionType_HeaderMatch, "foo", "bar")),
		ExpectOneValidationError(ValidationErrorCode_NoCall, map[string]string{}),
	)
}

func Test_Validate_HeaderNotIncluded(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecordWithHeaders("foobar", [][]string{[]string{"some_header_key", "some_header_value"}}),
		Validate("foobar", AssertOptionsWithKeyValue(ConditionType_HeaderMatch, "foo", "bar")),
		ExpectOneValidationError(ValidationErrorCode_HeaderNotIncluded, map[string]string{
			"missing_header_key": "foo",
		}),
	)
}

func Test_Validate_HeaderNotIncludedMany(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecordWithHeaders("foobar", [][]string{[]string{"some_header_key", "some_header_value"}}),
		Validate("foobar", AssertOptionsWithKeyValues(ConditionType_HeaderMatch, map[string]interface{}{
			"foo":  "bar",
			"foo2": "bar2",
		})),
		ExpectValidationErrorsCount(2),
		ExpectValidationErrorNth(0, ValidationErrorCode_HeaderNotIncluded, map[string]string{
			"missing_header_key": "foo",
		}),
		ExpectValidationErrorNth(1, ValidationErrorCode_HeaderNotIncluded, map[string]string{
			"missing_header_key": "foo2",
		}),
	)
}

func Test_Validate_HeaderMismatch_Single(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecordWithHeaders("foobar", [][]string{[]string{"some_header_key", "some_header_value"}}),
		Validate("foobar", AssertOptionsWithKeyValue(ConditionType_HeaderMatch, "some_header_key", "a_different_header_value")),
		ExpectOneValidationError(ValidationErrorCode_HeaderValueMismatch, map[string]string{
			"header_key":             "some_header_key",
			"header_value_requested": "some_header_value",
			"header_value_expected":  "a_different_header_value",
		}),
	)
}

func Test_Validate_HeaderMismatch_Many(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecordWithHeaders("foobar", [][]string{[]string{"some_header_key", "some_header_value"}}),
		Validate("foobar", AssertOptionsWithKeyValues(ConditionType_HeaderMatch, map[string]interface{}{
			"some_header_key": "a_different_header_value",
		})),
		ExpectOneValidationError(ValidationErrorCode_HeaderValueMismatch, map[string]string{
			"header_key":             "some_header_key",
			"header_value_requested": "some_header_value",
			"header_value_expected":  "a_different_header_value",
		}),
	)
}

func Test_Validate_WithAndChainingAssertingMethodAndHeader_Fail(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecordWithHeaders("foobar", [][]string{[]string{"some_header_key", "some_header_value"}}),
		Validate("foobar",
			&Condition{
				Type: ConditionType_HeaderMatch,
				KeyValues: map[string]interface{}{
					"some_header_key": "some_header_value",
				},
				And: &Condition{
					Type:  ConditionType_MethodMatch,
					Value: "post",
				},
			},
		),
		ExpectOneValidationError(ValidationErrorCode_MethodMismatch, map[string]string{
			"method_requested": "get",
			"method_expected":  "post",
		}),
	)
}

func Test_Validate_WithAndChainingAssertingMethodAndHeader(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecordWithHeaders("foobar", [][]string{[]string{"some_header_key", "some_header_value"}}),
		Validate("foobar",
			&Condition{
				Type: ConditionType_HeaderMatch,
				KeyValues: map[string]interface{}{
					"some_header_key": "some_header_value",
				},
				And: &Condition{
					Type:  ConditionType_MethodMatch,
					Value: "get",
				},
			},
		),
		ExpectZeroValidationErrors,
	)
}

func Test_Validate_JsonBodyAssertion_Match(t *testing.T) {
	RunUnitTest(
		t,
		AddPostRequestRecordWithPayload("foobar", `{"foo":"bar", "some_key": "some_value"}`),
		Validate("foobar", AssertOptionsWithData(ConditionType_JsonBodyMatch, map[string]interface{}{
			"foo":      "bar",
			"some_key": "some_value",
		})),
		ExpectZeroValidationErrors,
	)
}

func Test_Validate_JsonBodyAssertion_Mismatch(t *testing.T) {
	RunUnitTest(
		t,
		AddPostRequestRecordWithPayload("foobar", `{"foo":"bar", "some_key": "some_value"}`),
		Validate("foobar", AssertOptionsWithData(ConditionType_JsonBodyMatch, map[string]interface{}{
			"foo":         "bar",
			"some_key":    "some_value",
			"another_key": "another_value",
		})),
		ExpectOneValidationError(ValidationErrorCode_BodyMismatch, map[string]string{
			"body_requested": `{"foo":"bar","some_key":"some_value"}`,
			"body_expected":  `{"another_key":"another_value","foo":"bar","some_key":"some_value"}`,
		}),
	)
}

func Test_Validate_JsonBodyAssertion_Fail_RequestHasNoBody(t *testing.T) {
	RunUnitTest(
		t,
		AddPostRequestRecordWithPayload("foobar", ``),
		Validate("foobar", AssertOptionsWithData(ConditionType_JsonBodyMatch, map[string]interface{}{
			"foo": "bar",
		})),
		ExpectOneValidationError(ValidationErrorCode_RequestHasNoBody, map[string]string{}),
	)
}

func Test_Validate_Nth(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecord("foobar"),
		AddPostRequestRecordWithPayload("foobar", `{"foo":"bar", "some_key": "some_value"}`),

		Validate("foobar", AssertOptionsWithValue(ConditionType_MethodMatch, "get")),
		ExpectZeroValidationErrors,

		RemoveValidationErrors,
		ValidateNth(2, "foobar", AssertOptionsWithValue(ConditionType_MethodMatch, "get")),
		ExpectOneValidationError(ValidationErrorCode_MethodMismatch, map[string]string{
			"method_requested": "post",
			"method_expected":  "get",
		}),

		RemoveValidationErrors,
		ValidateNth(1, "foobar", AssertOptionsWithValue(ConditionType_MethodMatch, "get")),
		ExpectZeroValidationErrors,
	)
}

func Test_Validate_Nth_OutOfRange(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecord("foobar"),
		AddPostRequestRecordWithPayload("foobar", `{"foo":"bar", "some_key": "some_value"}`),
		ValidateNth(3, "foobar", AssertOptionsWithValue(ConditionType_MethodMatch, "get")),
		ExpectOneValidationError(ValidationErrorCode_NthOutOfRange, map[string]string{}),
	)
}

func Test_Validate_FormMatch_FormKeyNotExisting(t *testing.T) {
	RunUnitTest(
		t,
		AddPostRequestRecordWithPayload("foobar", `{"foo":"bar", "hello": "world"}`),
		Validate("foobar", AssertOptionsWithKeyValues(ConditionType_FormMatch, map[string]interface{}{
			"some_key": "some_value",
		})),
		ExpectOneValidationError(ValidationErrorCode_FormKeyDoesNotExist, map[string]string{
			"form_key": "some_key",
		}),
	)
}

func Test_Validate_FormMatch_FormValueMismatch(t *testing.T) {
	RunUnitTest(
		t,
		AddPostRequestRecordWithPayload("foobar", `foo=bar&hello=world`),
		Validate("foobar", AssertOptionsWithKeyValues(ConditionType_FormMatch, map[string]interface{}{
			"foo": "not_bar",
		})),
		ExpectOneValidationError(ValidationErrorCode_FormValueMismatch, map[string]string{
			"form_key":             "foo",
			"form_value_requested": "bar",
			"form_value_expected":  "not_bar",
		}),
	)
}

func Test_Validate_Querystring_FailBecauseRequestHasNoQuerystring(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecord("foobar"),
		Validate("foobar", AssertOptionsWithKeyValue(ConditionType_QuerystringMatch, "foo", "bar")),
		ExpectOneValidationError(ValidationErrorCode_RequestHasNoQuerystring, map[string]string{}),
	)
}

func Test_Validate_Querystring_FailBecauseQuerystringDoesNotMatch(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecord("foobar?foo=not_bar"),
		Validate("foobar", AssertOptionsWithKeyValue(ConditionType_QuerystringMatch, "foo", "bar")),
		ExpectOneValidationError(ValidationErrorCode_QuerystringMismatch, map[string]string{
			"querystring_key":             "foo",
			"querystring_value_expected":  "bar",
			"querystring_value_requested": "not_bar",
		}),
	)
}

func Test_Validate_Querystring_Matching_WithOne(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecord("foobar?foo=bar"),
		Validate("foobar", AssertOptionsWithKeyValue(ConditionType_QuerystringMatch, "foo", "bar")),
		ExpectValidationErrorsCount(0),
	)
}

func Test_Validate_Querystring_Failing_WithMany(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecord("foobar?foo=not_bar&hello=ola"),
		Validate("foobar", AssertOptionsWithKeyValues(ConditionType_QuerystringMatch, map[string]interface{}{
			"foo":   "bar",
			"hello": "world",
		})),
		ExpectValidationErrorsCount(2),
		ExpectValidationErrorNth(0, ValidationErrorCode_QuerystringMismatch, map[string]string{
			"querystring_key":             "foo",
			"querystring_value_expected":  "bar",
			"querystring_value_requested": "not_bar",
		}),
		ExpectValidationErrorNth(1, ValidationErrorCode_QuerystringMismatch, map[string]string{
			"querystring_key":             "hello",
			"querystring_value_expected":  "world",
			"querystring_value_requested": "ola",
		}),
	)
}

func Test_Validate_Querystring_Passing_WithMany(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecord("foobar?foo=bar&hello=world&somekey=somevalue"),
		Validate("foobar", AssertOptionsWithKeyValues(ConditionType_QuerystringMatch, map[string]interface{}{
			"foo":   "bar",
			"hello": "world",
		})),
		ExpectZeroValidationErrors,
	)
}

func Test_Validate_Querystring_FailBecauseExpectedQuerystringKeyWasNotInTheRequest(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecord("foobar?foo=bar"),
		Validate("foobar", AssertOptionsWithKeyValue(ConditionType_QuerystringExactMatch, "hello", "world")),
		ExpectOneValidationError(ValidationErrorCode_QuerystringKeyNotSet, map[string]string{
			"querystring_key": "hello",
		}),
	)
}

func Test_Validate_QuerystringExact_Passing_WithOne(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecord("foobar?foo=bar"),
		Validate("foobar", AssertOptionsWithKeyValues(ConditionType_QuerystringExactMatch, map[string]interface{}{
			"foo": "bar",
		})),
		ExpectZeroValidationErrors,
	)
}

func Test_Validate_QuerystringExact_Passing_WithMany(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecord("foobar?foo=bar&some_key=some_value&another_key=another_value"),
		Validate(
			"foobar",
			&Condition{
				Type:  ConditionType_QuerystringExactMatch,
				Key:   "foo",
				Value: "bar",
				KeyValues: map[string]interface{}{
					"some_key":    "some_value",
					"another_key": "another_value",
				},
			},
		),
		ExpectZeroValidationErrors,
	)
}

func Test_Validate_QuerystringExact_FailingBecauseValuesDontMatch(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecord("foobar?foo=not_bar"),
		Validate("foobar", AssertOptionsWithKeyValues(ConditionType_QuerystringExactMatch, map[string]interface{}{
			"foo": "bar",
		})),
		ExpectOneValidationError(ValidationErrorCode_QuerystringMismatch, map[string]string{
			"querystring_key":             "foo",
			"querystring_value_expected":  "bar",
			"querystring_value_requested": "not_bar",
		}),
	)
}

func Test_Validate_QuerystringExact_FailingBecauseExpectedQuerystringKeyWasNotInTheRequest(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecord("foobar?hello=world"),
		Validate("foobar", AssertOptionsWithKeyValue(ConditionType_QuerystringExactMatch, "foo", "bar")),
		ExpectOneValidationError(ValidationErrorCode_QuerystringKeyNotSet, map[string]string{
			"querystring_key": "foo",
		}),
	)
}

func Test_Validate_QuerystringExact_FailingBecauseExpectedQuerystringKeyWasNotInTheRequest_WithMany(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecord("foobar?foo=bar&some_key=some_value&another_key=another_value"),
		Validate("foobar", AssertOptionsWithKeyValues(ConditionType_QuerystringExactMatch, map[string]interface{}{
			"foo":         "bar",
			"another_key": "another_value",
		})),
		ExpectOneValidationError(ValidationErrorCode_QuerystringMismatch, map[string]string{
			"querystring_keys_expected":  "another_key,foo",
			"querystring_keys_requested": "another_key,foo,some_key",
		}),
	)
}

func Test_Validate_QuerystringExact_FailingBecauseItDoesNotMatchExactly(t *testing.T) {
	RunUnitTest(
		t,
		AddGetRequestRecord("foobar?foo=bar&hello=world"),
		Validate("foobar", AssertOptionsWithKeyValue(ConditionType_QuerystringExactMatch, "foo", "bar")),
		ExpectOneValidationError(ValidationErrorCode_QuerystringMismatch, map[string]string{
			"querystring_keys_expected":  "foo",
			"querystring_keys_requested": "foo,hello",
		}),
	)
}
