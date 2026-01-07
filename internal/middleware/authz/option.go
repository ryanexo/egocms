package authz

import (
    `strings`
    
    `github.com/gin-gonic/gin`
)

func WithWhitelist(path ...string) Option {
    return func(s *acl) {
        for _, p := range path {
            s.whitelistTree.Insert(p, nil)
        }
    }
}

func WithRouterWhitelist(path ...string) RouterOption {
    return func(g *gin.RouterGroup, s *acl) {
        basePath := g.BasePath()
        for _, p := range path {
            finalPath, _ := strings.CutPrefix(p, "/")
            s.whitelistTree.Insert(basePath+"/"+finalPath, nil)
        }
    }
}

func WithPermission(path string, perm string) Option {
    return func(s *acl) {
        s.permTree.Insert(path, perm)
    }
}

func WithRouterPermission(path string, perm string) RouterOption {
    return func(g *gin.RouterGroup, s *acl) {
        finalPath, _ := strings.CutPrefix(path, "/")
        s.permTree.Insert(finalPath+"/"+perm, perm)
    }
}
