package validate

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Validator is the validator for the application
type Validator struct {
	validator *validator.Validate
}

func New() *Validator {
	v := validator.New()
	v.RegisterValidation("future", func(fl validator.FieldLevel) bool {
		t, ok := fl.Field().Interface().(time.Time)
		if !ok {
			return false
		}
		return t.After(time.Now())
	})

	return &Validator{validator: v}
}

// Validate validates the input
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
