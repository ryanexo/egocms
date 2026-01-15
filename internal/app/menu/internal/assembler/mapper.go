package assembler

import (
    `cms/internal/app/menu/internal/dto`
    `cms/internal/infra/persistence/datatype`
    `cms/internal/infra/persistence/model`
    `cms/internal/util/types`
)

func ToMenuCreateCommand(data *dto.MenuCreateParams) *model.Menu {
    return &model.Menu{
        ParentID: data.ParentID,
        Type:     data.Type,
        Name:     data.Name,
        Sequence: data.Sequence,
        Visible:  datatype.BoolInt8(1),
        URI:      data.URI,
        Template: data.Template,
        Remark:   data.Remark,
    }
}

func ToMenuUpdateCommand(data *dto.MenuUpdateParams) *model.Menu {
    return &model.Menu{
        Base: model.Base{
            ID: data.ID,
        },
        Type:     data.Type,
        Name:     data.Name,
        Sequence: data.Sequence,
        URI:      data.URI,
        Template: data.Template,
        Remark:   data.Remark,
    }
}

func ToMenuDTO(data *model.Menu) *dto.Menu {
    return &dto.Menu{
        Base: types.Base{
            ID:        data.ID,
            CreatedAt: data.CreatedAt,
            UpdatedAt: data.UpdatedAt,
        },
        ParentID: data.ParentID,
        Type:     data.Type,
        Name:     data.Name,
        Sequence: data.Sequence,
        Visible:  &data.Visible,
        URI:      data.URI,
        Template: data.Template,
        Remark:   data.Remark,
    }
}

func ToMenuListDTO(data []*model.Menu) []*dto.Menu {
    result := make([]*dto.Menu, 0, len(data))
    for _, item := range data {
        result = append(result, ToMenuDTO(item))
    }
    return result
}
