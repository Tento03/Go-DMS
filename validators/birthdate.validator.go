package validators

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func BirthDateValidator(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}
	return date.Before(time.Now())
}
