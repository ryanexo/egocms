package model

import `database/sql`

type ArticleModel struct {
    Base
    Name        string               `gorm:"type:varchar(255);not null" json:"name"`
    Description string               `gorm:"type:varchar(255);not null" json:"description"`
    Definition  []ArticleModelSchema `gorm:"foreignKey:ContentModelId;referenceKey:ID"`
}

type ArticleModelRelationship struct {
    ID             uint64 `gorm:"primaryKey"`
    ArticleId      uint64 `gorm:"index:idx_art,priority:1"`
    ContentModelId uint64 `gorm:"index:idx_art,priority:2"`
    Data           string `gorm:"type:text;not null" json:"data"`
}

type ArticleModelSchema struct {
    ID             uint64 `gorm:"primaryKey"`
    ContentModelId uint64 `gorm:"uniqueIndex:idx_field,priority:1;" json:"contentModelId"`
    FieldKey       string `gorm:"type:varchar(255);not null;uniqueIndex:idx_field_key,priority:2" json:"fieldKey"`
    FieldName      string `gorm:"type:varchar(255);not null;" json:"fieldName"`
    Description    string `gorm:"type:varchar(255);not null;default:''" json:"description"`
    MinLen         uint64 `gorm:"default:0" json:"minLen"`
    MaxLen         uint64 `gorm:"default:0" json:"maxLen"`
    MinValue       int64  `gorm:"default:0" json:"minValue"`
    MaxValue       int64  `gorm:"default:0" json:"maxValue"`
    Pattern        string `gorm:"type:varchar(255)" json:"pattern"`
    Sequence       int64  `gorm:"index;not null;default:0" json:"sequence"`
    Type           int16  `gorm:"type:smallint;not null" json:"type"`
    Required       int8   `gorm:"type:tinyint(1);default:0"`
    EnumOptions    string `gorm:"type:text" json:"enumOptions"`
    Hidden         int8   `gorm:"type:tinyint(1);default:0"`
    Enable         int8   `gorm:"type:tinyint(1);default:0"`
    Readonly       int8   `gorm:"type:tinyint(1);default:0"`
}

type ArticleModelData struct {
    ID             uint64          `gorm:"primaryKey"`
    ContentModelId uint64          `gorm:"index:idx_art,priority:1"`
    ArticleId      uint64          `gorm:"index:idx_art,priority:2"`
    Field          string          `gorm:"type:varchar(255)"`
    ValueString    sql.NullString  `gorm:"type:varchar(255)"`
    ValueNumber    sql.NullInt64   `gorm:"type:bigint"`
    ValueBool      sql.NullBool    `gorm:"type:tinyint(1)"`
    ValueFloat     sql.NullFloat64 `gorm:"type:double"`
    ValueTime      sql.NullTime    `gorm:"type:datetime"`
}
