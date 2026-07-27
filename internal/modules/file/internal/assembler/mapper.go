package assembler

import (
    `cms/internal/infra/persistence/gorm/model`
    `cms/internal/modules/file/internal/dto`
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
