package tests_e2e

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	. "github.com/dhuan/mock/tests/e2e/utils"
)

func Test_E2E_Response_Fileserver(t *testing.T) {
	RunTest(
		t,
		"config_with_static_files/config.json",
		"GET",
		"foo/bar/hello.txt",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(200),
		StringMatches("Hello world!\n"),
	)
}

func Test_E2E_Response_Fileserver_UnexistingFile(t *testing.T) {
	RunTest(
		t,
		"config_with_static_files/config.json",
		"GET",
		"foo/bar/this_file_does_not_exist.txt",
		nil,
		strings.NewReader(""),
		StatusCodeMatches(404),
		StringMatches("File does not exist: this_file_does_not_exist.txt"),
	)
}

func Test_E2E_Response_Fileserver_WithCmdParams(t *testing.T) {
	responseFormats := []string{
		"--response 'fs:data/config_with_static_files/public'",
		"--response-file-server 'data/config_with_static_files/public'",
		"--file-server 'data/config_with_static_files/public'",
	}

	for _, responseFormat := range responseFormats {
		RunTestWithNoConfigAndWithArgs(
			t,
			[]string{
				"--route foo/bar/*",
				responseFormat,
			},
			"GET",
			"foo/bar/hello.txt",
			nil,
			strings.NewReader(""),
			StatusCodeMatches(200),
			StringMatches("Hello world!\n"),
		)
	}
}

func Test_E2E_Response_Fileserver_WithCmdParams_WithAbsolutePath(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	folderPath := fmt.Sprintf("%s/data/config_with_static_files/public", pwd)

	responseFormats := []string{
		fmt.Sprintf("--response 'fs:%s'", folderPath),
		fmt.Sprintf("--response-file-server '%s'", folderPath),
		fmt.Sprintf("--file-server '%s'", folderPath),
	}

	for _, responseFormat := range responseFormats {
		RunTestWithNoConfigAndWithArgs(
			t,
			[]string{
				"--route foo/bar/*",
				responseFormat,
			},
			"GET",
			"foo/bar/hello.txt",
			nil,
			strings.NewReader(""),
			StatusCodeMatches(200),
			StringMatches("Hello world!\n"),
		)
	}
}

func Test_E2E_Response_Fileserver_AutomaticHeaders(t *testing.T) {
	some_text_file_data := "Hello, world!"
	json_data := `{"hello": "world"}`
	dummy_data := `dummy_data`

	path := CreateTmpEnvironment(
		FileEntry("some_text_file.txt", []byte(some_text_file_data)),
		FileEntry("data.json", []byte(json_data)),
		FileEntry("image.jpg", []byte(dummy_data)),
	)

	for _, tc := range []struct {
		fileToRequest       string
		expectedData        []byte
		expectedContentType string
	}{
		{"some_text_file.txt", []byte(some_text_file_data), "text/plain; charset=utf-8"},
		{"data.json", []byte(json_data), "application/json"},
		{"image.jpg", []byte(dummy_data), "image/jpeg"},
	} {
		RunTest2(
			t,
			[]string{
				"--route public/*",
				"--file-server .",
			},
			"GET",
			fmt.Sprintf("public/%s", tc.fileToRequest),
			nil,
			strings.NewReader(""),
			&RunMockOptions{
				Wd: path,
			},
			StatusCodeMatches(200),
			StringMatches(string(tc.expectedData)),
			HeadersMatch(http.Header{
				"Content-Type": {tc.expectedContentType},
			}),
		)
	}
}

func Test_E2E_Response_Fileserver_Navigation(t *testing.T) {
	some_text_file_data := "Hello, world!"
	json_data := `{"hello": "world"}`
	dummy_data := `dummy_data`

	path := CreateTmpEnvironment(
		FileEntry("some_text_file.txt", []byte(some_text_file_data)),
		FileEntry("data.json", []byte(json_data)),
		FileEntry("image.jpg", []byte(dummy_data)),
		DirEntry("some_folder", []FsEntry{}),
	)

	RunTest2(
		t,
		[]string{
			"--route public/*",
			"--file-server .",
		},
		"GET",
		"public/",
		nil,
		strings.NewReader(""),
		&RunMockOptions{
			Wd: path,
		},
		StatusCodeMatches(200),
		RemoveUntestableDataFromFileserverHtmlOutput,
		TidyUpHtmlResponse,
		MatchesFile("data/html_match/fileserver_navigation_index.html"),
	)
}
