package utils

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
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

func IndexOfRegex(list []string, subject string) int {
	for i := range list {
		if RegexTest(list[i], subject) {
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

func EndsWith(subject, find string) bool {
	lenSubject := len(subject)
	lenFind := len(find)

	return subject[lenSubject-lenFind:] == find
}

func JoinMap[K comparable, V comparable](mapDst map[K]V, mapSrc map[K]V) {
	for i, v := range mapSrc {
		mapDst[i] = v
	}
}

func MapGetKeyByValue[K comparable, V comparable](m map[K]V, search V, fallback K) K {
	keys := GetKeys(m)

	for _, key := range keys {
		if m[key] == search {
			return key
		}
	}

	return fallback
}

func MapContains[K comparable, V comparable](m map[K]V, key K, value V) bool {
	valueExtracted, ok := m[key]
	if !ok {
		return false
	}

	return value == valueExtracted
}

func MapContainsX[K comparable, V comparable](m map[K]V, key K, value, fallback V) (bool, bool, V) {
	valueExtracted, ok := m[key]
	if !ok {
		return false, false, fallback
	}

	if value == valueExtracted {
		return true, true, valueExtracted
	}

	return true, false, valueExtracted
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

	return []byte(""), fmt.Errorf(errorMessage, *value)
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

	return fmt.Errorf(errorMessage, assertTypeText)
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

func ToCommandParams(commandSplit []string) (string, []string) {
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

func CreateTempFile(content []byte) (string, error) {
	mktempResult, err := exec.Command("mktemp").Output()
	if err != nil {
		return "", err
	}
	fileName := strings.TrimSuffix(string(mktempResult), "\n")

	if err = os.WriteFile(fileName, content, 0644); err != nil {
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
	result := text

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
	result := make([]string, 0)
	current := ""
	quote := 0
	lastQuote := ' '
	quoteMatchLast := false
	isQuote := false

	for i, char := range command {
		isQuote = char == '\'' || char == '"'
		quoteMatchLast = isQuote && char == lastQuote

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

func ExtractNumbersFromString(str string) (int, error) {
	resultStr := ""

	for _, char := range str {
		if !charIsNumber(string(char)) {
			continue
		}

		resultStr = fmt.Sprintf("%s%s", resultStr, string(char))
	}

	return strconv.Atoi(resultStr)
}

var number_chars = "012346789"

func charIsNumber(char string) bool {
	return strings.Contains(number_chars, char)
}

func RegexTest(regex string, subject string) bool {
	match, err := regexp.MatchString(regex, subject)

	if err != nil {
		return false
	}

	return match
}

func ParseHeaderLine(text string) (string, string, bool) {
	splitResult := strings.Split(text, ":")

	if len(splitResult) < 2 {
		return "", "", false
	}

	return splitResult[0], strings.Join(splitResult[1:], ":"), true
}

func ExtractHeadersFromText(fileContent []byte) map[string]string {
	headers := make(map[string]string)

	fileContentText := RemoveEmptyLines(string(fileContent))

	if fileContentText == "" {
		return headers
	}

	headerLines := strings.Split(fileContentText, "\n")

	for i := range headerLines {
		headerKey, headerValue, ok := ParseHeaderLine(headerLines[i])
		if !ok {
			continue
		}

		headers[headerKey] = headerValue
	}

	return headers
}

func IsPortFree(portNumber int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	free := err == nil

	if err == nil {
		listener.Close()
	}

	return free
}

func GetFreePort() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	listener.Close()

	return listener.Addr().(*net.TCPAddr).Port
}

func MapFilterFileLines(filePath string, f func(string) (string, bool)) error {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	result := make([]string, 0)
	lines := strings.Split(string(fileContent), "\n")

	for i := range lines {
		modifiedLine, add := f(lines[i])
		if !add {
			continue
		}

		result = append(result, modifiedLine)
	}

	return os.WriteFile(filePath, []byte(strings.Join(result, "\n")), 0644)
}

func AddLineToFile(filePath string, newLine string) error {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(fileContent), "\n")
	lines = append(lines, newLine)

	return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
}
