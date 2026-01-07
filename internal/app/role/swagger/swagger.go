package swagger

import (
    `dpcms/internal/app/role/internal/dto`
    `dpcms/internal/types`
)

type Role = types.ApiResult[dto.Role]
type RoleList = types.ApiPaginatedResult[dto.Role]
