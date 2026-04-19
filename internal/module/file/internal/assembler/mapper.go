package assembler

import (
    `cms/internal/domain/file/internal/dto`
    `cms/internal/infra/persist/model`
    `cms/internal/util/types`
)

func ToFileInfoDTO(data *model.File, url string) *dto.FileInfo {
    return &dto.FileInfo{
        Base: types.Base{
            ID:        data.ID,
            CreatedAt: data.CreatedAt,
            UpdatedAt: data.UpdatedAt,
        },
        Name: data.OriginalName,
        URL:  url,
        Size: data.Size,
    }
}

func ToFileMetaDTO(data *model.File) *dto.FileMeta {
    return &dto.FileMeta{
        Name:    data.OriginalName,
        Size:    data.Size,
        IsImage: data.IsImage,
    }
}
