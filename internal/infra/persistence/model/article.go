package model

import (
    `database/sql`
)

type Article struct {
    Base
    Url            string              `gorm:"type:varchar(255);unique" json:"url"`
    CategoryID     uint64              `gorm:"index" json:"categoryId"`
    UserID         uint64              `gorm:"index" json:"userId"`
    Flag           int16               `gorm:"type:smallint;index;not null;default:0" json:"flag"`
    Title          string              `gorm:"type:varchar(255);not null;default:''" json:"title"`
    Description    string              `gorm:"type:text;not null;default:''" json:"description"`
    ClickCount     uint64              `gorm:"not null;default:0" json:"clickCount"`
    Status         int8                `gorm:"type:tinyint;index;not null;default:0" json:"status"`
    Target         sql.NullString      `json:"target"`
    Keywords       []ArticleKeywords   `gorm:"foreignKey:ArticleID;referenceKey:ID" json:"keywords"`
    Content        ArticleContent      `gorm:"foreignKey:ArticleID;referenceKey:ID" json:"content"`
    ContentModelId uint64              `gorm:"index" json:"-"`
    Extra          ArticleContentModel `gorm:"foreignKey:ArticleID;referenceKey:ID" json:"extra"`
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
