package dto

import (
    `dpcms/internal/infra/persistence/customvalue`
    `dpcms/internal/infra/persistence/model`
)

type ArticleCreateParams struct {
    Url              string                    `json:"url"`
    Title            string                    `validate:"required,max=255" json:"title"`
    Description      string                    `validate:"max=255" json:"description"`
    Content          string                    `validate:"max=65535" json:"content"`
    Target           string                    `validate:"http_url" json:"target"`
    Keywords         []string                  `validate:"max=10" json:"keywords"`
    ContentModelId   *customvalue.Uint64String `json:"contentModelId"`
    ContentModelData map[string]any            `json:"contentModelData"`
}

type Article struct {
    model.Base
    Url              string                   `json:"url"`
    CategoryId       customvalue.Uint64String `json:"categoryId"`
    CategoryName     string                   `json:"categoryName"`
    AuthorId         string                   `json:"authorId"`
    AuthorName       string                   `json:"authorName"`
    Flag             int16                    `json:"flag"`
    Title            string                   `json:"title"`
    Description      string                   `json:"description"`
    ClickCount       customvalue.Uint64String `json:"clickCount"`
    Status           int8                     `json:"status"`
    Target           string                   `json:"target"`
    Keywords         []string                 `json:"keywords"`
    ContentModelId   customvalue.Uint64String `json:"contentModelId"`
    ContentModelData map[string]any           `json:"contentModelData"`
}

type ArticleDetail struct {
    Article
    Content string `json:"content"`
}
