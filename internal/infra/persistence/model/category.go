package model

type CategoryContext struct {
    ID uint64 `gorm:"primaryKey"`
    ClosureTableModel
}

type Category struct {
    Base
    ParentID *uint64      `gorm:"index"`
    Sequence int64        `gorm:"not null;index"`
    Name     string       `gorm:"type:varchar(255);not null"`
    Path     string       `gorm:"type:varchar(64);not null;uniqueIndex"`
    URL      *string      `gorm:"type:varchar(255);default:null"`
    SEO      *CategorySeo `gorm:"foreignKey:CategoryID;references:ID"`
}

type CategorySeo struct {
    Base
    CategoryID  uint64 `gorm:"not null;index"`
    Title       string `gorm:"type:varchar(255);not null"`
    Keywords    string `gorm:"type:varchar(255);not null"`
    Description string `gorm:"type:varchar(255);not null"`
}
