package assembler

import (
    `dpcms/internal/app/permission/internal/dto`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/types`
)

func BuildPermissionCreateCommand(data *dto.PermissionCreateParams) *model.Permission {
    return &model.Permission{
        MenuID:      data.MenuID,
        Name:        data.Name,
        Description: data.Description,
        Resource:    data.Resource,
        Action:      data.Action,
    }
}

func BuildPermissionUpdateCommand(data *dto.PermissionUpdateParams) *model.Permission {
    result := BuildPermissionCreateCommand(&data.PermissionCreateParams)
    result.ID = data.ID
    return result
}

func BuildPermissionDTO(data *model.Permission) *dto.Permission {
    return &dto.Permission{
        Base: types.Base{
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

func BuildPermissionListDTO(data []*model.Permission) []*dto.Permission {
    result := make([]*dto.Permission, 0, len(data))
    for _, item := range data {
        result = append(result, BuildPermissionDTO(item))
    }
    return result
}
