package model

type Role struct {
    Base
    Name        string `gorm:"type:varchar(64);not null" json:"name"`
    Description string `gorm:"type:varchar(255);default:'';not null" json:"description"`
}
