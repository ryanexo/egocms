package model

type Role struct {
	Base
	Name        string `gorm:"type:varchar(64)"`
	Description string `gorm:"type:varchar(255)"`
}

func (Role) TableName() string {
	return "role"
}
