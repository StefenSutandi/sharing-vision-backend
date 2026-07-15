package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			switch err.Tag() {
			case "required":
				errors[field] = err.Field() + " is required"
			case "min":
				errors[field] = err.Field() + " must contain at least " + err.Param() + " characters"
			case "max":
				errors[field] = err.Field() + " must contain at most " + err.Param() + " characters"
			case "oneof":
				errors[field] = err.Field() + " must be one of: " + err.Param()
			default:
				errors[field] = err.Field() + " is invalid"
			}
		}
		return errors
	}
	return nil
}
