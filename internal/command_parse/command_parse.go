package command_parse

import (
	"fmt"
)

func ToCommandParameters(command string) []string {
	result := make([]string, 0)
	current := ""
	quote := 0
	lastQuote := ' '
	quoteMatchLast := false
	isQuote := false

	for i, char := range command {
		isQuote = char == '\'' || char == '"'
		quoteMatchLast = quote > 0 && isQuote && char == lastQuote

		if isQuote {
			lastQuote = char
		}

		if char == ' ' && quote == 0 {
			result = append(result, current)
			current = ""

			continue
		}

		if i == len(command)-1 {
			if !isQuote {
				current = fmt.Sprintf("%s%s", current, string(char))
			}
			result = append(result, current)
			current = ""

			continue
		}

		if isQuote {
			if !quoteMatchLast {
				quote = quote + 1

				continue
			}

			quote = quote - 1

			continue
		}

		current = fmt.Sprintf("%s%s", current, string(char))
	}

	return result
}
