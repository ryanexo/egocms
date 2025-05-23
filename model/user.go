package model

import (
    `database/sql`
    
    `dpcms/packages/database`
)

type User struct {
    database.Model
    Username   string       `gorm:"type:varchar(255);not null;uniqueIndex" json:"username"`
    Password   string       `gorm:"type:varchar(255);not null" json:"-"`
    Email      string       `gorm:"type:varchar(64);not null;uniqueIndex" json:"email"`
    VerifiedAt sql.NullTime `gorm:"default:null" json:"verifiedAt"`
    IP         string       `gorm:"type:varchar(255);not null;default:''" json:"ip"`
    Status     int8         `gorm:"not null;default:0" json:"status"`
}
