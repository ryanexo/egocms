package bootstrap

import (
    articleRepo `dpcms/internal/app/article/adapter`
    articleModelRepo `dpcms/internal/app/articlemodel/adapter`
    categoryRepo `dpcms/internal/app/category/adapter`
    config `dpcms/internal/app/setting/adapter`
    file `dpcms/internal/app/file/adapter`
    menuRepo `dpcms/internal/app/menu/adapter`
    permission `dpcms/internal/app/permission/adapter`
    roleRepo `dpcms/internal/app/role/adapter`
    tokenBlacklistRepo `dpcms/internal/app/token/adapter`
    userRepo `dpcms/internal/app/user/adapter`
    
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
