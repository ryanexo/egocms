package new_article

import (
    "time"
    
    `cms/internal/public/apitype`
    `cms/internal/public/jsontype`
)

type CreateParams struct {
    ContentTypeID *jsontype.SafeUint64  `json:"contentTypeId" apitype:"string"`
    URL           string                `json:"url" validate:"omitempty,http_url"`
    Slug          string                `json:"slug" validate:"omitempty,max=255"`
    Title         string                `json:"title" validate:"required,max=255"`
    Content       string                `json:"content" validate:"max=65535"`
    Summary       string                `json:"summary" validate:"max=500"`
    ChangeLog     string                `json:"changeLog" validate:"max=500"`
    CategoryIDs   []jsontype.SafeUint64 `json:"categoryIds"`
    Tags          []string              `json:"tags" validate:"max=10,dive,max=64"`
    ContentData   map[string]any        `json:"contentData"`
}

type UpdateParams struct {
    ID          jsontype.SafeUint64   `json:"id" validate:"required" apitype:"string"`
    URL         string                `json:"url" validate:"omitempty,http_url"`
    Slug        string                `json:"slug" validate:"omitempty,max=255"`
    Title       string                `json:"title" validate:"required,max=255"`
    Content     string                `json:"content" validate:"max=65535"`
    Summary     string                `json:"summary" validate:"max=500"`
    ChangeLog   string                `json:"changeLog" validate:"max=500"`
    CategoryIDs []jsontype.SafeUint64 `json:"categoryIds"`
    Tags        []string              `json:"tags" validate:"max=10,dive,max=64"`
    ContentData map[string]any        `json:"contentData"`
}

type ContentSchema struct {
    FieldKey    string `json:"fieldKey"`
    FieldName   string `json:"fieldName"`
    Description string `json:"description"`
    Type        int16  `json:"type"`
    Sequence    int64  `json:"sequence" apitype:"string"`
    Required    bool   `json:"required"`
    Visible     bool   `json:"visible"`
    Enabled     bool   `json:"enabled"`
}

type Revision struct {
    ID        jsontype.SafeUint64 `json:"id" apitype:"string"`
    VersionNo jsontype.SafeUint64 `json:"versionNo" apitype:"string"`
    Title     string              `json:"title"`
    Content   string              `json:"content"`
    Summary   string              `json:"summary"`
    ChangeLog string              `json:"changeLog"`
    CreatedAt time.Time           `json:"createdAt"`
}

type Article struct {
    apitype.Base
    ContentTypeID      *jsontype.SafeUint64  `json:"contentTypeId,omitempty" apitype:"string"`
    AuthorID           jsontype.SafeUint64   `json:"authorId" apitype:"string"`
    URL                *string               `json:"url,omitempty"`
    Slug               *string               `json:"slug,omitempty"`
    Status             int8                  `json:"status"`
    CurrentVersionID   jsontype.SafeUint64   `json:"currentVersionId" apitype:"string"`
    PublishedVersionID *jsontype.SafeUint64  `json:"publishedVersionId,omitempty" apitype:"string"`
    CurrentRevision    Revision              `json:"currentRevision"`
    CategoryIDs        []jsontype.SafeUint64 `json:"categoryIds" apitype:"array,string"`
    Tags               []string              `json:"tags"`
    PublishedAt        *time.Time            `json:"publishedAt,omitempty"`
    ContentSchemas     []ContentSchema       `json:"contentSchemas"`
    ContentData        map[string]any        `json:"contentData"`
}
