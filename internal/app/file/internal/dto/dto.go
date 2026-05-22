package dto

import (
	datatype2 "cms/internal/pkg/datatype"
	"io"

	"cms/internal/infra/persistence/datatype"
	"cms/internal/util/types"
)

type FileSaveCommand struct {
	Name string
	Data io.Reader
}

type FileInfo struct {
	types.Base
	Name string              `json:"name"`
	URL  string              `json:"url"`
	Size datatype2.SafeInt64 `json:"size"`
}

type FileMeta struct {
	Name    string
	Size    datatype2.SafeInt64
	IsImage datatype.BoolInt8
}
