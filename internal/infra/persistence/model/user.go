package model

import (
    `database/sql`
    
    `dpcms/internal/infra/persistence/datatype`
)

type User struct {
    Base
    Username   string              `gorm:"not null;uniqueIndex" json:"username"`
    Password   string              `gorm:"not null" json:"-"`
    Email      string              `gorm:"not null;uniqueIndex" json:"email"`
    VerifiedAt sql.NullTime        `gorm:"default:null" json:"-"`
    IP         string              `gorm:"not null;default:''" json:"ip"`
    Status     int8                `gorm:"not null;default:0" json:"status"`
    RoleID     datatype.SafeUint64 `gorm:"index" json:"roleId"`
    Profile    *UserProfile        `gorm:"foreignKey:UserID;reference:ID" json:"profile"`
}

type UserProfile struct {
    Base        `json:"-"`
    UserID      datatype.SafeUint64 `gorm:"not null;uniqueIndex" json:"userId"`
    Nickname    string              `json:"nickname"`
    Gender      int8                `gorm:"comment:0男性,1女性" json:"gender"`
    Description string              `json:"description"`
    Country     string              `json:"country"`
    Province    string              `json:"province"`
    City        string              `json:"city"`
}
