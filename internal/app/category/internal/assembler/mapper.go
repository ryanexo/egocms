package assembler

import (
    `cms/internal/app/category/internal/dto`
    `cms/internal/infra/persistence/model`
    `cms/internal/types`
)

func ToCategoryCreateCommand(data *dto.CategoryCreateParams) *model.Category {
    result := &model.Category{
        ParentID: data.ParentID,
        Sequence: data.Sequence,
        Name:     data.Name,
        Path:     data.Path,
        Type:     data.Type,
        SEO: &model.CategorySeo{
            Title:       data.SEO.Title,
            Keywords:    data.SEO.Keywords,
            Description: data.SEO.Description,
        },
    }
    return result
}

func ToCategoryUpdateCommand(data *dto.CategoryUpdateParams) *model.Category {
    result := &model.Category{
        Base: model.Base{
            ID: data.ID,
        },
        Sequence: data.Sequence,
        Name:     data.Name,
        Path:     data.Path,
        Type:     data.Type,
        SEO: &model.CategorySeo{
            CategoryID:  data.ID,
            Title:       data.SEO.Title,
            Keywords:    data.SEO.Keywords,
            Description: data.SEO.Description,
        },
    }
    if data.Visible != nil {
        result.Visible = *data.Visible
    }
    return result
}

func ToCategoryDTO(data *model.Category) *dto.Category {
    return &dto.Category{
        Base: types.Base{
            ID:        data.ID,
            CreatedAt: data.CreatedAt,
            UpdatedAt: data.UpdatedAt,
        },
        ParentID: data.ParentID,
        Sequence: data.Sequence,
        Name:     data.Name,
        Path:     data.Path,
        Type:     data.Type,
        Visible:  &data.Visible,
        SEO: &dto.CategorySEO{
            Title:       data.SEO.Title,
            Keywords:    data.SEO.Keywords,
            Description: data.SEO.Description,
        },
    }
}

func ToCategoryListDTO(data []*model.Category) []*dto.Category {
    result := make([]*dto.Category, 0, len(data))
    for _, item := range data {
        result = append(result, ToCategoryDTO(item))
    }
    return result
}
