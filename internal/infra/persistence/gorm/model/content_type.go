package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type ContentType struct {
	Base
	Name        string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:varchar(255)"`
}

func (ContentType) TableName() string {
	return "content_type"
}

type ContentTypeEntries struct {
	Base
	ArticleID     uint64         `gorm:"type:bigint unsigned"`
	ContentTypeID uint64         `gorm:"type:bigint unsigned"`
	Data          datatypes.JSON `gorm:"type:json"`
}

func (ContentTypeEntries) TableName() string {
	return "content_type_entries"
}

type ContentTypeSchema struct {
	Base
	ContentTypeID uint64           `gorm:"type:bigint unsigned"`
	FieldKey      string           `gorm:"type:varchar(255)"`
	FieldName     string           `gorm:"type:varchar(255)"`
	Description   string           `gorm:"type:varchar(255)"`
	MinLen        *uint64          `gorm:"type:bigint unsigned"`
	MaxLen        *uint64          `gorm:"type:bigint unsigned"`
	MinValue      *decimal.Decimal `gorm:"type:decimal(10,2)"`
	MaxValue      *decimal.Decimal `gorm:"type:decimal(10,2)"`
	MinTime       *time.Time       `gorm:"type:datetime"`
	MaxTime       *time.Time       `gorm:"type:datetime"`
	Pattern       *string          `gorm:"type:varchar(255)"`
	Sequence      int64            `gorm:"type:bigint"`
	Type          int16            `gorm:"type:smallint"`
	EnumOptions   *string          `gorm:"type:text"`
	Required      *bool            `gorm:"type:tinyint"`
	Visible       *bool            `gorm:"type:tinyint"`
	Enable        *bool            `gorm:"type:tinyint"`
}

func (ContentTypeSchema) TableName() string {
	return "content_type_schema"
}

type ContentFieldValues struct {
	Base
	ModelID     uint64           `gorm:"type:bigint unsigned"`
	ArticleID   uint64           `gorm:"type:bigint unsigned"`
	FieldKey    string           `gorm:"type:varchar(255)"`
	Type        int16            `gorm:"type:smallint"`
	StringValue *string          `gorm:"type:varchar(255)"`
	BoolValue   *bool            `gorm:"type:boolean"`
	NumberValue *decimal.Decimal `gorm:"type:decimal(10,2)"`
	TimeValue   *time.Time       `gorm:"type:datetime"`
}

func (ContentFieldValues) TableName() string {
	return "content_field_values"
}
