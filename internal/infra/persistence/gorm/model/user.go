package model

import "time"

type User struct {
	Base
	Username   string     `gorm:"type:varchar(255)"`
	Password   string     `gorm:"type:varchar(255)"`
	Email      string     `gorm:"type:varchar(255)"`
	VerifiedAt *time.Time `gorm:"type:datetime"`
	IP         string     `gorm:"type:varchar(255)"`
	Status     int8       `gorm:"type:tinyint"`
	RoleID     *uint64    `gorm:"type:bigint unsigned"`
}

func (User) TableName() string {
	return "user"
}

type UserProfile struct {
	Base
	Avatar      *string `gorm:"type:varchar(255)"`
	UserID      uint64  `gorm:"type:bigint unsigned"`
	Nickname    *string `gorm:"type:varchar(255)"`
	Gender      *int8   `gorm:"type:tinyint"`
	Description *string `gorm:"type:varchar(255)"`
	Country     *string `gorm:"type:varchar(255)"`
	Province    *string `gorm:"type:varchar(255)"`
	City        *string `gorm:"type:varchar(255)"`
}

func (UserProfile) TableName() string {
	return "user_profile"
}
