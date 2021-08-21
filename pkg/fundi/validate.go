package fundi

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type (
	// InputValidator defines Fails and Errors method.
	InputValidator interface {
		Fails(input interface{}) bool
		Errors() map[string]string
	}

	// InputValidation type implements InputValidator interface.
	InputValidation struct {
		validate *validator.Validate
		errors   error
		messages map[string]string
	}
)

// NewInputValidator returns an instance of InputValidation.
func NewInputValidator() *InputValidation {
	v := validator.New()

	// register function to get field name from json tags.
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// todo: add method to register custom validation tags

	return &InputValidation{
		validate: v,
		messages: map[string]string{
			"required": "%s is a required field",
		},
	}
}

// Fails returns true if the InputValidation struct has invalid data.
func (i *InputValidation) Fails(input interface{}) bool {
	i.errors = i.validate.Struct(input)
	return i.errors != nil
}

// Errors returns a map of error strings per field.
func (i *InputValidation) Errors() map[string]string {
	errs := make(map[string]string)

	for _, err := range i.errors.(validator.ValidationErrors) {
		field := err.Field()
		errs[field] = fmt.Sprintf(i.messages[err.Tag()], field)
	}

	return errs
}

// RegisterValidationTag adds a custom validation tag.
func (i InputValidation) RegisterValidationTag(tag string, fn func()) error {
	// todo: wrap fn and return expected fn
	return i.validate.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
		// todo
		return true
	})
}

// RegisterValidationMessages appends custom messages to the Validator.
func (i *InputValidation) RegisterValidationMessages(map[string]string) error {
	// todo: append validation messages
	return nil
}
