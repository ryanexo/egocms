package api

import (
	"time"

	"cms/internal/public/apitype"
	"cms/internal/public/jsontype"

	"github.com/go-playground/validator/v10"
)

type PasswordConfirmation struct {
	Password        string `json:"password" validate:"required,min=6,max=32" label:"密码"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,min=6,max=32" label:"确认密码"`
}

type UserCreateParams struct {
	PasswordConfirmation
	Username string `json:"username" validate:"required,alphanum,min=4,max=255" label:"用户名"`
	Email    string `json:"email" validate:"required,max=255,email" label:"邮箱"`
	IP       string `json:"-"`
}

type UserPasswordUpdateParams struct {
	PasswordConfirmation
	OldPassword string `json:"oldPassword" validate:"required" label:"原密码"`
}

type UserCredentialParams struct {
	Username string `json:"username" validate:"required" label:"用户名"`
	Password string `json:"password" validate:"required" label:"密码"`
}

type UserPasswordResetParams struct {
	ID       jsontype.SafeUint64 `json:"id" validate:"required" label:"用户ID" apitype:"string"`
	Password string              `json:"password" validate:"required,min=6,max=32" label:"密码"`
}

type UserListParams struct {
	apitype.Pagination
	Username        *string              `json:"username" label:"用户名"`
	Status          *int8                `json:"status" label:"状态"`
	VerifyStartTime *time.Time           `json:"verifyStartTime" label:"验证时间起始"`
	VerifyEndTime   *time.Time           `json:"verifyEndTime" label:"验证时间结束"`
	Email           *string              `json:"email" label:"邮箱"`
	IP              *string              `json:"ip"`
	RoleID          *jsontype.SafeUint64 `json:"roleId" label:"角色" apitype:"string"`
	Nickname        *string              `json:"nickname" label:"昵称"`
	Gender          *int8                `json:"gender" validate:"omitnil,oneof=0 1 2" label:"性别"`
	Country         *string              `json:"country" label:"国家"`
	Province        *string              `json:"province" label:"省份"`
	City            *string              `json:"city" label:"城市"`
}

func (params *UserListParams) ValidationMessage(fieldError validator.FieldError) (string, bool) {
	messages := map[string]string{
		"Gender.oneof": "性别必须是未知、男、女的其中一个",
	}
	message, ok := messages[fieldError.StructField()+"."+fieldError.ActualTag()]
	return message, ok
}

type UserProfileParams struct {
	UserID      jsontype.SafeUint64 `json:"-" swaggerignore:"true"`
	Avatar      string              `json:"avatar" validate:"max=255"`
	Nickname    string              `json:"nickname" validate:"max=255"`
	Gender      int8                `json:"gender" validate:"oneof=0 1 2" label:"性别"`
	Description string              `json:"description" validate:"max=255"`
	Country     string              `json:"country" validate:"max=255"`
	Province    string              `json:"province" validate:"max=255"`
	City        string              `json:"city" validate:"max=255"`
}

type UserProfile struct {
	Avatar      string `json:"avatar"`
	Nickname    string `json:"nickname"`
	Gender      int8   `json:"gender"`
	Description string `json:"description"`
	Country     string `json:"country"`
	Province    string `json:"province"`
	City        string `json:"city"`
}

type User struct {
	apitype.Base
	Username   string               `json:"username"`
	Email      string               `json:"email"`
	VerifiedAt *time.Time           `json:"verifiedAt"`
	IP         *string              `json:"ip"`
	Status     int8                 `json:"status"`
	RoleID     *jsontype.SafeUint64 `json:"roleId" apitype:"string"`
	RoleName   string               `json:"roleName"`
	Profile    *UserProfile         `json:"profile"`
}

type UserAuthnResult struct {
	User  *User  `json:"detail"`
	Token string `json:"token"`
}
