package model

import (
	"database/sql"

	"cms/internal/infra/persistence/datatype"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type ArticleModel struct {
	Base
	Name        string                `gorm:"type:varchar(255);not null"`
	Description string                `gorm:"type:varchar(255);not null"`
	Schema      []*ArticleModelSchema `gorm:"foreignKey:ModelID;referenceKey:ID"`
}

type ArticleModelJsonData struct {
	Base
	ArticleID uint64 `gorm:"index:idx_artid_modelid,priority:1"`
	ModelID   uint64 `gorm:"index:idx_artid_modelid,priority:2"`
	Data      datatypes.JSON
}

type ArticleModelSchema struct {
	Base
	ModelID     uint64              `gorm:"uniqueIndex:idx_field_key,priority:1;"`
	FieldKey    string              `gorm:"type:varchar(255);not null;uniqueIndex:idx_field_key,priority:2"`
	FieldName   string              `gorm:"type:varchar(255);not null;"`
	Description string              `gorm:"type:varchar(255);not null;default:''"`
	MinLen      uint64              `gorm:"default:0"`
	MaxLen      uint64              `gorm:"default:0"`
	MinValue    decimal.NullDecimal `gorm:"decimal(10,2)"`
	MaxValue    decimal.NullDecimal `gorm:"decimal(10,2)"`
	MinTime     sql.NullTime
	MaxTime     sql.NullTime
	Pattern     string              `gorm:"type:varchar(255)"`
	Sequence    int64               `gorm:"index;default:0"`
	Type        int16               `gorm:"type:smallint;not null"`
	EnumOptions datatype.EnumValues `gorm:"type:text"`
	Required    bool                `gorm:"type:tinyint;default:0"`
	Hidden      bool                `gorm:"type:tinyint;default:0"`
	Enable      bool                `gorm:"type:tinyint;default:1"`
}

type ArticleModelData struct {
	Base
	ModelID     uint64              `gorm:"index:idx_artid_modelid_fieldkey,priority:1"`
	ArticleID   uint64              `gorm:"index:idx_artid_modelid_fieldkey,priority:2"`
	FieldKey    string              `gorm:"type:varchar(255);index:idx_artid_modelid_fieldkey,priority:3"`
	Type        int16               `gorm:"type:smallint;not null"`
	ValueString sql.NullString      `gorm:"type:varchar(255)"`
	ValueBool   sql.NullBool        `gorm:"type:bool"`
	ValueNumber decimal.NullDecimal `gorm:"type:decimal(10,2)"`
	ValueTime   sql.NullTime        `gorm:"type:datetime"`
}
