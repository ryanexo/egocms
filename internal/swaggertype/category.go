package swaggertype

import (
    `dpcms/internal/app/dto`
    `dpcms/internal/infra/persistence/model`
)

type CategoryList PaginatedResult[model.Category]
type Category Result[dto.Category]
