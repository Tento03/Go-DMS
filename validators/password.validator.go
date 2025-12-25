package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSymbol := regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)

	return len(password) >= 8 &&
		hasUpper &&
		hasLower &&
		hasNumber &&
		hasSymbol
}
