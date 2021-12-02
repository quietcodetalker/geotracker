package validation

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var usernameRegexp = regexp.MustCompile("[a-zA-Z0-9]{4,16}")

func ValidateUsername(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return usernameRegexp.Match([]byte(value))
}
