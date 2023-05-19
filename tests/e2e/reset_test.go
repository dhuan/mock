package tests_e2e

import (
	"strings"
	"testing"

	mocklib "github.com/dhuan/mock/pkg/mock"
	e2eutils "github.com/dhuan/mock/tests/e2e/utils"
	"github.com/stretchr/testify/assert"
)

func Test_E2E_Resetting(t *testing.T) {
	killMock, serverOutput, mockConfig := e2eutils.RunMockBg(
		e2eutils.NewState(),
		"serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}",
		nil,
	)
	defer killMock()

	e2eutils.Request(mockConfig, "POST", "foo/bar", strings.NewReader(`{"foo":"bar"}`), e2eutils.ContentTypeJsonHeaders)

	validationErrors := e2eutils.MockAssert(&mocklib.AssertConfig{
		Route: "foo/bar",
		Assert: &mocklib.Condition{
			Type: mocklib.ConditionType_JsonBodyMatch,
			KeyValues: map[string]interface{}{
				"foo": "bar",
			},
		},
	}, serverOutput)

	assert.Equal(t, 0, len(validationErrors))

	e2eutils.RequestApiReset(mockConfig)

	validationErrors = e2eutils.MockAssert(&mocklib.AssertConfig{
		Route: "foo/bar",
		Assert: &mocklib.Condition{
			Type: mocklib.ConditionType_JsonBodyMatch,
			KeyValues: map[string]interface{}{
				"foo": "bar",
			},
		},
	}, serverOutput)

	assert.Equal(
		t,
		[]mocklib.ValidationError{
			{Code: mocklib.ValidationErrorCode_NoCall, Metadata: map[string]string{}},
		},
		validationErrors,
	)
}
