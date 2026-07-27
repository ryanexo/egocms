package controller

import permission `cms/internal/modules/permission/service`

type PermissionController struct {
    permSrv *permission.PermissionService
}

func NewPermissionController(permSrv *permission.PermissionService) *PermissionController {
    return &PermissionController{permSrv: permSrv}
}
