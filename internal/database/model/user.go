package model

import (
    `database/sql`
    
    `dpcms/internal/packages/database`
)

type User struct {
    database.Model
    Username   string       `gorm:"not null;uniqueIndex" json:"username"`
    Password   string       `gorm:"not null" json:"-"`
    Email      string       `gorm:"not null;uniqueIndex" json:"email"`
    VerifiedAt sql.NullTime `gorm:"default:null" json:"verifiedAt"`
    IP         string       `gorm:"not null;default:''" json:"ip"`
    Status     int8         `gorm:"not null;default:0" json:"status"`
    RoleID     uint64       `gorm:"index" json:"roleId"`
    Profile    *UserProfile `gorm:"foreignKey:UserID;reference:ID" json:"profile"`
}

type UserProfile struct {
    database.Model
    UserID      uint64 `gorm:"not null;uniqueIndex" json:"userId"`
    Nickname    string `json:"nickname"`
    Gender      int8   `gorm:"comment:0男性,1女性" json:"gender"`
    Description string `json:"description"`
    Country     string `json:"country"`
    Province    string `json:"province"`
    City        string `json:"city"`
}
