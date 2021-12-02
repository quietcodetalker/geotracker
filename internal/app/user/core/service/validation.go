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

	if err = validate.RegisterValidation("validusername", validation.ValidateUsername); err != nil {
		log.Panicf("failed to register validation: %v", err)
	}
}
