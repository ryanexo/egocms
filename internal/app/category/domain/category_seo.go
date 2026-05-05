package domain

type CategorySEO struct {
    title       string
    keywords    []string
    description string
}

func NewCategorySeo(title string, keywords []string, description string) CategorySEO {
    return CategorySEO{
        title,
        keywords,
        description,
    }
}

func (c CategorySEO) Title() string {
    return c.title
}

func (c CategorySEO) Keywords() []string {
    return c.keywords
}

func (c CategorySEO) Description() string {
    return c.description
}
