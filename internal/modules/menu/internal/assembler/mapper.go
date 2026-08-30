package assembler

import (
    `cms/internal/infra/store/datatype`
    `cms/internal/infra/store/modeltype`
    `cms/internal/modules/menu/internal/dto`
    `cms/internal/public/apitype`
    `cms/internal/public/model`
)

func ToMenuCreateCommand(data *dto.MenuCreateParams) *model.Menu {
    return &model.Menu{
        Base:        modeltype.Base{},
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
        Base: modeltype.Base{
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
        Base: apitype.Base{
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
