package controller

import (
    `context`

    `GoBlog/internal/server/biz/service`
    `GoBlog/internal/server/dep`
    `GoBlog/internal/server/entity`
    `GoBlog/internal/server/erroz`
    `GoBlog/internal/server/pkg/password`
    `GoBlog/internal/server/pkg/token`
    `github.com/gin-gonic/gin`
)

type UserController interface {
    controller
}

type userController struct {
    service service.UserService
    deps    *dep.Dep
}

func (ctrl userController) RegisterForServer(s *gin.Engine) {
    g := s.Group("/user")
    g.POST("/register", ctrl.Register)
    g.POST("/login", ctrl.Login)
}

func NewUserController(s service.UserService, deps *dep.Dep) UserController {
    return userController{service: s, deps: deps}
}

func (ctrl userController) Register(ctx *gin.Context) {
    u := new(entity.UserCreatingForm)
    if err := ctx.ShouldBind(u); err != nil {
        erroz.Resolve(ctx, err)
        return
    }
    if u.Password != u.RePassword {
        erroz.ErrRePasswordNotEq.Apply(ctx)
    }
    result := &entity.User{
        Username: u.Username,
        Password: u.Password,
        Email:    u.Email,
    }
    if err := ctrl.service.Create(context.Background(), result); err != nil {
        erroz.Resolve(ctx, err)
        return
    }
    t, err := ctrl.deps.Token.CreateUserToken(token.UserTokenClaims{ID: result.ID})
    if err != nil {
        erroz.Resolve(ctx, err)
    }
    erroz.OK.WithOption(erroz.WithData(t)).Apply(ctx)
}

func (ctrl userController) Login(ctx *gin.Context) {
    u := new(entity.UserCredential)
    if err := ctx.ShouldBind(u); err != nil {
        erroz.Resolve(ctx, err)
        return
    }
    result, err := ctrl.service.FindByName(context.Background(), u.Username)
    if err != nil {
        erroz.Resolve(ctx, err)
    }
    if !password.Equal(result.Password, u.Password) {
        erroz.ErrWrongPassword.Apply(ctx)
    }
    t, err := ctrl.deps.Token.CreateUserToken(token.UserTokenClaims{ID: 0})
    if err != nil {
        erroz.Resolve(ctx, err)
    }
    erroz.OK.WithOption(erroz.WithData(t)).Apply(ctx)
}
