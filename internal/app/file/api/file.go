package api

import (
	"time"

	"cms/internal/public/apitype"
	"cms/internal/public/jsontype"
)

type File struct {
	ID           jsontype.SafeUint64 `json:"id" apitype:"string"`
	CreatedAt    time.Time           `json:"createdAt"`
	UpdatedAt    time.Time           `json:"updatedAt"`
	OriginalName string              `json:"originalName"`
	Ext          string              `json:"ext"`
	Size         jsontype.SafeUint64 `json:"size" apitype:"string"`
	Driver       string              `json:"driver"`
	URL          string              `json:"url"`
}

type FileListParams struct {
	apitype.Pagination
	OriginalName *string `json:"originalName"`
	Ext          *string `json:"ext"`
	Driver       *string `json:"driver"`
}

type Attachment struct {
	ID         jsontype.SafeUint64 `json:"id" apitype:"string"`
	CreatedAt  time.Time           `json:"createdAt"`
	UpdatedAt  time.Time           `json:"updatedAt"`
	FileID     jsontype.SafeUint64 `json:"fileId" apitype:"string"`
	EntityType string              `json:"entityType"`
	EntityID   jsontype.SafeUint64 `json:"entityId" apitype:"string"`
	Type       string              `json:"type"`
	Sort       uint32              `json:"sort"`
	File       *File               `json:"file,omitempty"`
}

type AttachmentCreateParams struct {
	FileID     jsontype.SafeUint64 `json:"fileId" validate:"required" apitype:"string"`
	EntityType string              `json:"entityType" validate:"required,max=32"`
	EntityID   jsontype.SafeUint64 `json:"entityId" validate:"required" apitype:"string"`
	Type       string              `json:"type" validate:"max=32"`
	Sort       uint32              `json:"sort"`
}

type AttachmentUpdateParams struct {
	apitype.ResourceID
	Type string `json:"type" validate:"max=32"`
	Sort uint32 `json:"sort"`
}

type AttachmentListParams struct {
	apitype.Pagination
	FileID     *jsontype.SafeUint64 `json:"fileId" apitype:"string"`
	EntityType *string              `json:"entityType" validate:"omitempty,max=32"`
	EntityID   *jsontype.SafeUint64 `json:"entityId" apitype:"string"`
	Type       *string              `json:"type" validate:"omitempty,max=32"`
}
