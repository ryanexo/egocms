package new_article

import (
	"time"

	"cms/internal/pkg/datatype"
	"cms/internal/util/types"
)

type CreateParams struct {
	ContentTypeID *datatype.SafeUint64 `json:"contentTypeId" swaggertype:"string"`
	URL           string               `json:"url" validate:"omitempty,http_url"`
	Slug          string               `json:"slug" validate:"omitempty,max=255"`
	Title         string               `json:"title" validate:"required,max=255"`
	Content       string               `json:"content" validate:"max=65535"`
	Summary       string               `json:"summary" validate:"max=500"`
	ChangeLog     string               `json:"changeLog" validate:"max=500"`
	CategoryIDs   []datatype.SafeUint64 `json:"categoryIds"`
	Tags          []string             `json:"tags" validate:"max=10,dive,max=64"`
	ContentData   map[string]any       `json:"contentData"`
}

type UpdateParams struct {
	ID          datatype.SafeUint64   `json:"id" validate:"required" swaggertype:"string"`
	URL         string                `json:"url" validate:"omitempty,http_url"`
	Slug        string                `json:"slug" validate:"omitempty,max=255"`
	Title       string                `json:"title" validate:"required,max=255"`
	Content     string                `json:"content" validate:"max=65535"`
	Summary     string                `json:"summary" validate:"max=500"`
	ChangeLog   string                `json:"changeLog" validate:"max=500"`
	CategoryIDs []datatype.SafeUint64 `json:"categoryIds"`
	Tags        []string              `json:"tags" validate:"max=10,dive,max=64"`
	ContentData map[string]any        `json:"contentData"`
}

type ContentSchema struct {
	FieldKey    string `json:"fieldKey"`
	FieldName   string `json:"fieldName"`
	Description string `json:"description"`
	Type        int16  `json:"type"`
	Sequence    int64  `json:"sequence" swaggertype:"string"`
	Required    bool   `json:"required"`
	Visible     bool   `json:"visible"`
	Enabled     bool   `json:"enabled"`
}

type Revision struct {
	ID        datatype.SafeUint64 `json:"id" swaggertype:"string"`
	VersionNo datatype.SafeUint64 `json:"versionNo" swaggertype:"string"`
	Title     string              `json:"title"`
	Content   string              `json:"content"`
	Summary   string              `json:"summary"`
	ChangeLog string              `json:"changeLog"`
	CreatedAt time.Time           `json:"createdAt"`
}

type Article struct {
	types.Base
	ContentTypeID      *datatype.SafeUint64 `json:"contentTypeId,omitempty" swaggertype:"string"`
	AuthorID           datatype.SafeUint64  `json:"authorId" swaggertype:"string"`
	URL                *string              `json:"url,omitempty"`
	Slug               *string              `json:"slug,omitempty"`
	Status             int8                 `json:"status"`
	CurrentVersionID   datatype.SafeUint64  `json:"currentVersionId" swaggertype:"string"`
	PublishedVersionID *datatype.SafeUint64 `json:"publishedVersionId,omitempty" swaggertype:"string"`
	CurrentRevision    Revision             `json:"currentRevision"`
	CategoryIDs        []datatype.SafeUint64 `json:"categoryIds" swaggertype:"array,string"`
	Tags               []string             `json:"tags"`
	PublishedAt        *time.Time           `json:"publishedAt,omitempty"`
	ContentSchemas     []ContentSchema      `json:"contentSchemas"`
	ContentData        map[string]any       `json:"contentData"`
}
