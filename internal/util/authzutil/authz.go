package authzutil

import (
    `context`
    
    `cms/internal/infra/persistence/model`
    `cms/internal/middleware/authz`
)

func GetAuthorizedUser(ctx context.Context) (*model.User, error) {
    return authz.GetCurrentUser[*model.User](ctx)
}
