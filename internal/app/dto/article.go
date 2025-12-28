package dto

import (
    `time`
    
    `dpcms/internal/infra/persistence/customvalue`
    `dpcms/internal/infra/persistence/model`
)

type ArticleCreateParams struct {
    Url         string                    `json:"url"`
    Title       string                    `validate:"required,max=255" json:"title"`
    Description string                    `validate:"max=255" json:"description"`
    Content     string                    `validate:"max=65535" json:"content"`
    Target      string                    `validate:"http_url" json:"target"`
    Keywords    []string                  `validate:"max=10" json:"keywords"`
    ModelId     *customvalue.Uint64String `json:"modelId"`
    ModelData   map[string]any            `json:"modelData"`
}

type Article struct {
    model.Base
    Url          string                   `json:"url"`
    CategoryId   customvalue.Uint64String `json:"categoryId"`
    CategoryName string                   `json:"categoryName"`
    AuthorId     customvalue.Uint64String `json:"authorId"`
    AuthorName   string                   `json:"authorName"`
    Flag         int16                    `json:"flag"`
    Title        string                   `json:"title"`
    Description  string                   `json:"description"`
    ClickCount   customvalue.Uint64String `json:"clickCount"`
    Status       int8                     `json:"status"`
    Target       string                   `json:"target"`
    Keywords     []string                 `json:"keywords"`
    ModelId      customvalue.Uint64String `json:"modelId"`
    ModelSchema  []ArticleModelSchema     `json:"modelSchema"`
    ModelData    map[string]any           `json:"modelData"`
    PublishAt    time.Time                `json:"publishAt"`
}

type ArticleDetail struct {
    Article
    Content string `json:"content"`
}

type ArticleModelSchema struct {
    FieldKey  string
    FieldName string
    Type      int16
    Sequence  int64
}
