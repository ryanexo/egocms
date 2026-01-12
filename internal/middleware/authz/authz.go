package authz

import `github.com/armon/go-radix`

type Factory struct {
    token TokenParser
    perm  PermissionChecker
}

func NewFactory(tokenParser TokenParser, permChecker PermissionChecker) *Factory {
    return &Factory{token: tokenParser, perm: permChecker}
}

func (s *Factory) Basic() AccessControl {
    return &acl{
        perm:          s.perm,
        token:         s.token,
        object:        "",
        whitelistTree: radix.New(),
        permTree:      radix.New(),
    }
}

func (s *Factory) AccessControl(obj string) AccessControl {
    return &acl{
        perm:          s.perm,
        token:         s.token,
        object:        obj,
        whitelistTree: radix.New(),
        permTree:      radix.New(),
    }
}
