package erroz

import (
    `strings`
    
    `dpcms/internal/erroz/internal/format`
    `dpcms/internal/erroz/internal/mod`
    `dpcms/internal/erroz/internal/typ`
)

const (
    oK = iota
    unknown
    auth
    notFound
    invalidState
    parameter
    conflict
)

const (
    server = iota
    client
    user
    role
    category
    menu
    article
    articleModel
    token
    permission
)

var (
    ModuleServer       = mod.Module{Code: server}
    ModuleClient       = mod.Module{Code: client}
    ModuleUser         = mod.Module{Code: user}
    ModuleRole         = mod.Module{Code: role}
    ModuleCategory     = mod.Module{Code: category}
    ModuleMenu         = mod.Module{Code: menu}
    ModuleArticle      = mod.Module{Code: article}
    ModuleArticleModel = mod.Module{Code: articleModel}
    ModuleToken        = mod.Module{Code: token}
    ModulePermission   = mod.Module{Code: permission}
)

var (
    TypOK           = typ.Typ{Code: oK}
    TypUnknown      = typ.Typ{Code: unknown}
    TypAuth         = typ.Typ{Code: auth}
    TypNotFound     = typ.Typ{Code: notFound}
    TypInvalidState = typ.Typ{Code: invalidState}
    TypParameter    = typ.Typ{Code: parameter}
    TypConflict     = typ.Typ{Code: conflict}
)

func Code(modCode mod.Module, typCode typ.Typ, code int64) string {
    return strings.Join([]string{
        modCode.String(),
        typCode.String(),
        format.Code(code, 6),
    }, "-")
}
