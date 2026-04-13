package model

import "cms/internal/infra/persist/datatype"

type Menu struct {
    Base
    ParentID    datatype.SafeUint64 `gorm:"not null"`
    Type        int8                `gorm:"type:tinyint;not null;comment:'菜单类型,0菜单,1外链'"`
    Name        string              `gorm:"type:varchar(64);not null;comment:'名称'"`
    Affix       datatype.BoolInt8   `gorm:"type:tinyint;default:0;not null;comment:'是否固定'"`
    Icon        string              `gorm:"type:varchar(64);comment:'图标'"`
    ExternalURL string              `gorm:"type:varchar(255);comment:'外链URL'"`
    Sequence    datatype.SafeInt64  `gorm:"default:0;not null"`
    Visible     datatype.BoolInt8   `gorm:"type:tinyint;default:1;not null"`
    URI         string              `gorm:"type:varchar(255);not null;comment:'菜单URI'"`
    Template    string              `gorm:"type:varchar(255);not null;comment:'菜单前台模板'"`
    Remark      string              `gorm:"type:varchar(255);default:'';not null"`
    Action      []*Permission       `gorm:"foreignKey:MenuID;referenceKey:ID"`
}

type MenuContext struct {
    ID datatype.SafeUint64 `gorm:"primaryKey"`
    ClosureTableModel
}
