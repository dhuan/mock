package tests_e2e

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_GetPayload_AllPayload(t *testing.T) {
	RunTest4(
		t,
		[]string{
			"--route foo/bar",
			"--method POST",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} get-payload | {{MOCK_EXECUTABLE}} write`,
			}, ";")),
		},
		Post("foo/bar", nil, []byte("Hello, world. This is the payload.")),
		StringMatches("Hello, world. This is the payload."),
	)
}

func Test_E2E_GetPayload_GetJsonField_OK(t *testing.T) {
	RunTest4(
		t,
		[]string{
			"--route foo/bar",
			"--method POST",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} get-payload foo | {{MOCK_EXECUTABLE}} write`,
			}, ";")),
		},
		Post("foo/bar", JSON_HEADER, []byte(`{"foo": "bar"}`)),
		StringMatches("bar\n"),
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
			t,
			[]string{
				"--route foo/bar",
				"--method POST",
				fmt.Sprintf("--exec '%s'", strings.Join([]string{
					fmt.Sprintf(`{{MOCK_EXECUTABLE}} get-payload %s | {{MOCK_EXECUTABLE}} write`, tc.path),
				}, ";")),
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
			StringMatches(fmt.Sprintf("%+v\n", tc.expect)),
		)
	}
}

func Test_E2E_GetPayload_GetJsonField_FieldDoesNotExist(t *testing.T) {
	RunTest4(
		t,
		[]string{
			"--route foo/bar",
			"--method POST",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} get-payload foo | {{MOCK_EXECUTABLE}} write`,
				`{{MOCK_EXECUTABLE}} get-payload foo`,
				`printf $? | {{MOCK_EXECUTABLE}} write -a`,
			}, ";")),
		},
		Post("foo/bar", JSON_HEADER, []byte(`{"hello": "world"}`)),
		StringMatches("1"),
	)
}

func Test_E2E_GetPayload_GetJsonField_WithEmptyPayload_Exit1(t *testing.T) {
	RunTest4(
		t,
		[]string{
			"--route foo/bar",
			"--method POST",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} get-payload foo | {{MOCK_EXECUTABLE}} write`,
				`{{MOCK_EXECUTABLE}} get-payload foo`,
				`printf $? | {{MOCK_EXECUTABLE}} write -a`,
			}, ";")),
		},
		Post("foo/bar", JSON_HEADER, nil),
		StringMatches("1"),
	)
}

func Test_E2E_GetPayload_GetFieldFromUrlEncodedForm_Ok(t *testing.T) {
	RunTest4(
		t,
		[]string{
			"--route foo/bar",
			"--method POST",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} get-payload foo | {{MOCK_EXECUTABLE}} write`,
				`echo $? | {{MOCK_EXECUTABLE}} write -a`,
			}, ";")),
		},
		PostUrlEncodedForm("foo/bar", map[string]string{
			"foo": "bar",
		}),
		StringMatches("bar\n0\n"),
	)
}

func Test_E2E_GetPayload_GetFieldFromUrlEncodedForm_FieldDoesNotExist(t *testing.T) {
	RunTest4(
		t,
		[]string{
			"--route foo/bar",
			"--method POST",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} get-payload hello | {{MOCK_EXECUTABLE}} write`,
				`{{MOCK_EXECUTABLE}} get-payload hello`,
				`echo $? | {{MOCK_EXECUTABLE}} write -a`,
			}, ";")),
		},
		PostUrlEncodedForm("foo/bar", map[string]string{
			"foo": "bar",
		}),
		StringMatches("1\n"),
	)
}

func Test_E2E_GetPayload_GetFieldFromMultipartForm_Ok(t *testing.T) {
	RunTest4(
		t,
		[]string{
			"--route foo/bar",
			"--method POST",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} get-payload foo | {{MOCK_EXECUTABLE}} write`,
				`{{MOCK_EXECUTABLE}} get-payload foo`,
				`echo $? | {{MOCK_EXECUTABLE}} write -a`,
			}, ";")),
		},
		PostMultipart("foo/bar", map[string]string{
			"foo": "bar",
		}),
		StringMatches("bar\n0\n"),
	)
}

func Test_E2E_GetPayload_GetFieldFromMultipartForm_FieldDoesNotExist(t *testing.T) {
	RunTest4(
		t,
		[]string{
			"--route foo/bar",
			"--method POST",
			fmt.Sprintf("--exec '%s'", strings.Join([]string{
				`{{MOCK_EXECUTABLE}} get-payload doesnotexist | {{MOCK_EXECUTABLE}} write`,
				`{{MOCK_EXECUTABLE}} get-payload doesnotexist`,
				`echo $? | {{MOCK_EXECUTABLE}} write -a`,
			}, ";")),
		},
		PostMultipart("foo/bar", map[string]string{
			"foo": "bar",
		}),
		StringMatches("1\n"),
	)
}
