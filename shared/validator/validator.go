package validator

import "github.com/go-playground/validator/v10"

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func NewValidator() *structValidator {
	return &structValidator{
		validate: validator.New(),
	}
}
