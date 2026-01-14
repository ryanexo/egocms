package controller

import permission `cms/internal/app/permission/service`

type PermissionController struct {
    permSrv *permission.PermissionService
}

func NewPermissionController(permSrv *permission.PermissionService) *PermissionController {
    return &PermissionController{permSrv: permSrv}
}
