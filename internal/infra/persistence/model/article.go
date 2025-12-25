package model

import (
    `database/sql`
    `time`
)

type Article struct {
    Base
    Url         string                   `gorm:"type:varchar(255);unique" json:"url"`
    CategoryID  uint64                   `gorm:"index" json:"categoryId"`
    AuthorId    uint64                   `gorm:"index" json:"authorId"`
    Flag        int16                    `gorm:"type:smallint;index;not null;default:0" json:"flag"`
    Title       string                   `gorm:"type:varchar(255);not null;default:''" json:"title"`
    Description string                   `gorm:"type:text;not null;default:''" json:"description"`
    ClickCount  uint64                   `gorm:"not null;default:0" json:"clickCount"`
    Status      int8                     `gorm:"type:tinyint;index;not null;default:0" json:"status"`
    Target      sql.NullString           `json:"target"`
    Keywords    []ArticleKeywords        `gorm:"foreignKey:ArticleID;referenceKey:ID" json:"keywords"`
    Content     ArticleContent           `gorm:"foreignKey:ArticleID;referenceKey:ID" json:"content"`
    ModelId     uint64                   `gorm:"index" json:"-"`
    Extra       ArticleModelRelationship `gorm:"foreignKey:ArticleID;referenceKey:ID" json:"extra"`
    SubmitAt    time.Time                `gorm:"type:datetime" json:"-"`
    PublishAt   time.Time                `gorm:"type:datetime" json:"publishAt"`
    OfflineAt   time.Time                `gorm:"type:datetime" json:"-"`
    RejectAt    time.Time                `gorm:"type:datetime" json:"-"`
}

type ArticleKeywords struct {
    ID        uint64 `gorm:"primaryKey" json:"-"`
    ArticleID uint64 `gorm:"index" json:"articleId"`
    Keyword   string `gorm:"type:varchar(255);not null" json:"keyword"`
}

type ArticleContent struct {
    ID        uint64 `gorm:"primaryKey" json:"-"`
    ArticleID uint64 `gorm:"index" json:"-"`
    Content   string `gorm:"type:text;not null;default:''" json:"content"`
}
