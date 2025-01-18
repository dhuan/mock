package tests_e2e

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_GetPayload_AllPayload(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			"--route foo/bar",
			"--method POST",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} get-payload > $MOCK_RESPONSE_BODY`,
			}, ";")),
		},
		Post("foo/bar", nil, []byte("Hello, world. This is the payload.")),
		StringMatches("Hello, world. This is the payload."),
	)
}

func Test_E2E_GetPayload_GetJsonField_OK(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			"--route foo/bar",
			"--method POST",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} get-payload foo > $MOCK_RESPONSE_BODY`,
			}, ";")),
		},
		Post("foo/bar", JSON_HEADER, []byte(`{"foo": "bar"}`)),
		StringMatches("bar"),
	)
}

func Test_E2E_GetPayload_GetJsonField_ArrayRoot(t *testing.T) {
	for _, tc := range []struct {
		path             string
		expect           interface{}
		expectStatusCode int
	}{
		{"[0].location", "earth", 0},
		{"[1].location", "mars", 0},
		{"[1]", "{\"location\":\"mars\"}", 0},
		{"[2]", "", 1},
	} {
		RunTest4(
			t, nil,
			[]string{
				"--route foo/bar",
				"--method POST",
				CmdExec(fmt.Sprintf(`({{MOCK_EXECUTABLE}} get-payload %s) > $MOCK_RESPONSE_BODY`, tc.path)),
			},
			Post("foo/bar", JSON_HEADER, []byte(`[{"location":"earth"},{"location":"mars"}]`)),
			StringMatches(fmt.Sprintf("%+v", tc.expect)),
			ExitCodeHeaderMatches(fmt.Sprintf("%d", tc.expectStatusCode)),
		)
	}
}

func Test_E2E_GetPayload_GetJsonField_InvalidJson(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			"--route foo/bar",
			"--method POST",
			CmdExec(`{{MOCK_EXECUTABLE}} get-payload foo > $MOCK_RESPONSE_BODY`),
		},
		Post("foo/bar", JSON_HEADER, []byte(`{This is invalid JSON}`)),
		StringMatches(""),
		ExitCodeHeaderMatches("1"),
	)
}

func Test_E2E_GetPayload_GetJsonField_Nested_OK(t *testing.T) {
	for _, tc := range []struct {
		path   string
		expect interface{}
	}{
		{"users[0].location", "earth"},
		{"users[1].location", "mars"},
		{"users[1].age", 30},
		{"users[0]", `{"age":20,"likes":[],"location":"earth"}`},
		{"users[1].likes", `["food","music"]`},
		{"users[1].likes[1]", `music`},
	} {
		RunTest4(
			t, nil,
			[]string{
				"--route foo/bar",
				"--method POST",
				CmdExec(fmt.Sprintf(`{{MOCK_EXECUTABLE}} get-payload %s > $MOCK_RESPONSE_BODY`, tc.path)),
			},
			Post("foo/bar", JSON_HEADER, []byte(`{
  "users": [
    {
      "location": "earth",
      "age": 20,
      "likes": []
    },
    {
      "location": "mars",
      "age": 30,
      "likes": [
        "food",
        "music"
      ]
    }
  ]
}`)),
			StringMatches(fmt.Sprintf("%+v", tc.expect)),
		)
	}
}

func Test_E2E_GetPayload_GetJsonField_Nested_InvalidFields(t *testing.T) {
	for _, tc := range []struct {
		path string
	}{
		{"users[1].location"},
		{"users[0].foo"},
		{"users[0].likes[20]"},
		{"foo"},
		{"[0]"},
	} {
		RunTest4(
			t, nil,
			[]string{
				"--route foo/bar",
				"--method POST",
				CmdExec(fmt.Sprintf(`{{MOCK_EXECUTABLE}} get-payload %s > $MOCK_RESPONSE_BODY`, tc.path)),
			},
			Post("foo/bar", JSON_HEADER, []byte(`{
  "users": [
    {
      "location": "earth",
      "age": 20,
      "likes": ["movies"]
    }
  ]
}`)),
			StringMatches(""),
			ExitCodeHeaderMatches("1"),
		)
	}
}

func Test_E2E_GetPayload_GetJsonField_FieldDoesNotExist(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			"--route foo/bar",
			"--method POST",
			CmdExec(`{{MOCK_EXECUTABLE}} get-payload foo > $MOCK_RESPONSE_BODY`),
		},
		Post("foo/bar", JSON_HEADER, []byte(`{"hello": "world"}`)),
		StringMatches(""),
		ExitCodeHeaderMatches("1"),
	)
}

func Test_E2E_GetPayload_GetJsonField_WithEmptyPayload_Exit1(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			"--route foo/bar",
			"--method POST",
			CmdExec(`{{MOCK_EXECUTABLE}} get-payload foo > $MOCK_RESPONSE_BODY`),
		},
		Post("foo/bar", JSON_HEADER, nil),
		StringMatches(""),
		ExitCodeHeaderMatches("1"),
	)
}

func Test_E2E_GetPayload_GetFieldFromUrlEncodedForm_Ok(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			"--route foo/bar",
			"--method POST",
			CmdExec(`{{MOCK_EXECUTABLE}} get-payload foo > $MOCK_RESPONSE_BODY`),
		},
		PostUrlEncodedForm("foo/bar", map[string]string{
			"foo": "bar",
		}),
		StringMatches("bar"),
		ExitCodeHeaderMatches("0"),
	)
}

func Test_E2E_GetPayload_GetFieldFromUrlEncodedForm_FieldDoesNotExist(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			"--route foo/bar",
			"--method POST",
			CmdExec(`{{MOCK_EXECUTABLE}} get-payload hello > $MOCK_RESPONSE_BODY`),
		},
		PostUrlEncodedForm("foo/bar", map[string]string{
			"foo": "bar",
		}),
		StringMatches(""),
		ExitCodeHeaderMatches("1"),
	)
}

func Test_E2E_GetPayload_GetFieldFromMultipartForm_Ok(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			"--route foo/bar",
			"--method POST",
			CmdExec(`{{MOCK_EXECUTABLE}} get-payload foo > $MOCK_RESPONSE_BODY`),
		},
		PostMultipart("foo/bar", map[string]string{
			"foo": "bar",
		}),
		StringMatches("bar"),
		ExitCodeHeaderMatches("0"),
	)
}

func Test_E2E_GetPayload_GetFieldFromMultipartForm_FieldDoesNotExist(t *testing.T) {
	RunTest4(
		t, nil,
		[]string{
			"--route foo/bar",
			"--method POST",
			CmdExec(`{{MOCK_EXECUTABLE}} get-payload doesnotexist > $MOCK_RESPONSE_BODY`),
		},
		PostMultipart("foo/bar", map[string]string{
			"foo": "bar",
		}),
		StringMatches(""),
		ExitCodeHeaderMatches("1"),
	)
}
