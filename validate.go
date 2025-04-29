package echoext

import "github.com/go-playground/validator/v10"

type Validator struct {
	v *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.v.Struct(i); err != nil {
		return err
	}

	return nil
}
