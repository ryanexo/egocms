package httpserver

import "github.com/gin-gonic/gin"

func Handler(fn func(ctx *gin.Context) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := fn(ctx); err != nil {
			_ = ctx.Error(err)
		}
	}
}
