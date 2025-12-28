package model

import (
    `database/sql`
    
    `dpcms/internal/infra/persistence/customvalue`
    
    `github.com/shopspring/decimal`
)

type ArticleModel struct {
    Base
    Name        string               `gorm:"type:varchar(255);not null" json:"name"`
    Description string               `gorm:"type:varchar(255);not null" json:"description"`
    Definition  []ArticleModelSchema `gorm:"foreignKey:ModelId;referenceKey:ID"`
}

type ArticleModelJsonData struct {
    ID        uint64              `gorm:"primaryKey"`
    ArticleId uint64              `gorm:"index:idx_art,priority:1"`
    ModelId   uint64              `gorm:"index:idx_art,priority:2"`
    Data      customvalue.JSONMap `gorm:"type:text;" json:"data"`
}

type ArticleModelSchema struct {
    ID          uint64                 `gorm:"primaryKey"`
    ModelId     uint64                 `gorm:"uniqueIndex:idx_field,priority:1;" json:"modelId"`
    FieldKey    string                 `gorm:"type:varchar(255);not null;uniqueIndex:idx_field_key,priority:2" json:"fieldKey"`
    FieldName   string                 `gorm:"type:varchar(255);not null;" json:"fieldName"`
    Description string                 `gorm:"type:varchar(255);not null;default:''" json:"description"`
    MinLen      int                    `gorm:"default:0" json:"minLen"`
    MaxLen      int                    `gorm:"default:0" json:"maxLen"`
    MinValue    decimal.NullDecimal    `gorm:"decimal(10,2)" json:"minValue"`
    MaxValue    decimal.NullDecimal    `gorm:"decimal(10,2)" json:"maxValue"`
    MinTime     sql.NullTime           `json:"minTime"`
    MaxTime     sql.NullTime           `json:"maxTime"`
    Pattern     string                 `gorm:"type:varchar(255)" json:"pattern"`
    Sequence    int64                  `gorm:"index;default:0" json:"sequence"`
    Type        int16                  `gorm:"type:smallint;not null" json:"type"`
    EnumOptions customvalue.EnumValues `gorm:"type:text" json:"enumOptions"`
    Required    bool                   `gorm:"type:bool;default:false" json:"required"`
    Hidden      bool                   `gorm:"type:bool;default:false" json:"hidden"`
    Enable      bool                   `gorm:"type:bool;default:true" json:"enable"`
}

type ArticleModelData struct {
    ID          uint64              `gorm:"primaryKey"`
    ModelId     uint64              `gorm:"index:idx_art,priority:1"`
    ArticleId   uint64              `gorm:"index:idx_art,priority:2"`
    FieldKey    string              `gorm:"type:varchar(255)"`
    Type        int16               `gorm:"type:smallint;not null"`
    ValueString sql.NullString      `gorm:"type:varchar(255)"`
    ValueBool   sql.NullBool        `gorm:"type:bool"`
    ValueNumber decimal.NullDecimal `gorm:"type:decimal(10,2)"`
    ValueTime   sql.NullTime        `gorm:"type:datetime"`
}
