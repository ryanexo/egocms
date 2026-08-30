package model

import (
	"net/url"
	"path"
	"strings"

	"cms/internal/app/category/errno"
	"cms/internal/infra/store/modeltype"
)

type CategoryMeta struct {
	modeltype.Base
	CategoryID  uint64
	Title       string
	Keywords    string
	Description string
	Thumb       string
}

func (s *CategoryMeta) SetThumb(thumb string) error {
	if len(thumb) > CategoryThumbLength {
		return errno.ErrThumbInvalidLength
	}
	if thumb == "" {
		s.Thumb = ""
		return nil
	}
	if strings.HasPrefix(thumb, "/") {
		s.Thumb = path.Clean(thumb)
		return nil
	}
	parsedURL, err := url.ParseRequestURI(thumb)
	if err == nil && parsedURL.IsAbs() {
		s.Thumb = thumb
		return nil
	}
	return errno.ErrThumbInvalidPath
}

func (s *CategoryMeta) SetTitle(t string) error {
	if len(t) > CategoryTitleLength {
		return errno.ErrTitleInvalidLength
	}
	s.Title = t
	return nil
}

func (s *CategoryMeta) SetKeywords(k []string) error {
	result := make([]string, 0, len(k))
	for _, v := range k {
		keyword := strings.TrimSpace(v)
		if keyword != "" {
			result = append(result, keyword)
		}
	}
	keywords := strings.Join(result, ",")
	if len(keywords) > CategoryKeywordsLength {
		return errno.ErrKeywordsInvalidLength
	}
	s.Keywords = keywords
	return nil
}

func (s *CategoryMeta) SetDescription(d string) error {
	if len(d) > CategoryDescriptionLength {
		return errno.ErrDescriptionInvalidLength
	}
	s.Description = d
	return nil
}
