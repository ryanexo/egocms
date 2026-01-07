package contextutil

import (
    `context`
    
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/middleware/authz`
)

func GetAuthorizedUser(ctx context.Context) (*model.User, error) {
    return authz.GetCurrentUser[*model.User](ctx)
}
