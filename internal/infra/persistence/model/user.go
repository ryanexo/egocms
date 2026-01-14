package model

import (
    "database/sql"
    
    "cms/internal/infra/persistence/datatype"
)

type User struct {
    Base
    Username   string              `gorm:"not null;uniqueIndex"`
    Password   string              `gorm:"not null"`
    Email      string              `gorm:"not null;uniqueIndex"`
    VerifiedAt sql.NullTime        `gorm:"default:null"`
    IP         string              `gorm:"not null;default:''"`
    Status     int8                `gorm:"not null;default:0"`
    RoleID     datatype.SafeUint64 `gorm:"index"`
    Role       *Role               `gorm:"foreignKey:RoleID;reference:ID"`
    Profile    *UserProfile        `gorm:"foreignKey:UserID;reference:ID"`
}

type UserProfile struct {
    Base
    Avatar      string              `gorm:"type:varchar(255)"`
    UserID      datatype.SafeUint64 `gorm:"not null;uniqueIndex"`
    Nickname    string              `gorm:"type:varchar(255)"`
    Gender      int8                `gorm:"comment:0男性,1女性"`
    Description string              `gorm:"type:varchar(255)"`
    Country     string              `gorm:"type:varchar(255)"`
    Province    string              `gorm:"type:varchar(255)"`
    City        string              `gorm:"type:varchar(255)"`
}
