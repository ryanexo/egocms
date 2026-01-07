package swagger

import (
    `dpcms/internal/app/user/internal/dto`
    `dpcms/internal/types`
)

type UserList = types.ApiPaginatedResult[dto.User]
type User = types.ApiResult[dto.User]
