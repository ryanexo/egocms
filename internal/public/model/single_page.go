package model

import (
    `cms/internal/infra/store/modeltype`
)

type SinglePage struct {
    modeltype.Base
    Name        *string `gorm:"type:varchar(255)"`
    Keywords    *string `gorm:"type:varchar(255)"`
    Description *string `gorm:"type:varchar(500)"`
    Content     *string `gorm:"type:text"`
}

func (SinglePage) TableName() string {
    return "single_page"
}
