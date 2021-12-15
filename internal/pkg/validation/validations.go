package validation

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var usernameRegexp = regexp.MustCompile("^[a-zA-Z0-9]{4,16}$")

func ValidateUsername(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return usernameRegexp.Match([]byte(value))
}

func ValidateLongitude(fl validator.FieldLevel) bool {
	value := fl.Field().Float()
	return -180.0 <= value && value <= 180.0
}
func ValidateLatitude(fl validator.FieldLevel) bool {
	value := fl.Field().Float()
	return -90.0 <= value && value <= 90.0
}
