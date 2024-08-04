package model

import (
    "time"
)

type User struct {
    Model
    Username   string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"username"`
    Password   string    `gorm:"type:varchar(255);not null" json:"-"`
    Email      string    `gorm:"type:varchar(64);not null;uniqueIndex" json:"email"`
    RoleID     uint      `gorm:"not null;default:0" json:"roleID"`
    VerifiedAt time.Time `gorm:"default:null" json:"verifiedAt"`
    IP         string    `gorm:"type:varchar(255);not null;default:''" json:"ip"`
    Status     uint      `gorm:"not null;default:0" json:"status"`
}

type Role struct {
    Model
    Name        string `gorm:"type:varchar(64);not null" json:"name"`
    Description string `gorm:"type:varchar(255);default:'';not null" json:"description"`
}

type RoleUserRelation struct {
    Model
    RoleID uint `gorm:"not null"`
    UserID uint `gorm:"not null"`
}

type Menu struct {
    Model
    ParentID   uint          `gorm:"not null" json:"parentID"`
    Name       string        `gorm:"type:varchar(64);not null" json:"name"`
    Sequence   uint          `gorm:"default:0;not null" json:"sequence"`
    Display    bool          `gorm:"type:tinyint;default:1;not null" json:"display"`
    URI        string        `gorm:"type:varchar(255);not null" json:"URI"`
    Type       uint          `gorm:"default:0;not null;comment:'菜单类型,0:菜单,1:按钮'" json:"type"`
    Template   string        `gorm:"type:varchar(255);not null;" json:"template"`
    Remark     string        `gorm:"type:varchar(255);default:'';not null" json:"remark"`
    Children   []*Menu       `gorm:"foreignKey:ParentID;references:ID" json:"children"`
    Permission []*Permission `gorm:"foreignKey:MenuID;references:ID" json:"permission"`
}

type MenuContext struct {
    Model
    ClosureTableModel
}

type RoleMenuRelation struct {
    Model
    RoleID uint `gorm:"not null"`
    MenuID uint `gorm:"not null"`
}

type Permission struct {
    Model
    MenuID uint `gorm:"not null"`
    Name   string
    Flag   string
}
