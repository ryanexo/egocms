package authzutil

import (
    `context`
    
    `cms/internal/middleware/authz`
    `cms/internal/public/model`
)

func GetAuthorizedUser(ctx context.Context) (*model.User, error) {
    return authz.GetCurrentUser[*model.User](ctx)
}
