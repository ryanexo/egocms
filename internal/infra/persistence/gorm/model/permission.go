package model

type Permission struct {
	Base
	MenuID      *uint64 `gorm:"type:bigint unsigned"`
	Name        string  `gorm:"type:varchar(255)"`
	Description string  `gorm:"type:varchar(255)"`
	Resource    string  `gorm:"type:varchar(255)"`
	Action      string  `gorm:"type:varchar(64)"`
}

func (Permission) TableName() string {
	return "permission"
}
