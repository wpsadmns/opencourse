package stringutil

import "strings"

const (
	EmptyString = ""
)

func IsNotEmpty(value string) bool {
	return len(strings.TrimSpace(value)) != 0
}

func IsEmpty(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}
