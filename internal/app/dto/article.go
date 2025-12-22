package dto

type ArticleCreateParams struct {
    Url              string   `json:"url"`
    Title            string   `validate:"required,max=255" json:"title"`
    Description      string   `validate:"max=255" json:"description"`
    Content          string   `validate:"max=65535" json:"content"`
    Target           string   `validate:"http_url" json:"target"`
    Keywords         []string `validate:"max=10" json:"keywords"`
    ContentModelId   uint64   `json:"contentModelId"`
    ContentModelData string   `json:"contentModelData"`
}
