package validator

import `github.com/go-playground/validator/v10`

type CustomValidationMessage interface {
    ValidationMessage(e validator.FieldError) (string, bool)
}
