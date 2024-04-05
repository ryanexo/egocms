package validator

import (
    `github.com/gin-gonic/gin/binding`
    `github.com/gookit/validate`
)

type CustomValidator struct{}

func (CustomValidator) Engine() any {
    return nil
}

func (CustomValidator) ValidateStruct(data any) error {
    v := validate.Struct(data)
    v.Validate()
    return v.Errors
}

var _ binding.StructValidator = new(CustomValidator)
