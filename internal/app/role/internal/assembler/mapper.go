package assembler

import (
    `dpcms/internal/app/role/internal/dto`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/types`
)

func BuildRoleDTO(data *model.Role) *dto.Role {
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

func BuildRoleListDTO(data []*model.Role) []*dto.Role {
    result := make([]*dto.Role, 0, len(data))
    for _, item := range data {
        result = append(result, BuildRoleDTO(item))
    }
    return result
}
