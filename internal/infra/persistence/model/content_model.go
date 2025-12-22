package model

import `database/sql`

type ContentModel struct {
    Base
    Name        string                   `gorm:"type:varchar(255);not null" json:"name"`
    Description string                   `gorm:"type:varchar(255);not null" json:"description"`
    Definition  []ContentModelDefinition `gorm:"foreignKey:ContentModelId;referenceKey:ID"`
}

type ArticleContentModel struct {
    ID             uint64 `gorm:"primaryKey"`
    ArticleId      uint64 `gorm:"index:idx_art,priority:1"`
    ContentModelId uint64 `gorm:"index:idx_art,priority:2"`
    Data           string `gorm:"type:text;not null" json:"data"`
}

type ContentModelDefinition struct {
    ID             uint64 `gorm:"primaryKey"`
    ContentModelId uint64 `gorm:"uniqueIndex:idx_field,priority:1;" json:"contentModelId"`
    Field          string `gorm:"type:varchar(255);not null;uniqueIndex:idx_field,priority:2" json:"field"`
    Name           string `gorm:"type:varchar(255);not null;" json:"name"`
    Description    string `gorm:"type:varchar(255);not null;default:''" json:"description"`
    Sequence       int64  `gorm:"index;not null;default:0" json:"sequence"`
    Type           int16  `gorm:"type:smallint;not null" json:"type"`
    Required       int8   `gorm:"type:tinyint(1);default:0"`
    ConfigJSON     string `gorm:"type:text" json:"config"`
}

type ContentModelData struct {
    ID             uint64          `gorm:"primaryKey"`
    ContentModelId uint64          `gorm:"index"`
    ArticleId      uint64          `gorm:"index"`
    Name           string          `gorm:"type:varchar(255)"`
    Field          string          `gorm:"type:varchar(255)"`
    ValueString    sql.NullString  `gorm:"type:varchar(255)"`
    ValueNumber    sql.NullInt64   `gorm:"type:bigint"`
    ValueBool      sql.NullBool    `gorm:"type:tinyint(1)"`
    ValueFlag      sql.NullFloat64 `gorm:"type:double"`
    ValueTime      sql.NullTime    `gorm:"type:datetime"`
}
