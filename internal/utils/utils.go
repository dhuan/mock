package utils

import (
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
