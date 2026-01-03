package dto

import (
	"dpcms/internal/app/dto/type"
	"dpcms/internal/infra/persistence/customvalue"
	"dpcms/internal/infra/persistence/datatype"
)

type ArticleModel struct {
	dtotype.Base
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ArticleModelCreateParams struct {
	Name        string              `validate:"required" json:"name"`
	Description string              `validate:"required" json:"description"`
}

type ArticleModelUpdateParams struct {
	ID          datatype.SafeUint64 `json:"id" swaggertype:"string"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
}

type ArticleModelSchemaUpdateParams struct {
	ModelId datatype.SafeUint64         `json:"modelId" swaggertype:"string"`
	Data    []*ArticleModelSchemaParams `json:"data"`
}

type ArticleModelSchemaParams struct {
	ID          *datatype.SafeUint64   `json:"id" swaggertype:"string"`
	FieldKey    string                 `validate:"required,max=255" json:"fieldKey"`
	FieldName   string                 `validate:"required,max=255" json:"fieldName"`
	Description string                 `validate:"max=255" json:"description"`
	Sequence    int32                  `json:"sequence"`
	Type        int16                  `validate:"required" json:"type"`
	Required    bool                   `json:"required"`
	MinLen      datatype.SafeUint64    `json:"minLen" swaggertype:"string"`
	MaxLen      datatype.SafeUint64    `json:"maxLen" swaggertype:"string"`
	MinValue    int32                  `json:"minValue"`
	MaxValue    int32                  `json:"maxValue"`
	Pattern     string                 `json:"pattern"`
	EnumOptions customvalue.EnumValues `json:"enumOptions"`
	Hidden      bool                   `json:"hidden"`
	Enable      bool                   `json:"enable"`
}
