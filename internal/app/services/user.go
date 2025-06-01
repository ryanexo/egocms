package services

import (
    "context"
    `database/sql`
    
    `dpcms/internal/app/errors/user_error`
    `dpcms/internal/app/services/types`
    `dpcms/internal/app/services/types/user`
    `dpcms/internal/database/model`
    `dpcms/internal/database/query`
    `dpcms/internal/infra`
    `dpcms/internal/packages/password`
    `dpcms/internal/utils/dbscopes`
    `github.com/jinzhu/copier`
)

type UserService struct {
    infra *infra.Infra
    query *query.Query
}

func NewUserService(infra *infra.Infra) UserService {
    return UserService{infra: infra, query: query.Use(infra.DB)}
}

func (srv UserService) isUnique(ctx context.Context, username string, email string) error {
    q := srv.query.User
    findResult, err := q.WithContext(ctx).Where(q.Username.Eq(username)).Or(q.Email.Eq(email)).First()
    if err != nil {
        return err
    }
    if findResult != nil {
        if findResult.Username == username {
            return user_error.ErrUsernameExists.ToError()
        }
        if findResult.Email == email {
            return user_error.ErrEMailExists.ToError()
        }
    }
    return nil
}

func (srv UserService) Create(ctx context.Context, params user.RegisterParams) (result user.AccountDetail, err error) {
    err = srv.isUnique(ctx, params.Username, params.Email)
    if err != nil {
        return
    }
    
    hashedPwd, err := password.Password(params.Password).Generate()
    if err != nil {
        return
    }
    params.Password = hashedPwd
    
    u := model.User{}
    err = copier.Copy(&u, params)
    if err != nil {
        return
    }
    
    err = srv.query.User.WithContext(ctx).Create(&u)
    if err != nil {
        return
    }
    profile := model.UserProfile{}
    err = srv.query.User.Profile.Model(&u).Append(&profile)
    if err != nil {
        return
    }
    
    err = copier.Copy(&result, &u)
    return
}

func (srv UserService) IsValidCredential(ctx context.Context, username string, password string) (bool, error) {
    _, valid, err := srv.FindByCredential(ctx, user.CredentialParams{Username: username, Password: password})
    return valid, err
}

func (srv UserService) FindByCredential(ctx context.Context, params user.CredentialParams) (result user.AccountDetail, valid bool, err error) {
    u, err := srv.query.User.WithContext(ctx).Where(srv.query.User.Username.Eq(params.Username)).First()
    if err != nil {
        return
    }
    if !password.Password(u.Password).Compare(params.Password) {
        return
    }
    err = copier.Copy(&result, &u)
    if err != nil {
        valid = true
    }
    return
}

func (srv UserService) FindByID(ctx context.Context, id int64) (result user.AccountDetail, err error) {
    q := srv.query.User
    u, err := q.WithContext(ctx).Where(q.ID.Eq(id)).First()
    if err != nil {
        return
    }
    err = copier.Copy(&result, &u)
    return
}

func (srv UserService) FindByName(ctx context.Context, name string) (result user.AccountDetail, err error) {
    q := srv.query.User
    u, err := q.WithContext(ctx).Where(q.Username.Eq(name)).First()
    if err != nil {
        return
    }
    err = copier.Copy(&result, &u)
    return
}

func (srv UserService) ResetPassword(ctx context.Context, id int64, pwd string) error {
    q := srv.query.User
    hashedPwd, err := password.Password(pwd).Generate()
    if err != nil {
        return err
    }
    _, err = q.WithContext(ctx).Where(q.ID.Eq(id)).Update(q.Password, hashedPwd)
    return err
}

func (srv UserService) Delete(ctx context.Context, id int64) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        q := tx.WithContext(ctx)
        _, err := q.User.Where(srv.query.User.ID.Eq(id)).Delete()
        if err != nil {
            return err
        }
        _, err = q.UserProfile.Where(srv.query.UserProfile.UserID.Eq(id)).Delete()
        if err != nil {
            return err
        }
        return nil
    })
}

func (srv UserService) List(ctx context.Context, params user.ListRetrieveParams) (result types.Pagination[user.AccountDetail], err error) {
    userDAO := srv.query.User
    profileDAO := srv.query.UserProfile
    q := userDAO.WithContext(ctx).LeftJoin(profileDAO, profileDAO.UserID.EqCol(userDAO.ID))
    if params.Username != nil {
        q = q.Where(userDAO.Username.Like(*params.Username + "%"))
    }
    if params.Email != nil {
        q = q.Where(userDAO.Email.Eq(*params.Email))
    }
    if params.Status != nil {
        q = q.Where(userDAO.Status.Eq(*params.Status))
    }
    if params.VerifyStartTime != nil && params.VerifyEndTime != nil {
        q = q.Where(
            userDAO.VerifiedAt.Gte(sql.NullTime{
                Valid: true,
                Time:  *params.VerifyStartTime,
            }),
        ).Where(
            userDAO.VerifiedAt.Lte(sql.NullTime{
                Valid: true,
                Time:  *params.VerifyEndTime,
            }),
        )
    }
    if params.IP != nil {
        q = q.Where(userDAO.IP.Like(*params.IP))
    }
    if params.RoleID != nil {
        q = q.Where(userDAO.RoleID.Eq(*params.RoleID))
    }
    if params.Nickname != nil {
        q = q.Where(profileDAO.Nickname.Eq(*params.Nickname))
    }
    if params.Gender != nil {
        q = q.Where(profileDAO.Gender.Eq(*params.Gender))
    }
    if params.Country != nil {
        q = q.Where(profileDAO.Country.Eq(*params.Country))
    }
    if params.Province != nil {
        q = q.Where(profileDAO.Province.Eq(*params.Province))
    }
    if params.City != nil {
        q = q.Where(profileDAO.City.Eq(*params.City))
    }
    
    count, err := q.Count()
    if err != nil {
        return
    }
    users, err := q.Scopes(dbscopes.Paginate(params.PageNo, params.PageSize)).Find()
    if err != nil {
        return
    }
    
    userList := make([]user.AccountDetail, len(users))
    err = copier.Copy(&userList, &users)
    
    return types.Pagination[user.AccountDetail]{
        PageNo:   params.PageNo,
        PageSize: params.PageSize,
        Total:    count,
        List:     userList,
    }, nil
}
