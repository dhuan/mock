package utils_test

import (
	"github.com/dhuan/mock/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ToCommandStrings_Simple(t *testing.T) {
	assert.Equal(
		t,
		[]string{
			"echo",
			"foo",
			"bar",
		},
		utils.ToCommandStrings("echo foo bar"),
	)
}

func Test_ToCommandStrings_WithQuotes(t *testing.T) {
	type testCase struct {
		command  string
		expected []string
	}

    tcs := []testCase{
		{"echo 'foo bar'", []string{"echo", "foo bar"}},
		{"echo \"foo bar\"", []string{"echo", "foo bar"}},
		{"echo 'foo' bar", []string{"echo", "foo", "bar"}},
		{"echo \"foo\" bar", []string{"echo", "foo", "bar"}},
	}

    for _, tc := range tcs {
        assert.Equal(
            t,
            tc.expected,
            utils.ToCommandStrings(tc.command),
        )
    }
}
