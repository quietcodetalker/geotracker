package service

import (
	"log"

	"github.com/go-playground/validator/v10"
	"gitlab.com/spacewalker/geotracker/internal/app/location/core/port"
	"gitlab.com/spacewalker/geotracker/internal/pkg/validation"
)

var validate *validator.Validate

func init() {
	var err error

	validate = validator.New()

	if err = validate.RegisterValidation("validusername", validation.ValidateUsername); err != nil {
		log.Panicf("failed to register validation: %v", err)
	}

	if err = validate.RegisterValidation("validlongitude", validation.ValidateLongitude); err != nil {
		log.Panicf("failed to register validation: %v", err)
	}

	if err = validate.RegisterValidation("validlatitude", validation.ValidateLatitude); err != nil {
		log.Panicf("failed to register validation: %v", err)
	}

	if err = validate.RegisterValidation("validgeopoint", validation.ValidateGeoPoint); err != nil {
		log.Panicf("failed to register validation: %v", err)
	}

	if err = validate.RegisterValidation("valid_page_token", validation.ValidatePageToken); err != nil {
		log.Panicf("failed to register validation: %v", err)
	}

	validate.RegisterStructValidation(validation.ValidateListMethod, port.UserServiceListUsersInRadiusRequest{})
}
