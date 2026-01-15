package model

import (
    `database/sql`
    `time`
    
    `cms/internal/infra/persist/datatype`
)

type Article struct {
    Base
    Url         string                `gorm:"type:varchar(255);unique"`
    CategoryID  datatype.SafeUint64   `gorm:"index"`
    Category    *Category             `gorm:"foreignKey:ID;referenceKey:CategoryID"`
    AuthorID    datatype.SafeUint64   `gorm:"index"`
    Author      *UserProfile          `gorm:"foreignKey:UserID;referenceKey:AuthorID"`
    Flag        int16                 `gorm:"type:smallint;index;not null;default:0"`
    Title       string                `gorm:"type:varchar(255);not null;default:''"`
    Description string                `gorm:"type:varchar(500);not null;default:''"`
    ClickCount  datatype.SafeUint64   `gorm:"not null;default:0"`
    Status      int8                  `gorm:"type:tinyint;index;not null;default:0"`
    Target      sql.NullString        `gorm:"type:varchar(500)"`
    Keywords    []*ArticleKeywords    `gorm:"foreignKey:ArticleID;referenceKey:ID"`
    Content     *ArticleContent       `gorm:"foreignKey:ArticleID;referenceKey:ID"`
    ModelID     *datatype.SafeUint64  `gorm:"index"`
    ModelSchema []*ArticleModelSchema `gorm:"foreignKey:ModelID;referenceKey:ModelID"`
    ModelData   *ArticleModelJsonData `gorm:"foreignKey:ArticleID;referenceKey:ID"`
    SubmitAt    time.Time             `gorm:"type:datetime"`
    PublishAt   time.Time             `gorm:"type:datetime"`
    OfflineAt   time.Time             `gorm:"type:datetime"`
    RejectAt    time.Time             `gorm:"type:datetime"`
}

type ArticleKeywords struct {
    Base
    ArticleID datatype.SafeUint64 `gorm:"index"`
    Keyword   string              `gorm:"type:varchar(255);not null"`
}

type ArticleContent struct {
    Base
    ArticleID datatype.SafeUint64 `gorm:"index"`
    Content   string              `gorm:"type:text;not null;default:''"`
}
