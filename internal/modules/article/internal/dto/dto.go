package dto

import (
    `cms/internal/public/apitype`
    `cms/internal/public/jsontype`
    
    "time"
)

type ArticleCreateParams struct {
    Url         string               `json:"url"`
    CategoryID  jsontype.SafeUint64  `validate:"required" json:"categoryId" apitype:"string"`
    AuthorID    jsontype.SafeUint64  `json:"-" swaggerignore:"true"`
    Flag        int16                `json:"flag"`
    Title       string               `validate:"required,max=255" json:"title"`
    Description string               `validate:"max=255" json:"description"`
    Content     string               `validate:"max=65535" json:"content"`
    Target      string               `validate:"omitempty,http_url" json:"target"`
    Keywords    []string             `validate:"max=10" json:"keywords"`
    ModelID     *jsontype.SafeUint64 `json:"modelID" apitype:"string"`
    ModelData   map[string]any       `json:"modelData"`
}

type ArticleUpdateParams struct {
    ID          jsontype.SafeUint64 `json:"id" apitype:"string"`
    Url         string              `json:"url"`
    CategoryID  jsontype.SafeUint64 `json:"categoryId" apitype:"string"`
    Flag        int16               `json:"flag"`
    Title       string              `json:"title"`
    Description string              `json:"description"`
    Content     string              `json:"content"`
    Target      string              `json:"target"`
    Keywords    []string            `json:"keywords"`
    ModelData   map[string]any      `json:"modelData"`
}

type Article struct {
    apitype.Base
    Url          string                `json:"url"`
    CategoryID   jsontype.SafeUint64   `json:"categoryId" apitype:"string"`
    CategoryName string                `json:"categoryName"`
    AuthorID     jsontype.SafeUint64   `json:"authorId" apitype:"string"`
    AuthorName   string                `json:"authorName"`
    Flag         int16                 `json:"flag"`
    Title        string                `json:"title"`
    Description  string                `json:"description"`
    Content      *string               `json:"content,omitempty"`
    ClickCount   jsontype.SafeUint64   `json:"clickCount" apitype:"string"`
    Status       int8                  `json:"status"`
    Target       string                `json:"target"`
    Keywords     []string              `json:"keywords"`
    ModelID      *jsontype.SafeUint64  `json:"modelId" apitype:"string"`
    ModelSchema  []*ArticleModelSchema `json:"modelSchema"`
    ModelData    map[string]any        `json:"modelData"`
    PublishAt    time.Time             `json:"publishAt"`
}

type ArticleModelSchema struct {
    FieldKey  string
    FieldName string
    Type      int16
    Sequence  jsontype.SafeInt64 `apitype:"string"`
}
