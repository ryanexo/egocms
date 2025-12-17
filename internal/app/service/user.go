package service

import (
    "context"
    `database/sql`
    
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/helper/dbscope`
    `dpcms/internal/app/service/internal/common`
    `dpcms/internal/app/service/srvparams`
    `dpcms/internal/database/model`
    `dpcms/internal/database/query`
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

func (srv UserService) isUnique(ctx context.Context, username string, email string) error {
    q := srv.query.User
    findResult, err := q.WithContext(ctx).Where(q.Username.Eq(username)).Or(q.Email.Eq(email)).First()
    if err != nil {
        return err
    }
    if findResult != nil {
        if findResult.Username == username {
            return erroz.ErrUsernameExists.ToError()
        }
        if findResult.Email == email {
            return erroz.ErrEMailExists.ToError()
        }
    }
    return nil
}

func (srv UserService) Create(ctx context.Context, params srvparams.UserCreateParams) (result srvparams.User, err error) {
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

func (srv UserService) FindByCredential(ctx context.Context, params srvparams.UserCredentialParams) (*srvparams.User, error) {
    uo := srv.query.User
    po := srv.query.UserProfile
    ro := srv.query.Role
    
    result := srvparams.User{}
    err := uo.WithContext(ctx).LeftJoin(po, uo.ID.EqCol(po.ID)).LeftJoin(ro, uo.RoleID.EqCol(ro.ID)).Select(uo.ALL, po.ALL, ro.Name.As("RoleName")).Where(uo.Username.Eq(params.Username)).Scan(&result)
    if err != nil {
        return nil, err
    }
    if !password.Password(result.Password).Compare(params.Password) {
        return nil, erroz.ErrWrongPassword.ToError()
    }
    return &result, nil
}

func (srv UserService) FindByID(ctx context.Context, id uint64) (result srvparams.User, err error) {
    q := srv.query.User
    u, err := q.WithContext(ctx).Where(q.ID.Eq(id)).First()
    if err != nil {
        return
    }
    err = copier.Copy(&result, &u)
    return
}

func (srv UserService) FindByName(ctx context.Context, name string) (result srvparams.User, err error) {
    q := srv.query.User
    u, err := q.WithContext(ctx).Where(q.Username.Eq(name)).First()
    if err != nil {
        return
    }
    err = copier.Copy(&result, &u)
    return
}

func (srv UserService) ResetPassword(ctx context.Context, id uint64, pwd string) error {
    q := srv.query.User
    hashedPwd, err := password.Password(pwd).Generate()
    if err != nil {
        return err
    }
    _, err = q.WithContext(ctx).Where(q.ID.Eq(id)).Update(q.Password, hashedPwd)
    return err
}

func (srv UserService) Delete(ctx context.Context, id uint64) error {
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

func (srv UserService) UpdateProfile(ctx context.Context, params srvparams.UserProfileUpdateParams) error {
    profile := model.UserProfile{}
    err := copier.Copy(&profile, params)
    if err != nil {
        return err
    }
    _, err = srv.query.WithContext(ctx).UserProfile.Select(
        srv.query.UserProfile.Gender,
        srv.query.UserProfile.Country,
        srv.query.UserProfile.Province,
        srv.query.UserProfile.City,
        srv.query.UserProfile.Nickname,
        srv.query.UserProfile.Description,
    ).Where(srv.query.UserProfile.UserID.Eq(params.ID)).Updates(&profile)
    return err
}

func (srv UserService) List(ctx context.Context, p srvparams.UserListQueryParams) (result common.PaginatedResult[srvparams.User], err error) {
    userDAO := srv.query.User
    profileDAO := srv.query.UserProfile
    q := userDAO.WithContext(ctx).LeftJoin(profileDAO, profileDAO.UserID.EqCol(userDAO.ID))
    if p.Username != nil {
        q = q.Where(userDAO.Username.Like(*p.Username + "%"))
    }
    if p.Email != nil {
        q = q.Where(userDAO.Email.Eq(*p.Email))
    }
    if p.Status != nil {
        q = q.Where(userDAO.Status.Eq(*p.Status))
    }
    if p.VerifyStartTime != nil && p.VerifyEndTime != nil {
        q = q.Where(
            userDAO.VerifiedAt.Gte(sql.NullTime{
                Valid: true,
                Time:  *p.VerifyStartTime,
            }),
        ).Where(
            userDAO.VerifiedAt.Lte(sql.NullTime{
                Valid: true,
                Time:  *p.VerifyEndTime,
            }),
        )
    }
    if p.IP != nil {
        q = q.Where(userDAO.IP.Like(*p.IP))
    }
    if p.RoleID != nil {
        q = q.Where(userDAO.RoleID.Eq(*p.RoleID))
    }
    if p.Nickname != nil {
        q = q.Where(profileDAO.Nickname.Eq(*p.Nickname))
    }
    if p.Gender != nil {
        q = q.Where(profileDAO.Gender.Eq(*p.Gender))
    }
    if p.Country != nil {
        q = q.Where(profileDAO.Country.Eq(*p.Country))
    }
    if p.Province != nil {
        q = q.Where(profileDAO.Province.Eq(*p.Province))
    }
    if p.City != nil {
        q = q.Where(profileDAO.City.Eq(*p.City))
    }
    
    count, err := q.Count()
    if err != nil {
        return
    }
    users, err := q.Scopes(dbscope.Paginate(p.PageNo, p.PageSize)).Find()
    if err != nil {
        return
    }
    
    userList := make([]srvparams.User, len(users))
    err = copier.Copy(&userList, &users)
    
    return common.PaginatedResult[srvparams.User]{
        PageNo:   p.PageNo,
        PageSize: p.PageSize,
        Total:    count,
        List:     userList,
    }, nil
}
