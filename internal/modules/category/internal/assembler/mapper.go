package assembler

import (
    model2 `cms/internal/infra/persistence/gorm/model`
    `cms/internal/modules/category/internal/dto`
    `cms/internal/util/types`
)

func ToCategoryCreateCommand(data *dto.CategoryCreateParams) *model2.Category {
    result := &model2.Category{
        ParentID: data.ParentID,
        Sequence: data.Sequence,
        Name:     data.Name,
        Path:     data.Path,
        Type:     data.Type,
        SEO: &model2.CategorySeo{
            Title:       data.SEO.Title,
            Keywords:    data.SEO.Keywords,
            Description: data.SEO.Description,
        },
    }
    return result
}

func ToCategoryUpdateCommand(data *dto.CategoryUpdateParams) *model2.Category {
    result := &model2.Category{
        Base: model2.Base{
            ID: data.ID,
        },
        Sequence: data.Sequence,
        Name:     data.Name,
        Path:     data.Path,
        Type:     data.Type,
        SEO: &model2.CategorySeo{
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

func ToCategoryDTO(data *model2.Category) *dto.Category {
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

func ToCategoryListDTO(data []*model2.Category) []*dto.Category {
    result := make([]*dto.Category, 0, len(data))
    for _, item := range data {
        result = append(result, ToCategoryDTO(item))
    }
    return result
}
