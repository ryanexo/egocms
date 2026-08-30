package assembler

import (
    `cms/internal/modules/permission/internal/dto`
    `cms/internal/public/apitype`
    `cms/internal/public/model`
)

func ToPermissionCreateCommand(data *dto.PermissionCreateParams) *model.Permission {
    return &model.Permission{
        MenuID:      data.MenuID,
        Name:        data.Name,
        Description: data.Description,
        Resource:    data.Resource,
        Action:      data.Action,
    }
}

func ToPermissionUpdateCommand(data *dto.PermissionUpdateParams) *model.Permission {
    result := ToPermissionCreateCommand(&data.PermissionCreateParams)
    result.ID = data.ID
    return result
}

func ToPermissionDTO(data *model.Permission) *dto.Permission {
    return &dto.Permission{
        Base: apitype.Base{
            ID:        data.ID,
            CreatedAt: data.CreatedAt,
            UpdatedAt: data.UpdatedAt,
        },
        MenuID:      data.MenuID,
        Name:        data.Name,
        Description: data.Description,
        Resource:    data.Resource,
        Action:      data.Action,
    }
}

func ToPermissionListDTO(data []*model.Permission) []*dto.Permission {
    result := make([]*dto.Permission, 0, len(data))
    for _, item := range data {
        result = append(result, ToPermissionDTO(item))
    }
    return result
}
