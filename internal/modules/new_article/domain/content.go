package domain

import (
	"strings"
	"unicode/utf8"

	"github.com/microcosm-cc/bluemonday"
)

const (
    maxTitleLength     = 255
    maxSummaryLength   = 500
    maxChangeLogLength = 500
    maxContentBytes    = 65535
)

type Revision struct {
    title     string
    content   string
    summary   string
    changeLog string
}

func NewRevision(title, content, summary, changeLog string) (Revision, error) {
    plainTextPolicy := bluemonday.StrictPolicy()
    title = strings.TrimSpace(plainTextPolicy.Sanitize(title))
    if title == "" {
        return Revision{}, ErrTitleRequired
    }
    if utf8.RuneCountInString(title) > maxTitleLength {
        return Revision{}, ErrTitleTooLong
    }

    content = bluemonday.UGCPolicy().Sanitize(content)
    if len(content) > maxContentBytes {
        return Revision{}, ErrContentTooLong
    }

    summary = strings.TrimSpace(plainTextPolicy.Sanitize(summary))
    if summary == "" {
        summary = truncateRunes(strings.TrimSpace(plainTextPolicy.Sanitize(content)), maxSummaryLength)
    }
    if utf8.RuneCountInString(summary) > maxSummaryLength {
        return Revision{}, ErrSummaryTooLong
    }

    changeLog = strings.TrimSpace(plainTextPolicy.Sanitize(changeLog))
    if utf8.RuneCountInString(changeLog) > maxChangeLogLength {
        return Revision{}, ErrChangeLogTooLong
    }

    return Revision{
        title:     title,
        content:   content,
        summary:   summary,
        changeLog: changeLog,
    }, nil
}

func (revision Revision) Title() string {
    return revision.title
}

func (revision Revision) Content() string {
    return revision.content
}

func (revision Revision) Summary() string {
    return revision.summary
}

func (revision Revision) ChangeLog() string {
    return revision.changeLog
}

func truncateRunes(value string, limit int) string {
    runes := []rune(value)
    if len(runes) <= limit {
        return value
    }
    return string(runes[:limit])
}
