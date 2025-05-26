package controller

import `dpcms/internal/http/service`

type MenuController struct {
    srv *service.Services
}

func NewMenuController(srv *service.Services) *MenuController {
    return &MenuController{srv}
}
