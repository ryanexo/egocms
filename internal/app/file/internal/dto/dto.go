package dto

import (
    `io`
    
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/types`
)

type FileSaveCommand struct {
    Filename string
    Data     io.Reader
}

type FileInfo struct {
    types.Base
    Filename string             `json:"filename"`
    URL      string             `json:"url"`
    Size     datatype.SafeInt64 `json:"size"`
}
