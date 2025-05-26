package service

import (
    "context"
    
    `dpcms/internal/database/model`
    `dpcms/internal/database/query`
    `dpcms/internal/http/dto`
    `dpcms/internal/http/errors/user_error`
    `dpcms/internal/infra`
    `dpcms/internal/packages/password`
    `github.com/jinzhu/copier`
)

type UserService struct {
    infra *infra.Infra
    query *query.Query
}

func NewUserService(infra *infra.Infra) *UserService {
    return &UserService{infra: infra, query: query.Use(infra.DB)}
}

func (srv *UserService) isUnique(ctx context.Context, u *model.User) error {
    q := srv.query.User
    findResult, err := q.WithContext(ctx).Where(q.Username.Eq(u.Username)).Or(q.Email.Eq(u.Email)).First()
    if err != nil {
        return err
    }
    if findResult != nil {
        if findResult.Username == u.Username {
            return user_error.ErrUsernameExists.ToError()
        }
        if findResult.Email == u.Email {
            return user_error.ErrEMailExists.ToError()
        }
    }
    return nil
}

func (srv *UserService) Create(ctx context.Context, userDTO *dto.UserRegisterRequest) (*dto.UserInfo, error) {
    u := &model.User{}
    err := copier.Copy(u, userDTO)
    if err != nil {
        return nil, err
    }
    
    err = srv.isUnique(ctx, u)
    if err != nil {
        return nil, err
    }
    hashedPwd, err := password.Password(u.Password).Make()
    if err != nil {
        return nil, err
    }
    u.Password = hashedPwd
    err = srv.query.User.WithContext(ctx).Create(u)
    if err != nil {
        return nil, err
    }
    userInfo := &dto.UserInfo{}
    err = copier.Copy(userInfo, u)
    if err != nil {
        return nil, err
    }
    return userInfo, nil
}

func (srv *UserService) FindUserWithCredential(ctx context.Context, username string, pwd string) (*model.User, error) {
    result, err := srv.FindByName(ctx, username)
    if err != nil {
        return nil, err
    }
    if !password.Password(result.Password).Compare(pwd) {
        return nil, user_error.ErrWrongPassword.ToError()
    }
    return result, nil
}

func (srv *UserService) FindByID(ctx context.Context, id int64) (*model.User, error) {
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

func (srv *UserService) UpdatePassword(ctx context.Context, id int64, pwd string) error {
    q := srv.query.User
    hashedPwd, err := password.Password(pwd).Make()
    if err != nil {
        return err
    }
    _, err = q.WithContext(ctx).Where(q.ID.Eq(id)).Update(q.Password, hashedPwd)
    return err
}
