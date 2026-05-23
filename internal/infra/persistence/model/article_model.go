package model

import (
    `time`
    
    "github.com/shopspring/decimal"
    "gorm.io/datatypes"
)

type ContentType struct {
    Base
    Name        string               `gorm:"type:varchar(255);not null"`
    Description string               `gorm:"type:varchar(255);not null"`
    Schema      []*ContentTypeSchema `gorm:"foreignKey:ContentTypeID;referenceKey:ID"`
}

type ContentEntries struct {
    Base
    ArticleID     uint64 `gorm:"index:idx_artid_ctid,priority:1"`
    ContentTypeID uint64 `gorm:"index:idx_artid_ctid,priority:2"`
    Data          datatypes.JSON
}

type ContentTypeSchema struct {
    Base
    ContentTypeID uint64           `gorm:"uniqueIndex:idx_field_key,priority:1;"`
    FieldKey      string           `gorm:"type:varchar(255);not null;uniqueIndex:idx_field_key,priority:2"`
    FieldName     string           `gorm:"type:varchar(255);not null;"`
    Description   string           `gorm:"type:varchar(255);not null;default:''"`
    MinLen        *uint64          `gorm:"default:0"`
    MaxLen        *uint64          `gorm:"default:0"`
    MinValue      *decimal.Decimal `gorm:"decimal(10,2)"`
    MaxValue      *decimal.Decimal `gorm:"decimal(10,2)"`
    MinTime       *time.Time
    MaxTime       *time.Time
    Pattern       *string `gorm:"type:varchar(255)"`
    Sequence      int64   `gorm:"index;default:0"`
    Type          int16   `gorm:"type:smallint;not null"`
    EnumOptions   *string `gorm:"type:text"`
    Required      bool    `gorm:"type:tinyint;default:0"`
    Visible       bool    `gorm:"type:tinyint;default:1"`
    Enable        bool    `gorm:"type:tinyint;default:1"`
}

type ContentFieldValues struct {
    Base
    ModelID     uint64           `gorm:"index:idx_artid_ctid_fieldkey,priority:1"`
    ArticleID   uint64           `gorm:"index:idx_artid_ctid_fieldkey,priority:2"`
    FieldKey    string           `gorm:"type:varchar(255);index:idx_artid_ctid_fieldkey,priority:3"`
    Type        int16            `gorm:"type:smallint;not null"`
    StringValue *string          `gorm:"type:varchar(255)"`
    BoolValue   *bool            `gorm:"type:bool"`
    NumberValue *decimal.Decimal `gorm:"type:decimal(10,2)"`
    TimeValue   *time.Time       `gorm:"type:datetime"`
}
