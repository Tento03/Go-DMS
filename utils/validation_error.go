package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidationError(err error) map[string]string {
	errors := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())
		errors[field] = err.Tag()
	}

	return errors
}
