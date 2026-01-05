package category

import (
    `dpcms/internal/app/dto`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
)

func BuildMenuCreateCommand(data *dto.MenuCreateParams) *model.Menu {
    return &model.Menu{
        ParentID:   data.ParentID,
        Type:       data.Type,
        Name:       data.Name,
        Sequence:   data.Sequence,
        Visible:    datatype.BoolInt8(1),
        URI:        data.URI,
        Resource:   data.Resource,
        Permission: data.Permission,
        Template:   data.Template,
        Remark:     data.Remark,
    }
}

func BuildMenuUpdateCommand(data *dto.MenuUpdateParams) *model.Menu {
    return &model.Menu{
        Base: model.Base{
            ID: data.ID,
        },
        Type:       data.Type,
        Name:       data.Name,
        Sequence:   data.Sequence,
        URI:        data.URI,
        Resource:   data.Resource,
        Permission: data.Permission,
        Template:   data.Template,
        Remark:     data.Remark,
    }
}

func BuildMenuDTO(data *model.Menu) *dto.Menu {
    return &dto.Menu{
        Base: dto.Base{
            ID:        data.ID,
            CreatedAt: data.CreatedAt,
            UpdatedAt: data.UpdatedAt,
        },
        ParentID:   data.ParentID,
        Type:       data.Type,
        Name:       data.Name,
        Sequence:   data.Sequence,
        Visible:    &data.Visible,
        URI:        data.URI,
        Resource:   data.Resource,
        Permission: data.Permission,
        Template:   data.Template,
        Remark:     data.Remark,
    }
}

func BuildMenuListDTO(data []*model.Menu) []*dto.Menu {
    result := make([]*dto.Menu, 0, len(data))
    for _, item := range data {
        result = append(result, BuildMenuDTO(item))
    }
    return result
}
