package authz

import (
    `context`
    
    `cms/internal/constant`
)

func GetCurrentUser[T any](ctx context.Context) (T, error) {
    var zero T
    val := ctx.Value(constant.RequestUserKey)
    if val == nil {
        return zero, ErrAuthorized.ToError()
    }
    user, ok := val.(T)
    if !ok {
        return zero, ErrAuthorized.ToError()
    }
    return user, nil
}
