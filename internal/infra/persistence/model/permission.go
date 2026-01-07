package model

import `dpcms/internal/infra/persistence/datatype`

type Permission struct {
    Base
    MenuID      *datatype.SafeUint64 `gorm:"comment:'归类到指定菜单ID,未填写为通用权限'"`
    Name        string               `gorm:"type:varchar(255);not null"`
    Description string               `gorm:"type:varchar(255);not null"`
    Resource    string               `gorm:"type:varchar(255);not null"`
    Action      string               `gorm:"type:varchar(64);not null"`
}
