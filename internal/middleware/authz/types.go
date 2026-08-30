package authz

import (
    "context"
    
    "github.com/gin-gonic/gin"
)

type User interface {
    UserID() uint64
    Role() string
}

type TokenParser interface {
    Parse(ctx context.Context, token string) (User, error)
}

type PermissionChecker interface {
    Check(ctx context.Context, subject string, object string, action string) (bool, error)
}

type Factory struct {
    token   TokenParser
    checker PermissionChecker
}

func New(tokenParser TokenParser, permissionChecker PermissionChecker) *Factory {
    return &Factory{
        token:   tokenParser,
        checker: permissionChecker,
    }
}

type Middleware struct {
    token      TokenParser
    checker    PermissionChecker
    object     string
    action     string
    public     bool
    permission bool
}

func (factory *Factory) Resource(object string) Middleware {
    return Middleware{
        token:   factory.token,
        checker: factory.checker,
        object:  object,
    }
}

func (middleware Middleware) Public() Middleware {
    middleware.public = true
    middleware.permission = false
    middleware.action = ""
    return middleware
}

func (middleware Middleware) Permission(action string) Middleware {
    middleware.public = false
    middleware.permission = true
    middleware.action = action
    return middleware
}

func (middleware Middleware) Wrap(handler gin.HandlerFunc) gin.HandlerFunc {
    if middleware.public {
        return handler
    }
    
    return func(ctx *gin.Context) {
        user, ok := middleware.authenticate(ctx)
        if !ok {
            return
        }
        
        if middleware.permission && !middleware.authorize(ctx, user) {
            return
        }
        
        handler(ctx)
    }
}
