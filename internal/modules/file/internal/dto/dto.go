package dto

import (
    "io"
    
    `cms/internal/public/apitype`
    datatype2 `cms/internal/public/jsontype`
    
    "cms/internal/infra/store/datatype"
)

type FileSaveCommand struct {
    Name string
    Data io.Reader
}

type FileInfo struct {
    apitype.Base
    Name string              `json:"name"`
    URL  string              `json:"url"`
    Size datatype2.SafeInt64 `json:"size"`
}

type FileMeta struct {
    Name    string
    Size    datatype2.SafeInt64
    IsImage datatype.BoolInt8
}
