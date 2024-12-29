package map_extract

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func parsePathItem(path string) []string {
	if !regexTest(`\[[0-9{1,}]\]$`, path) || regexTest(`^\[[0-9{1,}]\]$`, path) {
		return []string{path}
	}

	words := strings.Split(path, "[")

	return []string{words[0], fmt.Sprintf("[%s", words[1])}
}

func regexTest(regex string, subject string) bool {
	match, err := regexp.MatchString(regex, subject)

	if err != nil {
		return false
	}

	return match
}

func extractDigits(str string) int {
	result := ""
	for _, c := range str {
		if regexTest("[0-9]", string(c)) {
			result = fmt.Sprintf("%s%s", result, string(c))
		}
	}

	num, err := strconv.Atoi(result)
	if err != nil {
		panic(err)
	}

	return num
}

func isArrayPath(path string) (int, bool) {
	if !regexTest(`\[[0-9{1,}]\]$`, path) {
		return -1, false
	}

	return extractDigits(path), true
}

func Extract(data interface{}, path string) (interface{}, bool) {
	paths := make([]string, 0)
	for _, pathItem := range strings.Split(path, ".") {
		newPaths := parsePathItem(pathItem)

		paths = append(paths, newPaths...)
	}

	var value interface{} = data

	for len(paths) > 0 {
		current := paths[0]

		if index, ok := isArrayPath(current); ok {
			array, ok := isArray(value)
			if !ok {
				return nil, false
			}

			if index > (len(array) - 1) {
				return nil, false
			}

			newValue := array[index]
			if newValue == nil {
				return nil, false
			}

			value = newValue

			paths = paths[1:]

			continue
		}

		obj, ok := value.(map[string]interface{})
		if !ok {
			return nil, false
		}

		newValue, ok := obj[current]
		if !ok {
			return nil, false
		}

		value = newValue

		paths = paths[1:]
	}

	kind := reflect.TypeOf(value).Kind()
	if kind == reflect.Map || kind == reflect.Array || kind == reflect.Slice {
		return toJsonString(value), true
	}

	return value, true
}

func toJsonString(value interface{}) string {
	result, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	return string(result)
}

func isArray(value interface{}) ([]interface{}, bool) {
	parsed, ok := value.([]interface{})

	return parsed, ok
}
