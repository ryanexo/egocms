package dto

import (
    `dpcms/internal/infra/persistence/customvalue`
)

type ArticleModelCreateParams struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}

type ArticleModelSchemaUpdateParams struct {
    ModelId customvalue.Uint64String   `json:"modelId"`
    Data    []ArticleModelSchemaParams `json:"data"`
}

type ArticleModelSchemaParams struct {
    ID          *customvalue.Uint64String `json:"id"`
    FieldKey    string                    `validate:"required,max=255" json:"fieldKey"`
    FieldName   string                    `validate:"required,max=255" json:"fieldName"`
    Description string                    `validate:"max=255" json:"description"`
    Sequence    int64                     `json:"sequence"`
    Type        int16                     `validate:"required" json:"type"`
    Required    bool                      `json:"required"`
    MinLen      customvalue.Uint64String  `json:"minLen"`
    MaxLen      customvalue.Uint64String  `json:"maxLen"`
    MinValue    int64                     `json:"minValue"`
    MaxValue    int64                     `json:"maxValue"`
    Pattern     string                    `json:"pattern"`
    EnumOptions customvalue.EnumValues    `json:"enumOptions"`
    Hidden      bool                      `json:"hidden"`
    Enable      bool                      `json:"enable"`
    Readonly    bool                      `json:"readonly"`
}
