package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() (*Validator, error) {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return fld.Name
		}
		return name
	})

	v.RegisterAlias("mmYYYY", "datetime=01-2006")

	return &Validator{validate: v}, nil
}

func (v *Validator) Validate(s interface{}) error {
	return v.validate.Struct(s)
}

func (v *Validator) ValidateStruct(s interface{}) map[string][]string {
	if err := v.validate.Struct(s); err != nil {
		ve, ok := err.(validator.ValidationErrors)
		if !ok {
			return map[string][]string{"_error": {"validation failed"}}
		}
		out := make(map[string][]string)
		for _, e := range ve {
			key := e.Field()
			out[key] = append(out[key], getErrorMessage(e))
		}
		return out
	}
	return nil
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "uuid4":
		return fmt.Sprintf("%s must be a valid %s", err.Field(), err.Tag())
	case "datetime", "mmYYYY":
		return fmt.Sprintf("%s must match MM-YYYY", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}
