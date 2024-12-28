package map_extract_test

import (
	"encoding/json"
	"github.com/dhuan/mock/internal/map_extract"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Extract_Basics(t *testing.T) {
	testExtract(t, `{
		"foo": "bar"
	}`, "foo", "bar", true)
}

func Test_Extract_Nested(t *testing.T) {
	for _, tc := range []struct {
		path   string
		expect interface{}
	}{
		{"user.location", "berlin"},
		{"user.age", float64(20)},
	} {
		testExtract(t, `{
			"user": {
				"location": "berlin",
				"age": 20
			}
		}`, tc.path, tc.expect, true)
	}
}

func Test_Extract_Arrays(t *testing.T) {
	for _, tc := range []struct {
		path   string
		expect interface{}
	}{
		{"users[0].location", "berlin"},
		{"users[1].location", "london"},
		{"users[1]", `{"age":30,"likes":["music"],"location":"london"}`},
		{"users[0].likes", `["food","movies"]`},
		{"users[0].likes[1]", `movies`},
	} {
		testExtract(t, `{
			"users": [{
				"location": "berlin",
				"age": 20,
				"likes": ["food", "movies"]
			}, {
				"location": "london",
				"age": 30,
				"likes": ["music"]
			}]
		}`, tc.path, tc.expect, true)
	}
}

func Test_Extract_ArrayRoot(t *testing.T) {
	for _, tc := range []struct {
		path   string
		expect interface{}
	}{
		{"[0].location", "berlin"},
		{"[1].location", "london"},
		{"[0]", `{"location":"berlin"}`},
	} {
		testExtract(t, `[{"location":"berlin"},{"location":"london"}]`, tc.path, tc.expect, true)
	}
}

func testExtract(t *testing.T, jsonData string, fieldName string, expectedValue interface{}, expectedOk bool) {
	result, ok := map_extract.Extract(parse(jsonData), fieldName)

	assert.Equal(
		t,
		expectedOk,
		ok,
	)

	assert.Equal(
		t,
		expectedValue,
		result,
	)
}

func parse(data string) interface{} {
	var result interface{}
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		panic(err)
	}

	return result
}
