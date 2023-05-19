package tests_e2e

import (
	"strings"
	"testing"
	"time"

	e2eutils "github.com/dhuan/mock/tests/e2e/utils"
	"github.com/stretchr/testify/assert"
)

// The "delay" flag is not used here, therefore the request ends
// quickly (in less than 2 seconds)
func Test_E2E_Delay_WithoutDelay(t *testing.T) {
	killMock, _, mockConfig := e2eutils.RunMockBg(
		e2eutils.NewState(),
		"serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}",
		nil,
	)
	defer killMock()

	timeBeforeRequest := time.Now()
	response := e2eutils.Request(mockConfig, "POST", "foo/bar", strings.NewReader(`{"foo":"bar"}`), e2eutils.ContentTypeJsonHeaders)
	timeAfterRequest := time.Now()

	assert.Equal(t, 200, response.StatusCode)

	e2eutils.AssertTimeDifferenceLessThanSeconds(
		t,
		timeBeforeRequest,
		timeAfterRequest,
		2,
	)
}

// The "delay" flag is used, set to 3 seconds
func Test_E2E_Delay_WithDelay(t *testing.T) {
	killMock, _, mockConfig := e2eutils.RunMockBg(
		e2eutils.NewState(),
		"serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}} --delay 3000",
		nil,
	)
	defer killMock()

	timeBeforeRequest := time.Now()
	response := e2eutils.Request(mockConfig, "POST", "foo/bar", strings.NewReader(`{"foo":"bar"}`), e2eutils.ContentTypeJsonHeaders)
	timeAfterRequest := time.Now()

	assert.Equal(t, 200, response.StatusCode)

	e2eutils.AssertTimeDifferenceEqualOrMoreThanSeconds(
		t,
		timeBeforeRequest,
		timeAfterRequest,
		3,
	)
}
