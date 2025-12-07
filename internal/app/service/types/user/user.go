package user

import (
    `time`
    
    `dpcms/internal/app/helper/dbscope`
    `dpcms/internal/app/service/types`
    
    `github.com/go-playground/validator/v10`
)

type PasswordConfirm struct {
    Password        string `validate:"required,min=6,max=32" json:"password"  label:"密码"`
    PasswordConfirm string `validate:"required,eqfield=Password" json:"passwordConfirm" label:"确认密码"`
}

func (u *PasswordConfirm) ValidationMessage(e validator.FieldError) (string, bool) {
    key := e.StructField() + "." + e.ActualTag()
    message := map[string]string{
        "PasswordConfirm.eqfield": "两次密码输入不一致",
    }
    errMsg, ok := message[key]
    return errMsg, ok
}

type CreateDTO struct {
    PasswordConfirm
    Username string `validate:"required,alphanum,min=4,max=32" json:"username" label:"用户名"`
    Email    string `validate:"required,max=64,email" json:"email" label:"邮箱"`
    IP       string `json:"-"`
}

type UpdatePasswordDTO struct {
    PasswordConfirm
    RawPassword string `validate:"required" json:"rawPassword" label:"原密码"`
}

type CredentialDTO struct {
    Username string `validate:"required" json:"username" label:"用户名"`
    Password string `validate:"required" json:"password" label:"密码"`
}

type ResetPasswordDTO struct {
    ID       int64  `validate:"required" json:"ID" label:"用户ID"`
    Password string `validate:"required,min=6,max=32" json:"password" label:"密码"`
}

type ListDTO struct {
    dbscope.Pagination
    Username        *string    `json:"username" label:"用户名"`
    Status          *int8      `json:"status" label:"状态"`
    VerifyStartTime *time.Time `json:"verifyStartTime" label:"验证时间起始"`
    VerifyEndTime   *time.Time `json:"verifyEndTime" label:"验证时间结束"`
    Email           *string    `json:"email" label:"邮箱"`
    IP              *string    `json:"ip"`
    RoleID          *int64     `json:"roleId" label:"角色"`
    Nickname        *string    `json:"nickname" label:"昵称"`
    Gender          *int8      `validate:"omitnil,oneof=0 1" json:"gender" label:"性别"`
    Country         *string    `json:"country" label:"国家"`
    Province        *string    `json:"province" label:"省份"`
    City            *string    `json:"city" label:"城市"`
}

func (u *ListDTO) ValidationMessage(e validator.FieldError) (string, bool) {
    key := e.StructField() + "." + e.ActualTag()
    messages := map[string]string{
        "Gender.oneof": "性别必须是[男 女]的其中一个",
    }
    errMsg, ok := messages[key]
    return errMsg, ok
}

type ProfileDTO struct {
    Nickname    string `json:"nickname"`
    Gender      int8   `json:"gender"`
    Description string `json:"description"`
    Country     string `json:"country"`
    Province    string `json:"province"`
    City        string `json:"city"`
}

type AccountDTO struct {
    types.Meta
    Username string     `json:"username"`
    RoleID   int64      `json:"roleId"`
    RoleName string     `json:"roleName"`
    IP       string     `json:"ip"`
    Profile  ProfileDTO `json:"profile"`
}

type LoginResult struct {
    Detail AccountDTO `json:"detail"`
    Token  string     `json:"token"`
}

type DeleteDTO struct {
    ID int64 `validate:"required" json:"id"`
}

type UpdateProfileDTO struct {
    ID int64 `validate:"required" copier:"-" json:"id"`
    ProfileDTO
}
