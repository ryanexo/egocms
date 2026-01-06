package dto

import (
	"time"

	"dpcms/internal/infra/persistence/datatype"
	"dpcms/internal/infra/persistence/dbscope"
)

type ArticleModel struct {
	Base
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ArticleModelCreateParams struct {
	Name        string `validate:"required" json:"name"`
	Description string `validate:"required" json:"description"`
}

type ArticleModelUpdateParams struct {
	ResourceID
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ArticleModelSchemaUpdateParams struct {
	ResourceID
	Data []*ArticleModelSchemaParams `json:"data"`
}

type ArticleModelSchemaParams struct {
	ID          *datatype.SafeUint64 `json:"id" swaggertype:"string"`
	FieldKey    string               `validate:"required,max=255" json:"fieldKey"`
	FieldName   string               `validate:"required,max=255" json:"fieldName"`
	Description string               `validate:"max=255" json:"description"`
	Sequence    datatype.SafeInt64   `json:"sequence"`
	Type        int16                `validate:"required" json:"type"`
	Required    datatype.BoolInt8    `json:"required" swaggertype:"boolean"`
	MinLen      datatype.SafeUint64  `json:"minLen" swaggertype:"string"`
	MaxLen      datatype.SafeUint64  `json:"maxLen" swaggertype:"string"`
	MinValue    datatype.SafeInt64   `json:"minValue"`
	MaxValue    datatype.SafeInt64   `json:"maxValue"`
	MinTime     time.Time            `json:"minTime"`
	MaxTime     time.Time            `json:"maxTime"`
	Pattern     string               `json:"pattern"`
	EnumOptions datatype.EnumValues  `json:"enumOptions"`
	Hidden      datatype.BoolInt8    `json:"hidden" swaggertype:"boolean"`
	Enable      datatype.BoolInt8    `json:"enable" swaggertype:"boolean"`
}

type ArticleModelListParams struct {
	dbscope.Pagination
	Name *string `json:"name"`
}
