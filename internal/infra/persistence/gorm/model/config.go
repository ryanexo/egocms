package model

import "time"

type Config struct {
	ID        uint64    `gorm:"type:bigint unsigned;primaryKey"`
	CreatedAt time.Time `gorm:"type:datetime"`
	UpdatedAt time.Time `gorm:"type:datetime"`
	Field     string    `gorm:"type:varchar(64)"`
	Value     string    `gorm:"type:varchar(255)"`
}

func (Config) TableName() string {
	return "config"
}
