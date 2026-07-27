package model

type Category struct {
	Base
	ParentID uint64 `gorm:"type:bigint unsigned"`
	Sequence int64  `gorm:"type:bigint"`
	Name     string `gorm:"type:varchar(255)"`
	Path     string `gorm:"type:varchar(64)"`
	Visible  bool   `gorm:"type:tinyint"`
}

func (Category) TableName() string {
	return "category"
}

type CategorySeo struct {
	Base
	CategoryID  uint64 `gorm:"type:bigint unsigned"`
	Title       string `gorm:"type:varchar(255)"`
	Keywords    string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:varchar(255)"`
}

func (CategorySeo) TableName() string {
	return "category_seo"
}

type CategoryContext struct {
	ID         uint64 `gorm:"type:bigint unsigned;primaryKey"`
	Ancestor   uint64 `gorm:"type:bigint unsigned"`
	Descendant uint64 `gorm:"type:bigint unsigned"`
	Distance   uint64 `gorm:"type:bigint unsigned"`
}

func (CategoryContext) TableName() string {
	return "category_context"
}
