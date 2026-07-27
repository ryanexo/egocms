package httpx

import (
    `net/http`
    
    `github.com/gin-gonic/gin`
)

type Router interface {
    MiddlewareRegistry
    
    Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup
    BasePath() string
    Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
    POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
    GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
    DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
    PATCH(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
    PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
    OPTIONS(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
    HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
    Any(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
    Match(methods []string, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
    StaticFile(relativePath, filepath string) gin.IRoutes
    StaticFileFS(relativePath, filepath string, fs http.FileSystem) gin.IRoutes
    Static(relativePath, root string) gin.IRoutes
    StaticFS(relativePath string, fs http.FileSystem) gin.IRoutes
}

type MiddlewareRegistry interface {
    Use(middleware ...gin.HandlerFunc) gin.IRoutes
}

type Route interface {
    Setup(Router)
}

type Middleware interface {
    Setup(MiddlewareRegistry)
}
