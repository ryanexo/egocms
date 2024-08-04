package service

import (
    "context"
    "errors"
    
    `dpcms/api/repository`
    `dpcms/erroz`
    `dpcms/model`
    `dpcms/packages/password`
    "gorm.io/gorm"
)

type User struct {
    repo repository.Repositories
}

func NewUserService(repo repository.Repositories) *User {
    return &User{repo: repo}
}

func (srv User) Create(ctx context.Context, u *model.User) error {
    findResult, err := srv.repo.User.FindByUsernameOrEmail(ctx, u.Username, u.Email)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return err
    }
    if findResult != nil {
        if findResult.Username == u.Username {
            return erroz.ErrUsernameExists.ToError()
        }
        if findResult.Email == u.Email {
            return erroz.ErrEMailExists.ToError()
        }
    }
    passwd, err := password.Make(u.Password)
    if err != nil {
        return err
    }
    u.Password = passwd
    return srv.repo.User.Create(ctx, u)
}

func (srv User) FindByID(ctx context.Context, id uint) (*model.User, error) {
    result, err := srv.repo.User.FindByID(ctx, id)
    if errors.Is(err, gorm.ErrRecordNotFound) {
        err = erroz.ErrUserIDNotExists.ToError()
    }
    return result, err
}

func (srv User) FindByName(ctx context.Context, name string) (*model.User, error) {
    result, err := srv.repo.User.FindByUsername(ctx, name)
    if errors.Is(err, gorm.ErrRecordNotFound) {
        err = erroz.ErrUsernameNotExists.ToError()
    }
    return result, err
}

func (srv User) UpdatePassword(ctx context.Context, id uint, pwd string) error {
    encrypted, err := password.Make(pwd)
    if err != nil {
        return err
    }
    _, err = srv.repo.User.UpdateUserInfo(ctx, id, map[string]any{"password": encrypted})
    return err
}
