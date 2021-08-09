package util

import (
	"regexp"
	"strings"
)

func SnakeCase(str string) string {
	snake := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(str, "${1}_${2}")
	snake = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
