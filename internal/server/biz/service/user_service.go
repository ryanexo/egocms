package service

import (
    `context`
    `errors`

    `GoBlog/internal/server/biz/repository`
    `GoBlog/internal/server/entity`
    `GoBlog/internal/server/erroz`
    `GoBlog/internal/server/pkg/password`
    `gorm.io/gorm`
)

type UserService struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
    return UserService{repo: repo}
}

func (srv UserService) Create(ctx context.Context, u *entity.User) error {
    _, err := srv.repo.Find(ctx, []string{"id"}, repository.Conditions{"username": u.Username})
    if err == nil {
        return erroz.ErrUnameExists.ToError()
    }
    if !errors.Is(err, gorm.ErrRecordNotFound) {
        return err
    }
    _, err = srv.repo.Find(ctx, []string{"id"}, repository.Conditions{"email": u.Email})
    if err == nil {
        return erroz.ErrMailExists.ToError()
    }
    if !errors.Is(err, gorm.ErrRecordNotFound) {
        return err
    }
    passwd, err := password.Make(u.Password)
    if err != nil {
        return err
    }
    u.Password = passwd
    return srv.repo.Create(ctx, u)
}

func (srv UserService) FindByID(ctx context.Context, id uint) (entity.User, error) {
    user, err := srv.repo.FindByID(ctx, id)
    if errors.Is(err, gorm.ErrRecordNotFound) {
        err = erroz.ErrAccountNotExists.ToError()
    }
    return user, err
}

func (srv UserService) FindByName(ctx context.Context, name string) (entity.User, error) {
    user, err := srv.repo.Find(ctx, nil, repository.Conditions{"username": name})
    if errors.Is(err, gorm.ErrRecordNotFound) {
        err = erroz.ErrAccountNotExists.ToError()
    }
    return user, err
}
