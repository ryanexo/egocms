package service

import (
    "context"
    `database/sql`
    
    `dpcms/internal/app/dto`
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/service/internal/common`
    `dpcms/internal/infra`
    `dpcms/internal/infra/password`
    `dpcms/internal/infra/persistence/dbscope`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    
    `github.com/jinzhu/copier`
)

type User struct {
    persist *query.Query
}

func NewUserService(i *infra.Infra) *User {
    return &User{persist: i.Query}
}

func (srv User) isUnique(ctx context.Context, username string, email string) error {
    q := srv.persist.User
    findResult, err := q.WithContext(ctx).Where(q.Username.Eq(username)).Or(q.Email.Eq(email)).First()
    if err != nil {
        return err
    }
    if findResult != nil {
        if findResult.Username == username {
            return erroz.UserNameExists.ToError()
        }
        if findResult.Email == email {
            return erroz.UserEmailExists.ToError()
        }
    }
    return nil
}

func (srv User) Create(ctx context.Context, params dto.UserCreateParams) (result dto.User, err error) {
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
    
    err = srv.persist.Transaction(func(tx *query.Query) error {
        return tx.User.WithContext(ctx).Create(&u)
    })
    if err != nil {
        return
    }
    profile := model.UserProfile{}
    err = srv.persist.User.Profile.Model(&u).Append(&profile)
    if err != nil {
        return
    }
    
    err = copier.Copy(&result, &u)
    return
}

func (srv User) FindByCredential(ctx context.Context, params dto.UserCredentialParams) (*dto.User, error) {
    uo := srv.persist.User
    po := srv.persist.UserProfile
    ro := srv.persist.Role
    
    result := dto.User{}
    err := uo.WithContext(ctx).LeftJoin(po, uo.ID.EqCol(po.ID)).LeftJoin(ro, uo.RoleID.EqCol(ro.ID)).Select(uo.ALL, po.ALL, ro.Name.As("RoleName")).Where(uo.Username.Eq(params.Username)).Scan(&result)
    if err != nil {
        return nil, err
    }
    if !password.Password(result.Password).Compare(params.Password) {
        return nil, erroz.UserWrongPasswd.ToError()
    }
    return &result, nil
}

func (srv User) FindByID(ctx context.Context, id uint64) (result dto.User, err error) {
    q := srv.persist.User
    u, err := q.WithContext(ctx).Where(q.ID.Eq(id)).First()
    if err != nil {
        return
    }
    err = copier.Copy(&result, &u)
    return
}

func (srv User) FindByName(ctx context.Context, name string) (result dto.User, err error) {
    q := srv.persist.User
    u, err := q.WithContext(ctx).Where(q.Username.Eq(name)).First()
    if err != nil {
        return
    }
    err = copier.Copy(&result, &u)
    return
}

func (srv User) ResetPassword(ctx context.Context, id uint64, pwd string) error {
    q := srv.persist.User
    hashedPwd, err := password.Password(pwd).Generate()
    if err != nil {
        return err
    }
    _, err = q.WithContext(ctx).Where(q.ID.Eq(id)).Update(q.Password, hashedPwd)
    return err
}

func (srv User) Delete(ctx context.Context, id uint64) error {
    return srv.persist.Transaction(func(tx *query.Query) error {
        q := tx.WithContext(ctx)
        _, err := q.User.Where(tx.User.ID.Eq(id)).Delete()
        if err != nil {
            return err
        }
        _, err = q.UserProfile.Where(tx.UserProfile.UserID.Eq(id)).Delete()
        if err != nil {
            return err
        }
        return nil
    })
}

func (srv User) UpdateProfile(ctx context.Context, params dto.UserProfileUpdateParams) error {
    profile := model.UserProfile{}
    err := copier.Copy(&profile, params)
    if err != nil {
        return err
    }
    profileDao := srv.persist.UserProfile
    _, err = srv.persist.WithContext(ctx).UserProfile.Select(
        profileDao.Gender,
        profileDao.Country,
        profileDao.Province,
        profileDao.City,
        profileDao.Nickname,
        profileDao.Description,
    ).Where(profileDao.UserID.Eq(params.ID)).Updates(&profile)
    return err
}

func (srv User) List(ctx context.Context, p dto.UserListQueryParams) (result common.PaginatedResult[dto.User], err error) {
    userDAO := srv.persist.User
    profileDAO := srv.persist.UserProfile
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
    
    userList := make([]dto.User, len(users))
    err = copier.Copy(&userList, &users)
    
    return common.PaginatedResult[dto.User]{
        PageNo:   p.PageNo,
        PageSize: p.PageSize,
        Total:    count,
        List:     userList,
    }, nil
}
