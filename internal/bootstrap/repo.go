package bootstrap

import (
    articleRepo `cms/internal/app/article/adapter`
    articleModelRepo `cms/internal/app/articlemodel/adapter`
    categoryRepo `cms/internal/app/category/adapter`
    config `cms/internal/app/setting/adapter`
    file `cms/internal/app/file/adapter`
    menuRepo `cms/internal/app/menu/adapter`
    permission `cms/internal/app/permission/adapter`
    roleRepo `cms/internal/app/role/adapter`
    tokenBlacklistRepo `cms/internal/app/token/adapter`
    userRepo `cms/internal/app/user/adapter`
    
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
