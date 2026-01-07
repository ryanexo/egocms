package service

import (
    "context"
    
    `dpcms/internal/app/user/errno`
    `dpcms/internal/app/user/internal/assembler`
    `dpcms/internal/app/user/internal/dto`
    `dpcms/internal/app/user/repo`
    "dpcms/internal/infra"
    "dpcms/internal/infra/password"
    "dpcms/internal/infra/persistence/datatype"
    "dpcms/internal/infra/persistence/model"
    "dpcms/internal/infra/persistence/query"
    `dpcms/internal/types`
)

type UserService struct {
    persist *query.Query
}

func NewUserService(i *infra.Infra) *UserService {
    return &UserService{persist: i.Query}
}

func (srv UserService) Create(ctx context.Context, params dto.UserCreateParams) (datatype.SafeUint64, error) {
    usrRepo := repo.NewUserRepo(srv.persist)
    user, err := usrRepo.FirstByUsernameOrEmail(ctx, params.Username, params.Email)
    if err != nil {
        return 0, err
    }
    if user.Username == params.Username {
        return 0, errno.UserNameExists.ToError()
    }
    if user.Email == params.Email {
        return 0, errno.UserEmailExists.ToError()
    }
    
    hashedPwd, err := password.Password(params.Password).Generate()
    if err != nil {
        return 0, err
    }
    data := assembler.BuildUserCreateCommand(&params)
    data.Password = hashedPwd
    
    if err = usrRepo.Create(ctx, data); err != nil {
        return 0, err
    }
    
    return data.ID, nil
}

func (srv UserService) FindByCredential(ctx context.Context, params dto.UserCredentialParams) (*dto.User, error) {
    data, err := repo.NewUserRepo(srv.persist).FindByUsername(ctx, params.Username)
    if err != nil {
        return nil, err
    }
    if !password.Password(data.Password).Compare(params.Password) {
        return nil, errno.UserWrongPasswd.ToError()
    }
    return assembler.BuildUserDTO(data), nil
}

func (srv UserService) FindByID(ctx context.Context, id datatype.SafeUint64) (*dto.User, error) {
    data, err := repo.NewUserRepo(srv.persist).FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return assembler.BuildUserDTO(data), nil
}

func (srv UserService) FindByName(ctx context.Context, name string) (*dto.User, error) {
    data, err := repo.NewUserRepo(srv.persist).FindByUsername(ctx, name)
    if err != nil {
        return nil, err
    }
    return assembler.BuildUserDTO(data), nil
}

func (srv UserService) ResetPassword(ctx context.Context, id datatype.SafeUint64, pwd string) error {
    hashedPwd, err := password.Password(pwd).Generate()
    if err != nil {
        return err
    }
    _, err = repo.NewUserRepo(srv.persist).UpdatePassword(ctx, id, hashedPwd)
    return err
}

func (srv UserService) ChangePassword(ctx context.Context, current *model.User, params dto.UserPasswdUpdateParams) error {
    if params.Password != params.PasswordConfirm {
        return errno.UserWrongConfirmPasswd.ToError()
    }
    if params.Password == params.OldPassword {
        return errno.UserEqualsOldPasswd.ToError()
    }
    _, err := srv.FindByCredential(ctx, dto.UserCredentialParams{
        Username: current.Username,
        Password: params.OldPassword,
    })
    if err != nil {
        return err
    }
    return srv.ResetPassword(ctx, current.ID, params.Password)
}

func (srv UserService) Delete(ctx context.Context, id datatype.SafeUint64) error {
    return srv.persist.Transaction(func(tx *query.Query) error {
        return repo.NewUserRepo(tx).Delete(ctx, id)
    })
}

func (srv UserService) UpdateProfile(ctx context.Context, params dto.UserProfile) error {
    data := assembler.BuildUserProfileModel(&params)
    _, err := repo.NewUserRepo(srv.persist).UpdateProfile(ctx, data)
    return err
}

func (srv UserService) List(ctx context.Context, params dto.UserListParams) (*types.PaginatedResult[*dto.User], error) {
    data, total, err := repo.NewUserRepo(srv.persist).List(ctx, &params)
    if err != nil {
        return nil, err
    }
    
    return &types.PaginatedResult[*dto.User]{
        Pagination: params.Pagination,
        Total:      total,
        List:       assembler.BuildUserListDTO(data),
    }, nil
}
