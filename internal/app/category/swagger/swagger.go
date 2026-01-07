package swagger

import (
    `dpcms/internal/app/category/internal/dto`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/types`
)

type CategoryList = types.ApiPaginatedResult[model.Category]
type Category = types.ApiResult[dto.Category]
