package service

import (
	"github.com/go-playground/validator/v10"
	"gitlab.com/spacewalker/locations/internal/pkg/validation"
	"log"
)

var validate *validator.Validate

func init() {
	var err error

	validate = validator.New()

	if err = validate.RegisterValidation("validlongitude", validation.ValidateLongitude); err != nil {
		log.Panicf("failed to register validation: %v", err)
	}

	if err = validate.RegisterValidation("validlatitude", validation.ValidateLatitude); err != nil {
		log.Panicf("failed to register validation: %v", err)
	}

	if err = validate.RegisterValidation("validgeopoint", validation.ValidateGeoPoint); err != nil {
		log.Panicf("failed to register validation: %v", err)
	}
}
