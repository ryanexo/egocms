package authz

import (
    "context"
    
    "github.com/gin-gonic/gin"
)

type contextKeyType struct{}

var contextKey contextKeyType

func GetCurrentUser(ctx context.Context) (User, error) {
    if ginContext, ok := ctx.(*gin.Context); ok {
        if user, found := ginContext.Get(contextKey); found {
            if authorizedUser, valid := user.(User); valid && authorizedUser != nil {
                return authorizedUser, nil
            }
        }
    }
    
    user, ok := ctx.Value(contextKey).(User)
    if !ok || user == nil {
        return nil, ErrAuthorized
    }
    
    return user, nil
}

func CurrentUser[T User](ctx context.Context) (T, error) {
    var zero T
    user, err := GetCurrentUser(ctx)
    if err != nil {
        return zero, err
    }
    typedUser, ok := user.(T)
    if !ok {
        return zero, ErrAuthorized
    }
    
    return typedUser, nil
}
