package controller

import (
    `github.com/gin-gonic/gin`
)

type controller interface {
    RegisterForServer(*gin.Engine)
}
