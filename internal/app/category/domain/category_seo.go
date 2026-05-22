package domain

type CategorySEO struct {
    title       string
    keywords    string
    description string
}

func NewCategorySEO() CategorySEO {
    return CategorySEO{}
}

func (c CategorySEO) SetTitle(title string) {
    c.title = title
}

func (c CategorySEO) SetKeywords(keywords string) {
    c.keywords = keywords
}

func (c CategorySEO) SetDescription(description string) {
    c.description = description
}

func (c CategorySEO) Title() string {
    return c.title
}

func (c CategorySEO) Keywords() string {
    return c.keywords
}

func (c CategorySEO) Description() string {
    return c.description
}
