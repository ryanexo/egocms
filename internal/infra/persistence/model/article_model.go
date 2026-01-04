package model

import (
    `database/sql`
    
    `dpcms/internal/infra/persistence/customvalue`
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
    ArticleId datatype.SafeUint64 `gorm:"index:idx_artid_modelid,priority:1"`
    ModelId   datatype.SafeUint64 `gorm:"index:idx_artid_modelid,priority:2"`
    Data      customvalue.JSONMap `gorm:"type:text;"`
}

type ArticleModelSchema struct {
    Base
    ModelId     datatype.SafeUint64 `gorm:"uniqueIndex:idx_field_key,priority:1;"`
    FieldKey    string              `gorm:"type:varchar(255);not null;uniqueIndex:idx_field_key,priority:2"`
    FieldName   string              `gorm:"type:varchar(255);not null;"`
    Description string              `gorm:"type:varchar(255);not null;default:''"`
    MinLen      int                 `gorm:"default:0"`
    MaxLen      int                 `gorm:"default:0"`
    MinValue    decimal.NullDecimal `gorm:"decimal(10,2)"`
    MaxValue    decimal.NullDecimal `gorm:"decimal(10,2)"`
    MinTime     sql.NullTime
    MaxTime     sql.NullTime
    Pattern     string                 `gorm:"type:varchar(255)"`
    Sequence    datatype.SafeInt64     `gorm:"index;default:0"`
    Type        int16                  `gorm:"type:smallint;not null"`
    EnumOptions customvalue.EnumValues `gorm:"type:text"`
    Required    bool                   `gorm:"type:bool;default:false"`
    Hidden      bool                   `gorm:"type:bool;default:false"`
    Enable      bool                   `gorm:"type:bool;default:true"`
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
