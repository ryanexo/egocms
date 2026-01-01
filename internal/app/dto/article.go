package dto

import (
    `time`
    
    `dpcms/internal/app/dto/type`
    `dpcms/internal/infra/persistence/datatype`
)

type ArticleCreateParams struct {
    Url         string               `json:"url"`
    Title       string               `validate:"required,max=255" json:"title"`
    Description string               `validate:"max=255" json:"description"`
    Content     string               `validate:"max=65535" json:"content"`
    Target      string               `validate:"http_url" json:"target"`
    Keywords    []string             `validate:"max=10" json:"keywords"`
    ModelId     *datatype.SafeUint64 `json:"modelId"`
    ModelData   map[string]any       `json:"modelData"`
}

type ArticleUpdateParams struct {
    ID          datatype.SafeUint64 `json:"id"`
    Url         string              `json:"url"`
    CategoryId  datatype.SafeUint64 `json:"categoryId"`
    Flag        int16               `json:"flag"`
    Title       string              `json:"title"`
    Description string              `json:"description"`
    Content     string              `json:"content"`
    Target      string              `json:"target"`
    Keywords    []string            `json:"keywords"`
    ModelData   map[string]any      `json:"modelData"`
}

type Article struct {
    dtotype.Base
    Url          string                `json:"url"`
    CategoryId   datatype.SafeUint64   `json:"categoryId"`
    CategoryName string                `json:"categoryName"`
    AuthorId     datatype.SafeUint64   `json:"authorId"`
    AuthorName   string                `json:"authorName"`
    Flag         int16                 `json:"flag"`
    Title        string                `json:"title"`
    Description  string                `json:"description"`
    ClickCount   datatype.SafeUint64   `json:"clickCount"`
    Status       int8                  `json:"status"`
    Target       string                `json:"target"`
    Keywords     []string              `json:"keywords"`
    ModelId      datatype.SafeUint64   `json:"modelId"`
    ModelSchema  []*ArticleModelSchema `json:"modelSchema"`
    ModelData    map[string]any        `json:"modelData"`
    PublishAt    time.Time             `json:"publishAt"`
}

type ArticleDetail struct {
    Article
    Content string `json:"content"`
}

type ArticleModelSchema struct {
    FieldKey  string
    FieldName string
    Type      int16
    Sequence  datatype.SafeInt64
}
