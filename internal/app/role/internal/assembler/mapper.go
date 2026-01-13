package assembler

import (
    `dpcms/internal/app/role/internal/dto`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/types`
)

func ToRoleDTO(data *model.Role) *dto.Role {
    return &dto.Role{
        Base: types.Base{
            ID:        data.ID,
            CreatedAt: data.CreatedAt,
            UpdatedAt: data.UpdatedAt,
        },
        Name:        data.Name,
        Description: data.Description,
    }
}

func ToRoleListDTO(data []*model.Role) []*dto.Role {
    result := make([]*dto.Role, 0, len(data))
    for _, item := range data {
        result = append(result, ToRoleDTO(item))
    }
    return result
}
