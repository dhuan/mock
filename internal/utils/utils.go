package utils

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

func ReplaceRegex(subject string, find []string, replaceWith string) string {
	if len(find) == 0 {
		return subject
	}

	re := regexp.MustCompile(find[0])

	return ReplaceRegex(
		re.ReplaceAllString(subject, replaceWith),
		find[1:],
		replaceWith,
	)
}

func Unquote(subject string) string {
	return ReplaceRegex(subject, []string{`^"`, `"$`}, "")
}

func ListsEqual[T comparable](listA []T, listB []T) bool {
	for i := range listA {
		if listA[i] != listB[i] {
			return false
		}
	}

	return true
}

func ListsEqualUnsorted[T comparable](listA, listB []T) bool {
	if len(listA) != len(listB) {
		return false
	}

	for i := range listA {
		if IndexOf(listB, listA[i]) == -1 {
			return false
		}
	}

	return true
}

func IndexOf[T comparable](list []T, value T) int {
	for i := range list {
		if list[i] == value {
			return i
		}
	}

	return -1
}

func MktempDir() (string, error) {
	result, err := exec.Command("mktemp", "-d").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(result), "\n"), nil
}

func HasHeaderWithValue(headers *http.Header, headerKeyToSearch, headerValueToSearch string) bool {
	for headerKey, headerValues := range *headers {
		for _, headerValue := range headerValues {
			if headerKey == headerKeyToSearch && headerValue == headerValueToSearch {
				return true
			}
		}
	}

	return false
}

func BeginsWith(subject, find string) bool {
	return strings.Index(subject, find) == 0
}

func JoinMap[K comparable, V comparable](mapDst map[K]V, mapSrc map[K]V) {
	for i, v := range mapSrc {
		mapDst[i] = v
	}
}

func MapContains[K comparable, V comparable](m map[K]V, key K, value V) bool {
	valueExtracted, ok := m[key]
	if !ok {
		return false
	}

	return value == valueExtracted
}

func AnyEquals[T comparable](list []T, value T) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}

	return false
}

func MarshalJsonHelper[T comparable](
	mapValues map[T]string,
	errorMessage string,
	value *T,
) ([]byte, error) {
	for i := range mapValues {
		if *value == i {
			return []byte(mapValues[i]), nil
		}
	}

	return []byte(""), errors.New(fmt.Sprintf(errorMessage, *value))
}

func UnmarshalJsonHelper[T comparable](
	value *T,
	mapValues map[T]string,
	data []byte,
	errorMessage string,
) error {
	assertTypeText := Unquote(string(data))

	for key := range mapValues {
		if assertTypeText == mapValues[key] {
			*value = key

			return nil
		}
	}

	return errors.New(fmt.Sprintf(errorMessage, assertTypeText))
}

func MapMapValueOnly[T_Key comparable, T_Value comparable, T_ValueB interface{}](
	subject map[T_Key]T_Value,
	transform func(value T_Value) T_ValueB,
) map[T_Key]T_ValueB {
	result := make(map[T_Key]T_ValueB)

	for key := range subject {
		result[key] = transform(subject[key])
	}

	return result
}

func WrapIn(wrapper string) func(subject string) string {
	return func(subject string) string {
		return fmt.Sprintf(`"%s"`, subject)
	}
}

func GetKeys[T_Key comparable, T_Value interface{}](subject map[T_Key]T_Value) []T_Key {
	keys := make([]T_Key, 0, len(subject))

	for key := range subject {
		keys = append(keys, key)
	}

	return keys
}

func GetSortedKeys[T interface{}](subject map[string]T) []string {
	keys := GetKeys(subject)
	sort.Strings(keys)

	return keys
}

func ToCommandParams(command string) (string, []string) {
	commandSplit := strings.Split(command, " ")

	if len(commandSplit) == 0 {
		return "", []string{}
	}

	if len(commandSplit) == 1 {
		return commandSplit[0], []string{}
	}

	return commandSplit[0], commandSplit[1:]
}

func ParseEnv(env map[string]string) []string {
	result := make([]string, len(env))
	keys := GetKeys(env)

	for i, key := range keys {
		result[i] = fmt.Sprintf("%s=%s", key, env[key])
	}

	return result
}

func ToHeadersText(headers http.Header) string {
	textLines := make([]string, len(headers))
	headerKeys := GetSortedKeys(headers)

	for i, headerKey := range headerKeys {
		textLines[i] = fmt.Sprintf("%s: %s", headerKey, strings.Join(headers[headerKey], ","))
	}

	return strings.Join(textLines, "\n")
}

func CreateTempFile(content string) (string, error) {
	mktempResult, err := exec.Command("mktemp").Output()
	if err != nil {
		return "", err
	}
	fileName := strings.TrimSuffix(string(mktempResult), "\n")

	if err = os.WriteFile(fileName, []byte(content), 0644); err != nil {
		return "", err
	}

	return fileName, err
}

func RemoveEmptyLines(text string) string {
	return FilterLines(text, HasText)
}

func HasText(text string) bool {
	return strings.TrimSpace(text) != ""
}

func FilterLines(text string, filterFunc func(line string) bool) string {
	linesFiltered := make([]string, 0)
	lines := strings.Split(text, "\n")

	for i := range lines {
		if filterFunc(lines[i]) {
			linesFiltered = append(linesFiltered, lines[i])
		}
	}

	return strings.Join(linesFiltered, "\n")
}

func GetWord(index int, str, fallback string) string {
	splitResult := strings.Split(str, " ")

	if (len(splitResult) - 1) < index {
		return fallback
	}

	return splitResult[index]
}

func ReplaceVars(
	text string,
	vars map[string]string,
	toVarPlaceholder func(varName string) string,
) string {
	result := fmt.Sprintf("%s", text)

	for varName, varValue := range vars {
		currentSearch := toVarPlaceholder(varName)

		result = ReplaceRegex(result, []string{
			currentSearch,
		}, varValue)
	}

	return result
}

func ToDolarSignVariablePlaceHolder(varName string) string {
	return fmt.Sprintf("\\$%s", varName)
}

func ToDolarSignWithWrapVariablePlaceHolder(varName string) string {
	return fmt.Sprintf("\\${%s}", varName)
}

func ToCommandStrings(command string) []string {
    result := make([]string, 0, 0)
    current := ""

    for _, char := range command {
        if char == ' ' {
            result = append(result, current)
        }

        current = fmt.Sprintf("%s%s", current, string(char))
    }

    return result
}
