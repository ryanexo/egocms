package model

import (
    `database/sql`
    
    `dpcms/internal/infra/persistence/datatype`
    
    `github.com/shopspring/decimal`
)

type ArticleModel struct {
    Base
    Name        string                `gorm:"type:varchar(255);not null"`
    Description string                `gorm:"type:varchar(255);not null"`
    Schema      []*ArticleModelSchema `gorm:"foreignKey:ModelID;referenceKey:ID"`
}

type ArticleModelJsonData struct {
    Base
    ArticleID datatype.SafeUint64 `gorm:"index:idx_artid_modelid,priority:1"`
    ModelID   datatype.SafeUint64 `gorm:"index:idx_artid_modelid,priority:2"`
    Data      datatype.JSONMap    `gorm:"type:text;"`
}

type ArticleModelSchema struct {
    Base
    ModelID     datatype.SafeUint64 `gorm:"uniqueIndex:idx_field_key,priority:1;"`
    FieldKey    string              `gorm:"type:varchar(255);not null;uniqueIndex:idx_field_key,priority:2"`
    FieldName   string              `gorm:"type:varchar(255);not null;"`
    Description string              `gorm:"type:varchar(255);not null;default:''"`
    MinLen      datatype.SafeUint64 `gorm:"default:0"`
    MaxLen      datatype.SafeUint64 `gorm:"default:0"`
    MinValue    decimal.NullDecimal `gorm:"decimal(10,2)"`
    MaxValue    decimal.NullDecimal `gorm:"decimal(10,2)"`
    MinTime     sql.NullTime
    MaxTime     sql.NullTime
    Pattern     string              `gorm:"type:varchar(255)"`
    Sequence    datatype.SafeInt64  `gorm:"index;default:0"`
    Type        int16               `gorm:"type:smallint;not null"`
    EnumOptions datatype.EnumValues `gorm:"type:text"`
    Required    datatype.BoolInt8   `gorm:"type:tinyint;default:0"`
    Hidden      datatype.BoolInt8   `gorm:"type:tinyint;default:0"`
    Enable      datatype.BoolInt8   `gorm:"type:tinyint;default:1"`
}

type ArticleModelData struct {
    Base
    ModelID     datatype.SafeUint64 `gorm:"index:idx_artid_modelid_fieldkey,priority:1"`
    ArticleID   datatype.SafeUint64 `gorm:"index:idx_artid_modelid_fieldkey,priority:2"`
    FieldName   string              `gorm:"type:varchar(255)"`
    FieldKey    string              `gorm:"type:varchar(255);index:idx_artid_modelid_fieldkey,priority:3"`
    Type        int16               `gorm:"type:smallint;not null"`
    ValueString sql.NullString      `gorm:"type:varchar(255)"`
    ValueBool   sql.NullBool        `gorm:"type:bool"`
    ValueNumber decimal.NullDecimal `gorm:"type:decimal(10,2)"`
    ValueTime   sql.NullTime        `gorm:"type:datetime"`
}
