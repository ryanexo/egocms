package authz

import (
    "context"
    "strings"
    
    "github.com/gin-gonic/gin"
)

func (middleware Middleware) authenticate(ctx *gin.Context) (User, bool) {
    credential := ctx.GetHeader("Authorization")
    token, found := strings.CutPrefix(credential, "Bearer ")
    if !found || token == "" || middleware.token == nil {
        abort(ctx, ErrAuthorized)
        return nil, false
    }
    
    user, err := middleware.token.Parse(ctx, token)
    if err != nil || user == nil {
        abort(ctx, ErrAuthorized)
        return nil, false
    }
    
    ctx.Set(contextKey, user)
    ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), contextKey, user))
    return user, true
}

func (middleware Middleware) authorize(ctx *gin.Context, user User) bool {
    if middleware.checker == nil || middleware.object == "" || middleware.action == "" {
        abort(ctx, ErrAccessDenied)
        return false
    }
    
    allowed, err := middleware.checker.Check(ctx, user.Role(), middleware.object, middleware.action)
    if err != nil || !allowed {
        abort(ctx, ErrAccessDenied)
        return false
    }
    
    return true
}

func abort(ctx *gin.Context, err error) {
    _ = ctx.Error(err)
    ctx.Abort()
}
