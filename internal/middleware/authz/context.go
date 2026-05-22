package authz

import (
	"context"

	"github.com/gin-gonic/gin"
)

type contextKeyType struct{}

var contextKey = contextKeyType{}

func setCurrentUser(ctx *gin.Context, user User) {
	ctx.Set(contextKey, user)
}

func GetCurrentUser(ctx context.Context) (User, error) {
	val := ctx.Value(contextKey)
	if val == nil {
		return nil, ErrAuthorized
	}
	user, ok := val.(User)
	if !ok {
		return nil, ErrAuthorized
	}
	return user, nil
}
