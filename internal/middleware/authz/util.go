package authz

import (
    `context`
    
    `dpcms/internal/constant`
    `dpcms/internal/erroz`
)

func GetCurrentUser[T any](ctx context.Context) (T, error) {
    var zero T
    val := ctx.Value(constant.RequestUserKey)
    if val == nil {
        return zero, erroz.Unauthorized.ToError()
    }
    user, ok := val.(T)
    if !ok {
        return zero, erroz.Unauthorized.ToError()
    }
    return user, nil
}
