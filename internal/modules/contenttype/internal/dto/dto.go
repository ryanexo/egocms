package dto

import (
    "time"
    
    `cms/internal/public/apitype`
    `cms/internal/public/jsontype`
    
    "cms/internal/infra/store/datatype"
)

type ArticleModel struct {
    apitype.Base
    Name        string `json:"name"`
    Description string `json:"description"`
}

type ArticleModelCreateParams struct {
    Name        string `validate:"required" json:"name"`
    Description string `validate:"required" json:"description"`
}

type ArticleModelUpdateParams struct {
    apitype.ResourceID
    Name        string `json:"name"`
    Description string `json:"description"`
}

type ArticleModelSchemaUpdateParams struct {
    apitype.ResourceID
    Data []*ArticleModelSchemaParams `json:"data"`
}

type ArticleModelSchemaParams struct {
    ID          *jsontype.SafeUint64 `json:"id" apitype:"string"`
    FieldKey    string               `validate:"required,max=255" json:"fieldKey"`
    FieldName   string               `validate:"required,max=255" json:"fieldName"`
    Description string               `validate:"max=255" json:"description"`
    Sequence    jsontype.SafeInt64   `json:"sequence"`
    Type        int16                `validate:"required" json:"type"`
    Required    datatype.BoolInt8    `json:"required" apitype:"boolean"`
    MinLen      jsontype.SafeUint64  `json:"minLen" apitype:"string"`
    MaxLen      jsontype.SafeUint64  `json:"maxLen" apitype:"string"`
    MinValue    jsontype.SafeInt64   `json:"minValue"`
    MaxValue    jsontype.SafeInt64   `json:"maxValue"`
    MinTime     time.Time            `json:"minTime"`
    MaxTime     time.Time            `json:"maxTime"`
    Pattern     string               `json:"pattern"`
    EnumOptions datatype.EnumValues  `json:"enumOptions"`
    Hidden      datatype.BoolInt8    `json:"hidden" apitype:"boolean"`
    Enable      datatype.BoolInt8    `json:"enable" apitype:"boolean"`
}

type ArticleModelListParams struct {
    apitype.Pagination
    Name *string `json:"name"`
}
