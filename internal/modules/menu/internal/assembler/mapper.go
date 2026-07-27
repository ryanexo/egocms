package assembler

import (
    `cms/internal/infra/persistence/datatype`
    model2 `cms/internal/infra/persistence/gorm/model`
    `cms/internal/modules/menu/internal/dto`
    `cms/internal/util/types`
)

func ToMenuCreateCommand(data *dto.MenuCreateParams) *model2.Menu {
    return &model2.Menu{
        Base:        model2.Base{},
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

func ToMenuUpdateCommand(data *dto.MenuUpdateParams) *model2.Menu {
    return &model2.Menu{
        Base: model2.Base{
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

func ToMenuDTO(data *model2.Menu) *dto.Menu {
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

func ToMenuListDTO(data []*model2.Menu) []*dto.Menu {
    result := make([]*dto.Menu, 0, len(data))
    for _, item := range data {
        result = append(result, ToMenuDTO(item))
    }
    return result
}
