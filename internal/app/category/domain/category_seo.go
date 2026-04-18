package domain

type CategorySeo struct {
    title       string
    keywords    []string
    description string
}

func NewCategorySeo(title string, keywords []string, description string) *CategorySeo {
    return &CategorySeo{
        title,
        keywords,
        description,
    }
}
