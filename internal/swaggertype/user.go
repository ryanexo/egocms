package swaggertype

import "dpcms/internal/app/dto"

type UserList = PaginatedResult[dto.User]
type User = Result[dto.User]
