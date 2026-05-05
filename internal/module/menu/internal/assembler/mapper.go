package assembler

import (
    `cms/internal/domain/menu/internal/dto`
    `cms/internal/infra/persistence/datatype`
    `cms/internal/infra/persistence/model`
    `cms/internal/util/types`
)

func ToMenuCreateCommand(data *dto.MenuCreateParams) *model.Menu {
    return &model.Menu{
        Base:        model.Base{},
        ParentID:    data.ParentID,
        Type:        data.Type,
        Name:        data.Name,
        Affix:       data.Affix,
        Icon:        data.Icon,
        ExternalURL: data.ExternalURL,
        Sequence:    data.Sequence,
        Visible:     datatype.BoolInt8(1),
        URI:         data.URI,
        Template:    data.Template,
        Remark:      data.Remark,
        Action:      nil,
    }
}

func ToMenuUpdateCommand(data *dto.MenuUpdateParams) *model.Menu {
    return &model.Menu{
        Base: model.Base{
            ID: data.ID,
        },
        Type:        data.Type,
        Name:        data.Name,
        Affix:       data.Affix,
        Icon:        data.Icon,
        ExternalURL: data.ExternalURL,
        Sequence:    data.Sequence,
        URI:         data.URI,
        Template:    data.Template,
        Remark:      data.Remark,
    }
}

func ToMenuDTO(data *model.Menu) *dto.Menu {
    return &dto.Menu{
        Base: types.Base{
            ID:        data.ID,
            CreatedAt: data.CreatedAt,
            UpdatedAt: data.UpdatedAt,
        },
        ParentID:    data.ParentID,
        Type:        data.Type,
        Name:        data.Name,
        Affix:       data.Affix,
        Icon:        data.Icon,
        ExternalURL: data.ExternalURL,
        Sequence:    data.Sequence,
        Visible:     data.Visible,
        URI:         data.URI,
        Template:    data.Template,
        Remark:      data.Remark,
    }
}

func ToMenuListDTO(data []*model.Menu) []*dto.Menu {
    result := make([]*dto.Menu, 0, len(data))
    for _, item := range data {
        result = append(result, ToMenuDTO(item))
    }
    return result
}
