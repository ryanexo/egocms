package dto

import (
    `time`
    
    `dpcms/internal/infra/persistence/datatype`
)

type ArticleCreateParams struct {
    Url         string               `json:"url"`
    CategoryId  datatype.SafeUint64  `validate:"required" json:"categoryId"`
    AuthorId    datatype.SafeUint64  `json:"-" swaggerignore:"true"`
    Flag        int16                `json:"flag"`
    Title       string               `validate:"required,max=255" json:"title"`
    Description string               `validate:"max=255" json:"description"`
    Content     string               `validate:"max=65535" json:"content"`
    Target      string               `validate:"omitempty,http_url" json:"target"`
    Keywords    []string             `validate:"max=10" json:"keywords"`
    ModelId     *datatype.SafeUint64 `json:"modelId" swaggertype:"string"`
    ModelData   map[string]any       `json:"modelData"`
}

type ArticleUpdateParams struct {
    ID          datatype.SafeUint64 `json:"id" swaggertype:"string"`
    Url         string              `json:"url"`
    CategoryId  datatype.SafeUint64 `json:"categoryId" swaggertype:"string"`
    Flag        int16               `json:"flag"`
    Title       string              `json:"title"`
    Description string              `json:"description"`
    Content     string              `json:"content"`
    Target      string              `json:"target"`
    Keywords    []string            `json:"keywords"`
    ModelData   map[string]any      `json:"modelData"`
}

type Article struct {
    Base
    Url          string                `json:"url"`
    CategoryID   datatype.SafeUint64   `json:"categoryId" swaggertype:"string"`
    CategoryName string                `json:"categoryName"`
    AuthorID     datatype.SafeUint64   `json:"authorId" swaggertype:"string"`
    AuthorName   string                `json:"authorName"`
    Flag         int16                 `json:"flag"`
    Title        string                `json:"title"`
    Description  string                `json:"description"`
    Content      *string               `json:"content,omitempty"`
    ClickCount   datatype.SafeUint64   `json:"clickCount" swaggertype:"string"`
    Status       int8                  `json:"status"`
    Target       string                `json:"target"`
    Keywords     []string              `json:"keywords"`
    ModelID      *datatype.SafeUint64  `json:"modelId" swaggertype:"string"`
    ModelSchema  []*ArticleModelSchema `json:"modelSchema"`
    ModelData    map[string]any        `json:"modelData"`
    PublishAt    time.Time             `json:"publishAt"`
}

type ArticleModelSchema struct {
    FieldKey  string
    FieldName string
    Type      int16
    Sequence  datatype.SafeInt64 `swaggertype:"string"`
}
