package category

import (
    `dpcms/internal/app/dto`
    `dpcms/internal/infra/persistence/model`
    
    `github.com/jinzhu/copier`
)

func ToCategoryCreateCommand(data *dto.CategoryCreateParams) (*model.Category, error) {
    result := &model.Category{}
    err := copier.Copy(&result, data)
    if err != nil {
        return nil, err
    }
    return result, nil
}

func ToCategoryUpdateCommand(data *dto.CategoryUpdateParams) (*model.Category, error) {
    result := &model.Category{}
    err := copier.Copy(&result, data)
    if err != nil {
        return nil, err
    }
    return result, nil
}

func ToCategoryDTO(category *model.Category) (*dto.Category, error) {
    var result *dto.Category
    if err := copier.Copy(&result, category); err != nil {
        return nil, err
    }
    return result, nil
}

func ToCategoryListDTO(categories []*model.Category) ([]*dto.Category, error) {
    var result []*dto.Category
    if err := copier.Copy(&result, categories); err != nil {
        return nil, err
    }
    return result, nil
}
