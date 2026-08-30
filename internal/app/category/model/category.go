package model

import (
	"cms/internal/app/category/errno"
	"cms/internal/infra/store/modeltype"
	"cms/internal/public/utils"
)

const (
	CategoryNameLength        = 255
	CategoryTitleLength       = 255
	CategoryKeywordsLength    = 255
	CategoryDescriptionLength = 255
	CategoryThumbLength       = 500
)

type Category struct {
	modeltype.Base
	ParentID uint64
	Sequence int64
	Name     string
	Path     string
	Meta     CategoryMeta `gorm:"foreignKey:CategoryID"`
}

func (s *Category) AsRoot() {
	s.ParentID = 0
}

func (s *Category) Move(id uint64) {
	s.ParentID = id
}

func (s *Category) SetName(name string) error {
	if name == "" {
		return errno.ErrNameRequired
	}
	if len(name) > CategoryNameLength {
		return errno.ErrInvalidNameLength
	}
	s.Name = name
	return nil
}

func (s *Category) SetPath(p string) error {
	if p == "" {
		return errno.ErrPathRequired
	}
	if !utils.IsValidPath(p) {
		return errno.ErrInvalidPath
	}
	s.Path = p
	return nil
}
