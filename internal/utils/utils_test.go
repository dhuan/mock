package utils_test

import (
	"testing"
	"github.com/dhuan/mock/internal/utils"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateEndpointConfigs_Duplicates(t *testing.T) {

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
