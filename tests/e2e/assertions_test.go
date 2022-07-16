package tests_e2e

import (
	"testing"

	mocklib "github.com/dhuan/mock/pkg/mock"
	e2eutils "github.com/dhuan/mock/tests/e2e/utils"
	"github.com/stretchr/testify/assert"
)

func Test_E2E_Assertion_NoCalls(t *testing.T) {
	e2eutils.RunMockBg(e2eutils.NewState(), "serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}")

	mockState := mocklib.Init("localhost:4000")
	validationErrors := mocklib.Assert(mockState, &mocklib.AssertConfig{
		Route:  "foo/bar",
		Method: "POST",
	})

	assert.Equal(
		t,
		[]mocklib.ValidationError{
			{Code: "no_call", Metadata: map[string]string{}},
		},
		validationErrors,
	)
}

func Test_E2E_Assertion_BasicAssertion(t *testing.T) {
	e2eutils.RunMockBg(e2eutils.NewState(), "serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}")

	mockState := mocklib.Init("localhost:4000")
	validationErrors := mocklib.Assert(mockState, &mocklib.AssertConfig{
		Route:  "foo/bar",
		Method: "POST",
	})

	assert.Equal(
		t,
		[]mocklib.ValidationError{
			{Code: "no_call", Metadata: map[string]string{}},
		},
		validationErrors,
	)
}
