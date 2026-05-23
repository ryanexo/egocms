package provider

import (
    articleRepo `cms/internal/app/article/provider`
    categoryRepo `cms/internal/app/category`
    articleModelRepo `cms/internal/app/contenttype/provider`
    file `cms/internal/app/file/provider`
    menuRepo `cms/internal/app/menu/provider`
    permission `cms/internal/app/permission/provider`
    roleRepo `cms/internal/app/role/provider`
    config `cms/internal/app/setting/provider`
    tokenBlacklistRepo `cms/internal/app/token/provider`
    userRepo `cms/internal/app/user/provider`
    
    `github.com/google/wire`
)

var RepoProvider = wire.NewSet(
    articleRepo.NewArticleRepo,
    articleModelRepo.NewArticleModelRepo,
    categoryRepo.NewCategoryRepo,
    menuRepo.NewMenuRepo,
    roleRepo.NewRoleRepo,
    userRepo.NewUserRepo,
    tokenBlacklistRepo.NewTokenBlacklistRepo,
    file.NewFileRepo,
    config.NewConfigRepo,
    permission.NewPermissionRepo,
)
