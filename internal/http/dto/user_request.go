package dto

import (
    `dpcms/internal/packages/validate`
    `github.com/go-playground/validator/v10`
    `github.com/iancoleman/strcase`
)

var _ validate.CustomValidationMessage = new(UserRegisterRequest)

type UserRegisterRequest struct {
    ID              *int64  `json:"id"`
    Username        *string `validate:"required,alphanum,min=4,max=32" json:"username" label:"用户名"`
    Password        *string `validate:"required,min=6,max=32" json:"password"  label:"密码"`
    PasswordConfirm *string `validate:"required,eqfield=Password" json:"passwordConfirm" label:"确认密码"`
    Email           *string `validate:"required,max=64,email" json:"email" label:"邮箱"`
    IP              string  `json:"-"`
}

func (u *UserRegisterRequest) ValidationMessage(e validator.FieldError) (string, bool) {
    key := strcase.ToLowerCamel(e.StructField() + "." + e.ActualTag())
    message := map[string]string{
        "PasswordConfirm.eqfield": "两次密码输入不一致",
    }
    errMsg, ok := message[key]
    return errMsg, ok
}

type UserLoginRequest struct {
    Username *string `validate:"required" json:"username"`
    Password *string `validate:"required" json:"password"`
}
