package app

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("nowhitespace", noWhitespace)

	return validate
}

func noWhitespace(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	return !regexp.MustCompile(`\s`).MatchString(value)
}
