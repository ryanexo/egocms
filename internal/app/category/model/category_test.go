package model

import (
	"strings"
	"testing"

	"cms/internal/app/category/errno"

	"github.com/stretchr/testify/require"
)

func TestCategorySetNameValidatesNewValue(t *testing.T) {
	category := &Category{}

	err := category.SetName(strings.Repeat("a", CategoryNameLength+1))

	require.ErrorIs(t, err, errno.ErrInvalidNameLength)
	require.Empty(t, category.Name)
}

func TestCategorySetPathValidatesNewValue(t *testing.T) {
	category := &Category{}

	err := category.SetPath("invalid/path")

	require.ErrorIs(t, err, errno.ErrInvalidPath)
	require.Empty(t, category.Path)
}

func TestCategoryMetaSetThumb(t *testing.T) {
	tests := []struct {
		name      string
		thumb     string
		expected  string
		wantError error
	}{
		{name: "empty", thumb: "", expected: ""},
		{name: "local path", thumb: "/images/../thumb.png", expected: "/thumb.png"},
		{name: "absolute URL", thumb: "https://example.com/thumb.png", expected: "https://example.com/thumb.png"},
		{name: "invalid", thumb: "images/thumb.png", wantError: errno.ErrThumbInvalidPath},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			meta := &CategoryMeta{}
			err := meta.SetThumb(test.thumb)
			if test.wantError != nil {
				require.ErrorIs(t, err, test.wantError)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.expected, meta.Thumb)
		})
	}
}

func TestCategoryMetaSetKeywords(t *testing.T) {
	meta := &CategoryMeta{}

	err := meta.SetKeywords([]string{" go ", "", "cms"})

	require.NoError(t, err)
	require.Equal(t, "go,cms", meta.Keywords)
}
