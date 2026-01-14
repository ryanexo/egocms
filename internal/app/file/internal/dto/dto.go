package dto

import (
    `io`
    
    `cms/internal/infra/persistence/datatype`
    `cms/internal/types`
)

type FileSaveCommand struct {
    Name string
    Data io.Reader
}

type FileInfo struct {
    types.Base
    Name string             `json:"name"`
    URL  string             `json:"url"`
    Size datatype.SafeInt64 `json:"size"`
}

type FileMeta struct {
    Name    string
    Size    datatype.SafeInt64
    IsImage datatype.BoolInt8
}
