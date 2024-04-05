package entity

import (
    `time`

    `gorm.io/gorm`
)

type User struct {
    gorm.Model
    LoginLog    UserLoginLog
    Group       Group
    Username    string    `gorm:"type:varchar(32);not null;uniqueIndex"`
    Password    string    `gorm:"type:varchar(255);not null"`
    Email       string    `gorm:"type:varchar(64);not null;uniqueIndex"`
    RoleID      uint64    `gorm:"not null;default:0" json:"role_id"`
    GroupID     uint64    `gorm:"not null;default:0" json:"group_id"`
    Points      uint64    `gorm:"not null;default:0"`
    Experiences uint64    `gorm:"not null;default:0"`
    Signature   string    `gorm:"not null;default:''"`
    VerifiedAt  time.Time `gorm:"not null;default:'0000-00-00 00:00:00'" json:"-"`
    CreatedIP   string    `gorm:"type:varchar(39);not null;default:''" json:"-"`
    LastIP      string    `gorm:"type:varchar(39);not null;default:''"`
    LastAt      time.Time `gorm:"not null;default:'0000-00-00 00:00:00'"`
    Status      byte      `gorm:"not null"`
}

type UserCredential struct {
    ID       uint
    Username string `validate:"required|alphaNum/isAlphaNum|minLen:6|maxLen:32" form:"username"`
    Password string `validate:"required|minLen:6|maxLen:32" form:"password"`
}

type UserCreatingForm struct {
    Username   string `validate:"required|alphaNum/isAlphaNum|minLen:6|maxLen:32" form:"username"`
    Password   string `validate:"required|minLen:6|maxLen:32" form:"password"`
    RePassword string `validate:"required|minLen:6|maxLen:32" form:"re_password"`
    Email      string `validate:"required|maxLen:64|email"`
}
