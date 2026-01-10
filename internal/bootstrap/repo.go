package bootstrap

import (
    articleRepo `dpcms/internal/app/article/adapter`
    articleModelRepo `dpcms/internal/app/articlemodel/adapter`
    categoryRepo `dpcms/internal/app/category/adapter`
    menuRepo `dpcms/internal/app/menu/adapter`
    roleRepo `dpcms/internal/app/role/adapter`
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
)
