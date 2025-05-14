package service

import (
    "context"
    
    `dpcms/api/infra`
    `dpcms/erroz`
    `dpcms/model`
    `dpcms/model/query`
    `dpcms/packages/password`
)

type UserService struct {
    infra *infra.Infra
    query *query.Query
}

func NewUserService(infra *infra.Infra) *UserService {
    return &UserService{infra: infra, query: query.Use(infra.DB)}
}

func (srv *UserService) Create(ctx context.Context, u *model.User) error {
    q := srv.query.User
    findResult, err := q.WithContext(ctx).Where(q.Username.Eq(u.Username)).Or(q.Email.Eq(u.Email)).First()
    if err != nil {
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
    return q.WithContext(ctx).Create(u)
}

func (srv *UserService) FindByID(ctx context.Context, id uint) (*model.User, error) {
    q := srv.query.User
    result, err := q.WithContext(ctx).Where(q.ID.Eq(id)).First()
    if err != nil {
        return nil, err
    }
    return result, err
}

func (srv *UserService) FindByName(ctx context.Context, name string) (*model.User, error) {
    q := srv.query.User
    result, err := q.WithContext(ctx).Where(q.Username.Eq(name)).First()
    if err != nil {
        return nil, err
    }
    return result, err
}

func (srv *UserService) UpdatePassword(ctx context.Context, id uint, pwd string) error {
    q := srv.query.User
    finalPassword, err := password.Make(pwd)
    if err != nil {
        return err
    }
    _, err = q.WithContext(ctx).Where(q.ID.Eq(id)).Update(q.Password, finalPassword)
    return err
}
