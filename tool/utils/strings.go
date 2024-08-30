package utils

import (
	"regexp"
	"strings"
)

func PascalCaseToSnakeCase(s string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(s, "${1}_${2}")

	return strings.ToLower(snake)
}
