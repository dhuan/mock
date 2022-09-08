package tests_e2e

import (
	"testing"

	mocklib "github.com/dhuan/mock/pkg/mock"
	e2eutils "github.com/dhuan/mock/tests/e2e/utils"
	"github.com/stretchr/testify/assert"
)

func Test_E2E_Resetting(t *testing.T) {
	killMock := e2eutils.RunMockBg(e2eutils.NewState(), "serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}")
	defer killMock()

	mockConfig := mocklib.Init("localhost:4000")
	e2eutils.Request(mockConfig, "POST", "foo/bar", `{"foo":"bar"}`, e2eutils.ContentTypeJsonHeaders)

	validationErrors := e2eutils.MockAssert(&mocklib.AssertConfig{
		Route: "foo/bar",
		Assert: &mocklib.AssertOptions{
			Type: mocklib.AssertType_JsonBodyMatch,
			Data: map[string]interface{}{
				"foo": "bar",
			},
		},
	})

	assert.Equal(t, 0, len(validationErrors))

	e2eutils.RequestApiReset(mockConfig)

	validationErrors = e2eutils.MockAssert(&mocklib.AssertConfig{
		Route: "foo/bar",
		Assert: &mocklib.AssertOptions{
			Type: mocklib.AssertType_JsonBodyMatch,
			Data: map[string]interface{}{
				"foo": "bar",
			},
		},
	})

	assert.Equal(
		t,
		[]mocklib.ValidationError{
			{Code: mocklib.ValidationErrorCode_NoCall, Metadata: map[string]string{}},
		},
		validationErrors,
	)
}
