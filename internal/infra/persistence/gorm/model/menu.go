package model

type Menu struct {
	Base
	ParentID uint64  `gorm:"type:bigint unsigned"`
	Type     int8    `gorm:"type:tinyint"`
	Name     string  `gorm:"type:varchar(64)"`
	Affix    bool    `gorm:"type:tinyint"`
	Icon     *string `gorm:"type:varchar(64)"`
	URL      *string `gorm:"type:varchar(255)"`
	Sequence int64   `gorm:"type:bigint"`
	Visible  bool    `gorm:"type:tinyint"`
	Path     string  `gorm:"type:varchar(255)"`
	Template string  `gorm:"type:varchar(255)"`
	Remark   string  `gorm:"type:varchar(255)"`
}

func (Menu) TableName() string {
	return "menu"
}

type MenuContext struct {
	ID         uint64 `gorm:"type:bigint unsigned;primaryKey"`
	Ancestor   uint64 `gorm:"type:bigint unsigned"`
	Descendant uint64 `gorm:"type:bigint unsigned"`
	Distance   uint64 `gorm:"type:bigint unsigned"`
}

func (MenuContext) TableName() string {
	return "menu_context"
}
