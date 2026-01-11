package erroz

import (
    `strings`
    
    `dpcms/internal/erroz/internal/format`
    `dpcms/internal/erroz/internal/mod`
    `dpcms/internal/erroz/internal/typ`
)

const (
    typeOK = iota
    typeUnknown
    typeAuth
    typeNotFound
    typeInvalidState
    typeParameter
    typeConflict
)

const (
    modServer = iota
    modClient
    modUser
    modRole
    modToken
    modPermission
    modCategory
    modMenu
    modArticle
    modArticleModel
    modConfig
    modFile
)

var (
    ModuleServer       = mod.Module{Code: modServer}
    ModuleClient       = mod.Module{Code: modClient}
    ModuleUser         = mod.Module{Code: modUser}
    ModuleRole         = mod.Module{Code: modRole}
    ModuleCategory     = mod.Module{Code: modCategory}
    ModuleMenu         = mod.Module{Code: modMenu}
    ModuleConfig       = mod.Module{Code: modConfig}
    ModuleFile         = mod.Module{Code: modFile}
    ModuleArticle      = mod.Module{Code: modArticle}
    ModuleArticleModel = mod.Module{Code: modArticleModel}
    ModuleToken        = mod.Module{Code: modToken}
    ModulePermission   = mod.Module{Code: modPermission}
)

var (
    TypOK           = typ.Typ{Code: typeOK}
    TypUnknown      = typ.Typ{Code: typeUnknown}
    TypAuth         = typ.Typ{Code: typeAuth}
    TypNotFound     = typ.Typ{Code: typeNotFound}
    TypInvalidState = typ.Typ{Code: typeInvalidState}
    TypParameter    = typ.Typ{Code: typeParameter}
    TypConflict     = typ.Typ{Code: typeConflict}
)

func Code(modCode mod.Module, typCode typ.Typ, code int64) string {
    return strings.Join([]string{
        modCode.String(),
        typCode.String(),
        format.Code(code, 6),
    }, "-")
}
