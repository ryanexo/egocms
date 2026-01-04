package category

import (
    `dpcms/internal/app/dto`
    `dpcms/internal/infra/persistence/model`
    
    `github.com/jinzhu/copier`
)

func BuildMenuCreateCommand(data *dto.MenuCreateParams) (*model.Menu, error) {
    result := &model.Menu{}
    err := copier.Copy(&result, data)
    if err != nil {
        return nil, err
    }
    return result, nil
}

func BuildMenuUpdateCommand(data *dto.MenuUpdateParams) (*model.Menu, error) {
    result := &model.Menu{}
    err := copier.Copy(&result, data)
    if err != nil {
        return nil, err
    }
    return result, nil
}

func BuildMenuDTO(category *model.Menu) (*dto.Menu, error) {
    var result *dto.Menu
    if err := copier.Copy(&result, category); err != nil {
        return nil, err
    }
    return result, nil
}

func BuildMenuListDTO(categories []*model.Menu) ([]*dto.Menu, error) {
    var result []*dto.Menu
    if err := copier.Copy(&result, categories); err != nil {
        return nil, err
    }
    return result, nil
}
