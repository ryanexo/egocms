package srvparams

import (
    `database/sql`
    `time`
    
    `dpcms/internal/app/helper/dbscope`
    `dpcms/internal/packages/database`
    
    `github.com/go-playground/validator/v10`
)

type UserPasswdConfirm struct {
    Password        string `validate:"required,min=6,max=32" json:"password"  label:"密码"`
    PasswordConfirm string `validate:"required,eqfield=Password" json:"passwordConfirm" label:"确认密码"`
}

func (u *UserPasswdConfirm) ValidationMessage(e validator.FieldError) (string, bool) {
    key := e.StructField() + "." + e.ActualTag()
    message := map[string]string{
        "UserPasswdConfirm.eqfield": "两次密码输入不一致",
    }
    errMsg, ok := message[key]
    return errMsg, ok
}

type UserCreateParams struct {
    UserPasswdConfirm
    Username string `validate:"required,alphanum,min=4,max=32" json:"username" label:"用户名"`
    Email    string `validate:"required,max=64,email" json:"email" label:"邮箱"`
    IP       string `json:"-"`
}

type UserPasswdUpdateParams struct {
    UserPasswdConfirm
    RawPassword string `validate:"required" json:"rawPassword" label:"原密码"`
}

type UserCredentialParams struct {
    Username string `validate:"required" json:"username" label:"用户名"`
    Password string `validate:"required" json:"password" label:"密码"`
}

type UserPasswdResetParams struct {
    ID       uint64 `validate:"required" json:"ID" label:"用户ID"`
    Password string `validate:"required,min=6,max=32" json:"password" label:"密码"`
}

type UserListQueryParams struct {
    dbscope.Pagination
    Username        *string    `json:"username" label:"用户名"`
    Status          *int8      `json:"status" label:"状态"`
    VerifyStartTime *time.Time `json:"verifyStartTime" label:"验证时间起始"`
    VerifyEndTime   *time.Time `json:"verifyEndTime" label:"验证时间结束"`
    Email           *string    `json:"email" label:"邮箱"`
    IP              *string    `json:"ip"`
    RoleID          *uint64    `json:"roleId" label:"角色"`
    Nickname        *string    `json:"nickname" label:"昵称"`
    Gender          *int8      `validate:"omitnil,oneof=0 1" json:"gender" label:"性别"`
    Country         *string    `json:"country" label:"国家"`
    Province        *string    `json:"province" label:"省份"`
    City            *string    `json:"city" label:"城市"`
}

func (u *UserListQueryParams) ValidationMessage(e validator.FieldError) (string, bool) {
    key := e.StructField() + "." + e.ActualTag()
    messages := map[string]string{
        "Gender.oneof": "性别必须是[男 女]的其中一个",
    }
    errMsg, ok := messages[key]
    return errMsg, ok
}

type UserProfile struct {
    Nickname    string `json:"nickname"`
    Gender      int8   `json:"gender"`
    Description string `json:"description"`
    Country     string `json:"country"`
    Province    string `json:"province"`
    City        string `json:"city"`
}

type User struct {
    database.Model
    Username   string       `json:"username"`
    Password   string       `json:"-"`
    Email      string       `json:"email"`
    VerifiedAt sql.NullTime `json:"verifiedAt"`
    IP         string       `json:"ip"`
    Status     int8         `json:"status"`
    RoleID     uint64       `json:"roleId"`
    RoleName   string       `json:"roleName"`
    Profile    *UserProfile `json:"profile"`
}

type UserAuthnResult struct {
    User  User   `json:"detail"`
    Token string `json:"token"`
}

type UserProfileUpdateParams struct {
    ID uint64 `validate:"required" json:"id"`
    UserProfile
}
