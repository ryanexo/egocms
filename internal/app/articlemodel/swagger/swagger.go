package swagger

import (
    `dpcms/internal/app/articlemodel/internal/dto`
    `dpcms/internal/types`
)

type ArticleModel = types.ApiResult[dto.ArticleModel]
type ArticleModelListResult = types.ApiPaginatedResult[dto.ArticleModel]
