package dto

import (
    `time`
    
    "github.com/golang-jwt/jwt/v5"
)

type User struct {
    ID              *uint   `json:"id"`
    Username        *string `validate:"required,alphanum,min=6,max=32" json:"username" label:"用户名"`
    Password        *string `validate:"required,min=6,max=32" json:"password"  label:"密码"`
    ConfirmPassword *string `validate:"required,eqfield=Password" json:"confirmPassword" label:"确认密码"`
    Email           *string `validate:"required,max=64,email" json:"email" label:"邮箱"`
    RoleID          *uint   `json:"roleID"`
    Status          *uint   `json:"status"`
}

type UserToken struct {
    jwt.RegisteredClaims
    UserID uint `json:"userID"`
}

func (u *UserToken) Create(key string) (string, error) {
    if u.ExpiresAt == nil {
        u.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 3))
    }
    return jwt.NewWithClaims(jwt.SigningMethodHS512, u).SignedString([]byte(key))
}
