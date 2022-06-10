package utils

import (
	"net/http"
	"os/exec"
	"regexp"
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
