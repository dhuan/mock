package tests_e2e

import (
	"fmt"
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
