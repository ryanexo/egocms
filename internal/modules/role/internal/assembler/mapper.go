package assembler

import (
    `cms/internal/infra/persistence/gorm/model`
    `cms/internal/modules/role/internal/dto`
    `cms/internal/util/types`
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
