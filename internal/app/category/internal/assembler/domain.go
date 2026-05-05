package assembler

import (
    `strings`
    
    `cms/internal/app/category/domain`
    `cms/internal/infra/persistence/model`
)

func Dehydrate(v domain.Category) model.Category {
    return model.Category{
        Base: model.Base{
            ID:        v.ID(),
            CreatedAt: v.CreatedAt(),
            UpdatedAt: v.UpdatedAt(),
            DeletedAt: v.DeletedAt(),
        },
        ParentID: v.ParentID(),
        Sequence: v.Sequence(),
        Name:     v.Name(),
        Path:     v.Path(),
        SEO: &model.CategorySeo{
            Title:       v.SEO().Title(),
            Keywords:    strings.Join(v.SEO().Keywords(), ","),
            Description: v.SEO().Description(),
        },
    }
}

func Hydrate(v model.Category) domain.Category {
    return domain.RestoreCategory(
        v.ID,
        v.Name,
        v.Path,
        v.Sequence,
        domain.NewCategorySeo(
            v.SEO.Title,
            strings.Split(v.SEO.Keywords, ","),
            v.SEO.Description,
        ),
        v.CreatedAt,
        v.UpdatedAt,
        v.DeletedAt,
    )
}
