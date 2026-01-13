package model

import `dpcms/internal/infra/persistence/datatype`

type Permission struct {
    Base
    MenuID      *datatype.SafeUint64 `gorm:"index:idx_perm_menu;comment:'归类到指定菜单ID,未填写为通用权限'"`
    Name        string               `gorm:"type:varchar(255);not null"`
    Description string               `gorm:"type:varchar(255);not null"`
    Resource    string               `gorm:"index:idx_perm_resource;priority:1;type:varchar(255);not null"`
    Action      string               `gorm:"index:idx_perm_resource;priority:2;type:varchar(64);not null"`
}
