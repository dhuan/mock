package tests_e2e

import (
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
		"",
		StringMatches("Hello world!\n"),
	)
}
