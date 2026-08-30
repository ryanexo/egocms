package model

import (
	"strings"
	"testing"
	
	"cms/internal/app/file/errno"
	
	"github.com/stretchr/testify/require"
)

func TestFileValidation(t *testing.T) {
    t.Run("accepts valid metadata", func(t *testing.T) {
        data := &File{}
        require.NoError(t, data.SetOriginalName("report.pdf"))
        require.NoError(t, data.SetExt(".pdf"))
        require.NoError(t, data.SetStorage("local", "2026/0826/report.pdf"))
        require.Equal(t, "report.pdf", data.OriginalName)
    })
    
    t.Run("rejects invalid metadata", func(t *testing.T) {
        data := &File{}
        require.ErrorIs(t, data.SetOriginalName(""), errno.ErrOriginalNameRequired)
        require.ErrorIs(t, data.SetOriginalName(strings.Repeat("文", OriginalNameMaxLength+1)), errno.ErrOriginalNameTooLong)
        require.ErrorIs(t, data.SetStorage("", "a"), errno.ErrInvalidDriver)
    })
}

func TestAttachmentValidation(t *testing.T) {
    attachment := &Attachment{}
    require.NoError(t, attachment.SetEntity("article", 1))
    require.NoError(t, attachment.SetType(""))
    require.Equal(t, DefaultAttachmentType, attachment.Type)
    require.ErrorIs(t, attachment.SetEntity("", 1), errno.ErrInvalidEntityType)
    require.ErrorIs(t, attachment.SetEntity("article", 0), errno.ErrEntityIDRequired)
}
