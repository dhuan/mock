package mock_test

import (
	"testing"

	"github.com/dhuan/mock/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func Test_ToReadableErrors_WithOneError(t *testing.T) {
	validationErrors := []mock.ValidationError{
		{
			Code: mock.ValidationErrorCode_MethodMismatch,
			Metadata: map[string]string{
				"method_requested": "get",
				"method_expected":  "post",
			},
		},
	}

	readableErrors := mock.ToReadableError(validationErrors)

	assert.Equal(
		t,
		`Error: method_mismatch
method_expected: post
method_requested: get`,
		readableErrors,
	)
}

func Test_ToReadableErrors_WithMultipleErrors(t *testing.T) {
	validationErrors := []mock.ValidationError{
		{
			Code: mock.ValidationErrorCode_MethodMismatch,
			Metadata: map[string]string{
				"method_requested": "get",
				"method_expected":  "post",
			},
		},
		{
			Code: mock.ValidationErrorCode_FormKeyDoesNotExist,
			Metadata: map[string]string{
				"form_key": "some_key",
			},
		},
	}

	readableErrors := mock.ToReadableError(validationErrors)

	assert.Equal(
		t,
		`Error: method_mismatch
method_expected: post
method_requested: get

Error: form_key_does_not_exist
form_key: some_key`,
		readableErrors,
	)
}

func Test_ToReadableErrors_WithMultipleErrors_WithoutMetadata(t *testing.T) {
	validationErrors := []mock.ValidationError{
		{
			Code:     mock.ValidationErrorCode_NoCall,
			Metadata: map[string]string{},
		},
		{
			Code: mock.ValidationErrorCode_FormKeyDoesNotExist,
			Metadata: map[string]string{
				"form_key": "some_key",
			},
		},
	}

	readableErrors := mock.ToReadableError(validationErrors)

	assert.Equal(
		t,
		`Error: no_call

Error: form_key_does_not_exist
form_key: some_key`,
		readableErrors,
	)
}
