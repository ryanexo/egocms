package biz

import (
    `GoBlog/internal/server/biz/controller`
    `GoBlog/internal/server/biz/repository`
    `GoBlog/internal/server/biz/service`
    `github.com/google/wire`
)

var ControllerSet = wire.NewSet(
    controller.NewUserController,
)

var ServiceWireSet = wire.NewSet(
    service.NewUserService,
)

var RepoWireSet = wire.NewSet(
    repository.NewUserRepo,
)
