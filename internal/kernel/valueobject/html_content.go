package valueobject

import `github.com/microcosm-cc/bluemonday`

type HTMLContent struct {
    content string
}

func NewHTMLContent(content string) HTMLContent {
    return HTMLContent{bluemonday.UGCPolicy().Sanitize(content)}
}

func (s HTMLContent) String() string {
    return s.content
}
