package article

import "github.com/microcosm-cc/bluemonday"

type Content struct {
	description string
	content     string
}

func NewContent(description, content string) Content {
	filter := bluemonday.UGCPolicy()
	return Content{description: filter.Sanitize(description), content: filter.Sanitize(content)}
}

func (s Content) Description() string {
	if s.description != "" {
		return s.description
	}
	return s.content[:200]
}

func (s Content) Content() string {
	return s.content
}
