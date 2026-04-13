package model

import (
    `database/sql`
    `time`
    
    `cms/internal/infra/persist/datatype`
)

type Article struct {
    Base
    Url         string                `gorm:"type:varchar(255);unique;comment:'自定义内容URL'"`
    CategoryID  datatype.SafeUint64   `gorm:"index;comment:'文章份额里'"`
    Category    *Category             `gorm:"foreignKey:ID;referenceKey:CategoryID"`
    AuthorID    datatype.SafeUint64   `gorm:"index;comment:'文章作者'"`
    Author      *UserProfile          `gorm:"foreignKey:UserID;referenceKey:AuthorID"`
    Flag        int16                 `gorm:"type:smallint;index;not null;default:0;comment:'文章flag，使用时转二进制'"`
    Title       string                `gorm:"type:varchar(255);not null;default:'';comment:'文章标题'"`
    Description string                `gorm:"type:varchar(500);not null;default:'';comment:'文章简介'"`
    ClickCount  datatype.SafeUint64   `gorm:"not null;default:0;comment:'点击数'"`
    Status      int8                  `gorm:"type:tinyint;index;not null;default:0;comment:'文章状态'"`
    Target      sql.NullString        `gorm:"type:varchar(500);comment:'外链类文章URL'"`
    Keywords    []*ArticleKeywords    `gorm:"foreignKey:ArticleID;referenceKey:ID"`
    Content     *ArticleContent       `gorm:"foreignKey:ArticleID;referenceKey:ID"`
    ModelID     *datatype.SafeUint64  `gorm:"index;comment:'内容模型'"`
    ModelSchema []*ArticleModelSchema `gorm:"foreignKey:ModelID;referenceKey:ModelID"`
    ModelData   *ArticleModelJsonData `gorm:"foreignKey:ArticleID;referenceKey:ID"`
    SubmitAt    time.Time             `gorm:"type:datetime;comment:'提交时间';"`
    PublishAt   time.Time             `gorm:"type:datetime;comment:'发布时间'"`
    OfflineAt   time.Time             `gorm:"type:datetime;comment:'下线时间'"`
    RejectAt    time.Time             `gorm:"type:datetime:comment:'拒审时间'"`
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
