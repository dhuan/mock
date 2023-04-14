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

func Test_WithQuotesMultipleTimes_2(t *testing.T) {
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

func Test_WithSubQuotes(t *testing.T) {
	assert.Equal(
		t,
		[]string{
			"serve",
			"--param-one",
			`here "goes some" text`,
		},
		command_parse.ToCommandParameters(`serve --param-one 'here "goes some" text'`),
	)
}

func Test_WithSubQuotes_2(t *testing.T) {
	assert.Equal(
		t,
		[]string{
			"serve",
			"--param-one",
			`here "goes some" text`,
			"--param-two",
			`here "goes another" text`,
		},
		command_parse.ToCommandParameters(`serve --param-one 'here "goes some" text' --param-two 'here "goes another" text'`),
	)
}
