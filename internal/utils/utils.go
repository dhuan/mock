package utils

import (
	"errors"
	"fmt"
	"net/http"
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
	for i, _ := range listA {
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

	for i, _ := range listA {
		if IndexOf[T](listB, listA[i]) == -1 {
			return false
		}
	}

	return true
}

func IndexOf[T comparable](list []T, value T) int {
	for i, _ := range list {
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
	for i, _ := range list {
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
	for i, _ := range mapValues {
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

	for key, _ := range mapValues {
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

	for key, _ := range subject {
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

	for key, _ := range subject {
		keys = append(keys, key)
	}

	return keys
}

func GetSortedKeys[T interface{}](subject map[string]T) []string {
	keys := GetKeys[string, T](subject)
	sort.Strings(keys)

	return keys
}
