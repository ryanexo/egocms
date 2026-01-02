package model

import (
    `database/sql`
    `time`
    
    `dpcms/internal/infra/persistence/datatype`
)

type Article struct {
    Base
    Url         string               `gorm:"type:varchar(255);unique" json:"url"`
    CategoryID  datatype.SafeUint64  `gorm:"index" json:"categoryId"`
    AuthorId    datatype.SafeUint64  `gorm:"index" json:"authorId"`
    Flag        int16                `gorm:"type:smallint;index;not null;default:0" json:"flag"`
    Title       string               `gorm:"type:varchar(255);not null;default:''" json:"title"`
    Description string               `gorm:"type:varchar(500);not null;default:''" json:"description"`
    ClickCount  datatype.SafeUint64  `gorm:"not null;default:0" json:"clickCount"`
    Status      int8                 `gorm:"type:tinyint;index;not null;default:0" json:"status"`
    Target      sql.NullString       `json:"target"`
    Keywords    []ArticleKeywords    `gorm:"foreignKey:ArticleID;referenceKey:ID" json:"keywords"`
    Content     ArticleContent       `gorm:"foreignKey:ArticleID;referenceKey:ID" json:"content"`
    ModelId     datatype.SafeUint64  `gorm:"index" json:"-"`
    ModelData   ArticleModelJsonData `gorm:"foreignKey:ArticleID;referenceKey:ID" json:"modelData"`
    SubmitAt    time.Time            `gorm:"type:datetime" json:"-"`
    PublishAt   time.Time            `gorm:"type:datetime" json:"publishAt"`
    OfflineAt   time.Time            `gorm:"type:datetime" json:"-"`
    RejectAt    time.Time            `gorm:"type:datetime" json:"-"`
}

type ArticleKeywords struct {
    Base
    ArticleID datatype.SafeUint64 `gorm:"index" json:"articleId"`
    Keyword   string              `gorm:"type:varchar(255);not null" json:"keyword"`
}

type ArticleContent struct {
    Base
    ArticleID datatype.SafeUint64 `gorm:"index" json:"-"`
    Content   string              `gorm:"type:text;not null;default:''" json:"content"`
}
