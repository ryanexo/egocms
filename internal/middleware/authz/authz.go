package authz

import `github.com/armon/go-radix`

type Builder struct {
    token TokenParser
    perm  PermissionChecker
}

func NewBuilder(tokenParser TokenParser, permChecker PermissionChecker) *Builder {
    return &Builder{token: tokenParser, perm: permChecker}
}

func (s *Builder) Basic() AccessControl {
    return &acl{
        perm:          s.perm,
        token:         s.token,
        object:        "",
        whitelistTree: radix.New(),
        permTree:      radix.New(),
    }
}

func (s *Builder) AccessControl(obj string) AccessControl {
    return &acl{
        perm:          s.perm,
        token:         s.token,
        object:        obj,
        whitelistTree: radix.New(),
        permTree:      radix.New(),
    }
}
