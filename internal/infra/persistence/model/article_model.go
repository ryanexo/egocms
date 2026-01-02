package model

import (
    `database/sql`
    
    `dpcms/internal/infra/persistence/customvalue`
    `dpcms/internal/infra/persistence/datatype`
    
    `github.com/shopspring/decimal`
)

type ArticleModel struct {
    Base
    Name        string                `gorm:"type:varchar(255);not null" json:"name"`
    Description string                `gorm:"type:varchar(255);not null" json:"description"`
    Schema      []*ArticleModelSchema `gorm:"foreignKey:ModelId;referenceKey:ID" json:"schema"`
}

type ArticleModelJsonData struct {
    Base
    ArticleId datatype.SafeUint64 `gorm:"index:idx_artid_modelid,priority:1"`
    ModelId   datatype.SafeUint64 `gorm:"index:idx_artid_modelid,priority:2"`
    Data      customvalue.JSONMap `gorm:"type:text;" json:"data"`
}

type ArticleModelSchema struct {
    Base
    ModelId     datatype.SafeUint64    `gorm:"uniqueIndex:idx_field_key,priority:1;" json:"modelId"`
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
    Sequence    datatype.SafeInt64     `gorm:"index;default:0" json:"sequence"`
    Type        int16                  `gorm:"type:smallint;not null" json:"type"`
    EnumOptions customvalue.EnumValues `gorm:"type:text" json:"enumOptions"`
    Required    bool                   `gorm:"type:bool;default:false" json:"required"`
    Hidden      bool                   `gorm:"type:bool;default:false" json:"hidden"`
    Enable      bool                   `gorm:"type:bool;default:true" json:"enable"`
}

type ArticleModelData struct {
    Base
    ModelId     datatype.SafeUint64 `gorm:"index:idx_artid_modelid_fieldkey,priority:1"`
    ArticleId   datatype.SafeUint64 `gorm:"index:idx_artid_modelid_fieldkey,priority:2"`
    FieldName   string              `gorm:"type:varchar(255)"`
    FieldKey    string              `gorm:"type:varchar(255);index:idx_artid_modelid_fieldkey,priority:3"`
    Type        int16               `gorm:"type:smallint;not null"`
    ValueString sql.NullString      `gorm:"type:varchar(255)"`
    ValueBool   sql.NullBool        `gorm:"type:bool"`
    ValueNumber decimal.NullDecimal `gorm:"type:decimal(10,2)"`
    ValueTime   sql.NullTime        `gorm:"type:datetime"`
}
