package utils

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.Validator.Struct(i)
}

func MapValidationErrors(err error) map[string]string {
	errorsMap := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errorsMap[e.Field()] = e.Tag()
		}
	}
	return errorsMap
}
