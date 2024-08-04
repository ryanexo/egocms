package dto

import (
    `time`
    
    `dpcms/packages/validate`
    "github.com/golang-jwt/jwt/v5"
)

type User struct {
    validate.PartialValidate `validate:"-"`
    ID                       *uint   `validate:"required,gte=1" json:"id"`
    Username                 *string `validate:"required,alphanum,min=6,max=32" json:"username"`
    Password                 *string `validate:"required,min=6,max=32" json:"password"`
    ConfirmPassword          *string `validate:"required,eqfield=Password" json:"confirmPassword"`
    Email                    *string `validate:"required,max=64,email" json:"email"`
    RoleID                   *uint   `validate:"required,gte=1" json:"roleID"`
    Status                   *uint   `validate:"required,gte=1" json:"status"`
}

func (u *User) WithLoginScene() {
    u.SetValidationFields("Username", "Password")
}

func (u *User) WithRegisterScene() {
    u.SetValidationFields("Username", "Password", "ConfirmPassword", "Email")
}

func (u *User) WithResetPassword() {
    u.SetValidationFields("Username", "Password", "ConfirmPassword")
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
