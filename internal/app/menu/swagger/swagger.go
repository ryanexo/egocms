package swagger

import (
    `dpcms/internal/app/menu/internal/dto`
    `dpcms/internal/types`
)

type Menu = types.ApiResult[dto.Menu]
type MenuList = types.ApiPaginatedResult[dto.Menu]
