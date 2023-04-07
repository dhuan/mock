package command_parse_test

import (
	"testing"

	"github.com/dhuan/mock/internal/command_parse"
	"github.com/stretchr/testify/assert"
)

func Test_SimpleCommand(t *testing.T) {
	assert.Equal(
		t,
		[]string{
			"echo",
			"Hello",
			"world",
		},
		command_parse.ToCommandParameters("echo Hello world"),
	)
}

func Test_WithQuotes(t *testing.T) {
	assert.Equal(
		t,
		[]string{
			"echo",
			"Hello world",
		},
		command_parse.ToCommandParameters("echo 'Hello world'"),
	)
}

func Test_WithQuotesMultipleTimes(t *testing.T) {
	assert.Equal(
		t,
		[]string{
			"command",
			"Hello world",
			"foo bar",
		},
		command_parse.ToCommandParameters("command 'Hello world' 'foo bar'"),
	)
}

func Test_Foobar(t *testing.T) {
	assert.Equal(
		t,
		[]string{
			"serve",
			"--param-one",
			"value one",
			"--param-two",
			"value two",
			"--param-three",
			"value three",
		},
		command_parse.ToCommandParameters("serve --param-one 'value one' --param-two 'value two' --param-three 'value three'"),
	)
}
