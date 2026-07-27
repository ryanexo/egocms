package domain

type ContentType struct {
    id          uint64
    name        string
    description string
}

func NewContentType(name string, description string) ContentType {
    return ContentType{name: name, description: description}
}
