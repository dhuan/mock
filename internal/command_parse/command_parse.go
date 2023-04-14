package command_parse

import (
	"fmt"
)

func isFakeQuote(isQuoteChar bool, quoteCount int, char rune, quoteStack []rune) bool {
	if !isQuoteChar {
		return false
	}

	if quoteCount == 0 {
		return false
	}

	if len(quoteStack) == 0 {
		return false
	}

	lastIndex := len(quoteStack) - 1
	firstQuote := quoteStack[lastIndex]

	if quoteCount > 0 && char == firstQuote {
		return false
	}

	return true
}

func ToCommandParameters(command string) []string {
	result := make([]string, 0)
	currentWord := ""
	quoteCount := 0
	isSpace := false
	isQuote := false
	quoteStack := make([]rune, 0)
	isOpeningQuote := false

	for _, char := range command {
		isSpace = char == ' '
		isQuote = char == '\'' || char == '"'
		if isQuote {
			if len(quoteStack) == 0 {
				quoteStack = make([]rune, 0)
				quoteStack = append(quoteStack, char)
			} else {
				quoteStack = append([]rune{char}, quoteStack...)
			}
		}
		// quoteEqualsLast = isQuote && len(quoteStack) > 1 && quoteStack[1] == char
		if isFakeQuote(isQuote, quoteCount, char, quoteStack) {
			isQuote = false
		}
		isOpeningQuote = false
		if isQuote && quoteCount == 0 {
			isOpeningQuote = true
		}
		if isQuote && isOpeningQuote {
			quoteCount = quoteCount + 1
		}
		if isQuote && !isOpeningQuote {
			quoteCount = quoteCount - 1
		}
		if isQuote {
			if isOpeningQuote {
				if len(quoteStack) == 0 {
					quoteStack = make([]rune, 0)
					quoteStack = append(quoteStack, char)
				} else {
					quoteStack = append([]rune{char}, quoteStack...)
				}
			} else {
				quoteStack = quoteStack[1:]
			}
		}
		//////////////////////////////////////////

		if isQuote && isOpeningQuote {
			continue
		}

		if isQuote && !isOpeningQuote && quoteCount == 0 {
			result = append(result, currentWord)
			currentWord = ""

			continue
		}

		if isSpace && quoteCount == 0 {
			if len(currentWord) > 0 {
				result = append(result, currentWord)
				currentWord = ""
			}

			continue
		}

		currentWord = fmt.Sprintf("%s%s", currentWord, string(char))
	}

	if len(currentWord) > 0 {
		result = append(result, currentWord)
	}

	return result
}
