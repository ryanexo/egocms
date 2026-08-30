package model

import (
    `time`
    
    `cms/internal/infra/store/modeltype`
)

type Article struct {
    modeltype.Base
    CurrentVersionID   *uint64 `gorm:"type:bigint unsigned"`
    PublishedVersionID *uint64 `gorm:"type:bigint unsigned"`
    ContentTypeID      *uint64 `gorm:"type:bigint unsigned"`
    AuthorID           uint64  `gorm:"type:bigint unsigned"`
    URL                *string `gorm:"type:varchar(255)"`
    Slug               *string `gorm:"type:varchar(255)"`
    Title              *string `gorm:"type:varchar(255)"`
    Summary            *string `gorm:"type:varchar(500)"`
    Status             int8    `gorm:"type:tinyint"`
}

func (Article) TableName() string {
    return "article"
}

type ArticleVersion struct {
    modeltype.Base
    ArticleID uint64  `gorm:"type:bigint unsigned"`
    VersionNo uint32  `gorm:"type:int unsigned"`
    Title     string  `gorm:"type:varchar(255)"`
    Content   string  `gorm:"type:text"`
    Summary   string  `gorm:"type:varchar(500)"`
    ChangeLog *string `gorm:"type:varchar(500)"`
}

func (ArticleVersion) TableName() string {
    return "article_version"
}

type ArticlePublish struct {
    modeltype.Base
    PublishAt   time.Time `gorm:"type:datetime"`
    ArticleID   uint64    `gorm:"type:bigint unsigned"`
    VersionID   uint64    `gorm:"type:bigint unsigned"`
    PublishType int8      `gorm:"type:tinyint"`
    Remark      *string   `gorm:"type:varchar(255)"`
}

func (ArticlePublish) TableName() string {
    return "article_publish"
}

type ArticleCategoryRel struct {
    ArticleID  uint64 `gorm:"type:bigint unsigned;primaryKey"`
    CategoryID uint64 `gorm:"type:bigint unsigned;primaryKey"`
}

func (ArticleCategoryRel) TableName() string {
    return "article_category_rel"
}

type ArticleTag struct {
    modeltype.Base
    Name string `gorm:"type:varchar(64)"`
}

func (ArticleTag) TableName() string {
    return "article_tag"
}

type ArticleTagRel struct {
    ArticleID uint64 `gorm:"type:bigint unsigned;primaryKey"`
    TagID     uint64 `gorm:"type:bigint unsigned;primaryKey"`
}

func (ArticleTagRel) TableName() string {
    return "article_tag_rel"
}

type ArticleComment struct {
    modeltype.Base
    ParentID  *uint64 `gorm:"type:bigint unsigned"`
    ArticleID uint64  `gorm:"type:bigint unsigned"`
    AuthorID  uint64  `gorm:"type:bigint unsigned"`
    Content   string  `gorm:"type:text"`
}

func (ArticleComment) TableName() string {
    return "article_comment"
}
