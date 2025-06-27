package validator

import "github.com/go-playground/validator/v10"

type Interface interface {
	Validate(data interface{}) error
}

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{validator: validator.New()}
}

func (v *Validator) Validate(data interface{}) error {
	return v.validator.Struct(data)
}
