package model

import (
    `dpcms/internal/infra/database`
)

type Role struct {
    database.Model
    Name        string `gorm:"type:varchar(64);not null" json:"name"`
    Description string `gorm:"type:varchar(255);default:'';not null" json:"description"`
}
