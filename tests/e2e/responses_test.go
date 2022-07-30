package tests_e2e

import (
	"encoding/json"
	"testing"

	mocklib "github.com/dhuan/mock/pkg/mock"
	e2eutils "github.com/dhuan/mock/tests/e2e/utils"
	"github.com/stretchr/testify/assert"
)

func Test_E2E_Response_FileResponse(t *testing.T) {
	runResponseTest(t, "POST", "foo/bar/1", stringMatches("Hello world! This is response A.\n"))
}

func Test_E2E_Response_ResponseInsideFolder(t *testing.T) {
	runResponseTest(t, "POST", "foo/bar/2", stringMatches("This test asserts that you can set response files inside folders.\n"))
}

func Test_E2E_Response_JsonResponse(t *testing.T) {
	runResponseTest(t, "POST", "foo/bar/3", jsonMatches(map[string]interface{}{
		"response_text": "This is a JSON response.",
	}))
}

func runResponseTest(t *testing.T, method, route string, assertionFunc func(t *testing.T, response []byte)) {
	killMock := e2eutils.RunMockBg(e2eutils.NewState(), "serve -c {{TEST_DATA_PATH}}/config_with_file_responses/config.json -p {{TEST_E2E_PORT}}")
	defer killMock()

	mockConfig := mocklib.Init("localhost:4000")
	responseBody := e2eutils.Request(mockConfig, method, route, "", map[string]string{})

	assertionFunc(t, responseBody)
}

func stringMatches(expected string) func(t *testing.T, response []byte) {
	return func(t *testing.T, responseBody []byte) {
		assert.Equal(t, expected, string(responseBody))
	}
}

func jsonMatches(expectedJson map[string]interface{}) func(t *testing.T, response []byte) {
	return func(t *testing.T, responseBody []byte) {
		jsonEncodedA, err := json.Marshal(expectedJson)
		if err != nil {
			t.Fatal("Failed to parse JSON from expected input!")
		}

		jsonEncodedB, err := encodeJsonAgain(responseBody)
		if err != nil {
			t.Fatal("Failed to parse JSON from response!")
		}

		assert.Equal(t, string(jsonEncodedA), string(jsonEncodedB))
	}
}

func encodeJsonAgain(encodedJson []byte) ([]byte, error) {
	var jsonTarget map[string]interface{}
	err := json.Unmarshal(encodedJson, &jsonTarget)
	if err != nil {
		return []byte(""), err
	}

	return json.Marshal(jsonTarget)
}
