package tests_e2e

import (
	"fmt"
	"strings"
	"testing"
	"time"

	. "github.com/dhuan/mock/tests/e2e/utils"
	"github.com/stretchr/testify/assert"
)

// The "delay" flag is not used here, therefore the request ends
// quickly (in less than 2 seconds)
func Test_E2E_Delay_WithoutDelay(t *testing.T) {
	killMock, serverOutput, mockConfig, _ := RunMockBg(
		NewState(),
		"serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}}",
		nil,
		true,
		nil,
	)
	defer killMock()

	timeBeforeRequest := time.Now()
	response := Request(mockConfig, "POST", "foo/bar", strings.NewReader(`{"foo":"bar"}`), ContentTypeJsonHeaders, serverOutput)
	timeAfterRequest := time.Now()

	assert.Equal(t, 200, response.StatusCode)

	AssertTimeDifferenceLessThanSeconds(
		t,
		timeBeforeRequest,
		timeAfterRequest,
		2,
	)
}

// The "delay" flag is used, set to 3 seconds
func Test_E2E_Delay_WithDelay(t *testing.T) {
	killMock, serverOutput, mockConfig, _ := RunMockBg(
		NewState(),
		"serve -c {{TEST_DATA_PATH}}/config_basic/config.json -p {{TEST_E2E_PORT}} --delay 3000",
		nil,
		true,
		nil,
	)
	defer killMock()

	timeBeforeRequest := time.Now()
	response := Request(mockConfig, "POST", "foo/bar", strings.NewReader(`{"foo":"bar"}`), ContentTypeJsonHeaders, serverOutput)
	timeAfterRequest := time.Now()

	assert.Equal(t, 200, response.StatusCode)

	AssertTimeDifferenceEqualOrMoreThanSeconds(
		t,
		timeBeforeRequest,
		timeAfterRequest,
		3,
	)
}

func Test_E2E_Delay_WithBaseApi(t *testing.T) {
	state := NewState()
	killMockBase, _, _, _ := RunMockBg(
		state,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			"--route 'foo/bar'",
			"--response 'Hello world! This is the base API.'",
		}, " "),
		nil,
		true,
		nil,
	)
	defer killMockBase()

	state2 := NewState()
	killMockBase2, serverOutput, mockConfig, _ := RunMockBg(
		state2,
		strings.Join([]string{
			"serve",
			"-p {{TEST_E2E_PORT}}",
			fmt.Sprintf("--base 'localhost:%d'", state.Port),
			"--delay 3000",
		}, " "),
		nil,
		true,
		nil,
	)
	defer killMockBase2()

	timeBeforeRequest := time.Now()
	response := Request(mockConfig, "GET", "foo/bar", nil, nil, serverOutput)
	timeAfterRequest := time.Now()

	assert.Equal(t, 200, response.StatusCode)

	AssertTimeDifferenceEqualOrMoreThanSeconds(
		t,
		timeBeforeRequest,
		timeAfterRequest,
		3,
	)
}
