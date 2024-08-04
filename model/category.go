package model

type CategoryContext struct {
    Model
    ClosureTableModel
}

type Category struct {
    Model
    ParentID uint         `gorm:"not null;index"`
    Sequence uint         `gorm:"not null;index"`
    Name     string       `gorm:"type:varchar(255);not null"`
    Alias    string       `gorm:"type:varchar(64);not null;uniqueIndex"`
    Type     uint         `gorm:"type:tinyint;not null;comment:'0:普通分类,1:单页型分类,2:链接'"`
    Display  uint         `gorm:"type:tinyint;not null"`
    SEO      *CategorySeo `gorm:"foreignKey:CategoryID;references:ID"`
    Children []*Category  `gorm:"foreignKey:ParentID;references:ID"`
}

type CategorySeo struct {
    Model
    CategoryID     uint   `gorm:"not null;index"`
    SeoTitle       string `gorm:"type:varchar(255);not null"`
    SeoKeywords    string `gorm:"type:varchar(255);not null"`
    SeoDescription string `gorm:"type:varchar(255);not null"`
}
