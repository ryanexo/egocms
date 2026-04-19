package domain

type CategorySeo struct {
    title       string
    keywords    []string
    description string
}

func NewCategorySeo(title string, keywords []string, description string) CategorySeo {
    return CategorySeo{
        title,
        keywords,
        description,
    }
}

func (c CategorySeo) Title() string {
    return c.title
}

func (c CategorySeo) Keywords() []string {
    return c.keywords
}

func (c CategorySeo) Description() string {
    return c.description
}
