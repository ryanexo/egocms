package provider

import (
    article `cms/internal/app/article/service`
    articlemodel `cms/internal/app/articlemodel/service`
    `cms/internal/app/category`
    fileService `cms/internal/app/file/service`
    menu `cms/internal/app/menu/service`
    permission `cms/internal/app/permission/service`
    role `cms/internal/app/role/service`
    setting `cms/internal/app/setting/service`
    token `cms/internal/app/token/service`
    user `cms/internal/app/user/service`
    
    "github.com/google/wire"
)

var ServiceProvider = wire.NewSet(
    article.NewArticleService,
    articlemodel.NewArticleModelService,
    category.NewCategoryService,
    menu.NewMenuService,
    user.NewUserService,
    role.NewRoleService,
    token.NewTokenService,
    setting.NewSettingService,
    fileService.NewFileService,
    fileService.NewFileDriverService,
    permission.NewPermissionService,
)
