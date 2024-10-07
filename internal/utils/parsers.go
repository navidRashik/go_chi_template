package utils

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type IError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

func ParseValidationErrors(err error) []IError {
	var errors []IError
	for _, err := range err.(validator.ValidationErrors) {
		var el IError
		el.Field = toSnakeCase(err.Field())
		el.Tag = err.ActualTag()
		el.Value = err.Param()
		errors = append(errors, el)
	}
	return errors
}

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func init() {
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
}

var matchFirstCap, matchAllCap *regexp.Regexp
