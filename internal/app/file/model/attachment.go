package model

import (
    "unicode/utf8"
    
    "cms/internal/app/file/errno"
    "cms/internal/infra/store/modeltype"
)

const (
    EntityTypeMaxLength     = 32
    AttachmentTypeMaxLength = 32
    DefaultAttachmentType   = "attachment"
)

type Attachment struct {
    modeltype.Base
    FileID     uint64
    EntityType string
    EntityID   uint64
    Type       string
    Sort       uint32
    File       File
}

func (attachment *Attachment) SetEntity(entityType string, entityID uint64) error {
    if entityType == "" || utf8.RuneCountInString(entityType) > EntityTypeMaxLength {
        return errno.ErrInvalidEntityType
    }
    if entityID == 0 {
        return errno.ErrEntityIDRequired
    }
    attachment.EntityType = entityType
    attachment.EntityID = entityID
    return nil
}

func (attachment *Attachment) SetType(attachmentType string) error {
    if attachmentType == "" {
        attachmentType = DefaultAttachmentType
    }
    if len(attachmentType) > AttachmentTypeMaxLength {
        return errno.ErrAttachmentTypeTooLong
    }
    attachment.Type = attachmentType
    return nil
}
